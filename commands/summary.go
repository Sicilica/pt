package commands

import (
	"fmt"
	"slices"
	"time"

	"github.com/pkg/errors"

	"github.com/sicilica/pt/types"
)

func init() {
	register("summary <time period> [tag]", commandSummary, "Info", "displays historical tasks that match the given query")
}

func commandSummary(c types.CommandContext) error {
	tp, err := parseRequestedTime(c)
	if err != nil {
		return err
	}

	tag, _ := c.Args().Pop()
	err = c.Args().MustBeEmpty()
	if err != nil {
		return err
	}

	tasks, err := c.PT().QueryTasks(tp.Start, tp.End, &types.Query{Tag: tag})
	if err != nil {
		return err
	}

	fmt.Printf("Tasks %s to %s\n", tp.Start.Format(time.DateOnly), tp.End.Format(time.DateOnly))
	fmt.Println("-----")

	var total time.Duration
	tagTimes := make(map[string]time.Duration)
	for _, t := range tasks {
		start := t.Start
		stop := t.Stop
		if start.Before(tp.Start) {
			start = tp.Start
		}
		if stop.After(tp.End) {
			stop = tp.End
		}
		d := stop.Sub(start)
		total += d

		tags, err := c.PT().GetTaskTags(t)
		if err != nil {
			return err
		}

		for _, tag := range tags {
			tagTimes[tag] += d
		}

		fmt.Printf("%s - %s %s\n", formatTime(t.Start), formatDuration(d), formatTags(tags))
	}

	fmt.Println("-----")
	untracked := tp.End.Sub(tp.Start) - total
	if tp.End.After(time.Now()) {
		untracked = time.Since(tp.Start) - total
	}
	fmt.Println("total:", formatDuration(total))
	if untracked > 0 {
		fmt.Printf("untracked: %s (%.0f%%)\n", formatDuration(untracked), 100*(float64(untracked)/float64(untracked+total)))
	}

	// Fetch all relevant tags and their parents
	tags := make([]string, 0, len(tagTimes))
	for t := range tagTimes {
		tags = append(tags, t)
	}
	tagParents := make(map[string]string)
	for len(tags) > 0 {
		t := tags[len(tags)-1]
		tags = tags[:len(tags)-1]

		if _, ok := tagParents[t]; ok {
			continue
		}

		p, err := c.PT().GetTagParent(t)
		if err != nil && err != types.ErrNoParent {
			return err
		}
		tagParents[t] = p
		if p != "" {
			tags = append(tags, p)
		}
	}

	// Aggregate everything
	tagTree := make(map[string][]string)
	for t, p := range tagParents {
		tagTree[p] = append(tagTree[p], t)
	}
	var addTagTimes func(string)
	addTagTimes = func(t string) {
		for _, c := range tagTree[t] {
			addTagTimes(c)
			tagTimes[t] += tagTimes[c]
		}
	}
	addTagTimes("")

	// And print
	fmt.Println("")
	fmt.Println("breakdown")
	fmt.Println("-----")
	if _, ok := tagTimes["(untracked)"]; !ok {
		tagTimes["(untracked)"] = untracked
		tagTree[""] = append(tagTree[""], "(untracked)")
	}
	var printTag func(string, string)
	printTag = func(t, prefix string) {
		if t != "" {
			fmt.Printf("%s%s %s\n", prefix, t, formatDuration(tagTimes[t]))
			prefix += "| "
		}

		slices.SortFunc(tagTree[t], func(a, b string) int {
			d := tagTimes[b] - tagTimes[a]
			if d > 0 {
				return 1
			}
			if d < 0 {
				return -1
			}
			return 0
		})

		for _, c := range tagTree[t] {
			printTag(c, prefix)
		}
	}
	printTag("", "")

	return nil
}

type timePeriod struct {
	Start time.Time
	End   time.Time
}

func (tp *timePeriod) Add(d time.Duration) *timePeriod {
	if tp != nil {
		tp.Start = tp.Start.Add(d)
		tp.End = tp.End.Add(d)
	}
	return tp
}

var timePeriodParseRules = map[string]func(c types.CommandContext) (*timePeriod, error){
	"last": func(c types.CommandContext) (*timePeriod, error) {
		nextArg, err := c.Args().Pop()
		if err != nil {
			return nil, err
		}

		switch nextArg {
		case "week":
			y, m, d := time.Now().Local().Date()
			today := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
			end := today.Add(time.Duration(today.Weekday()) * -24 * time.Hour)
			return &timePeriod{
				Start: end.Add(-7 * 24 * time.Hour),
				End:   end,
			}, nil
		case "month":
			y, m, _ := time.Now().Local().Date()
			return &timePeriod{
				Start: time.Date(y, m-1, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(y, m, 1, 0, 0, 0, 0, time.Local),
			}, nil
		case "quarter":
			y, m, _ := time.Now().Local().Date()
			return &timePeriod{
				Start: time.Date(y, ((m+2)/3-1)*3-2, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(y, ((m+2)/3)*3-2, 1, 0, 0, 0, 0, time.Local),
			}, nil
		case "year":
			y, _, _ := time.Now().Local().Date()
			return &timePeriod{
				Start: time.Date(y-1, 1, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(y, 1, 1, 0, 0, 0, 0, time.Local),
			}, nil
		default:
			return nil, errors.Errorf("unrecognized time period \"%s\"", nextArg)
		}
	},
	"this": func(c types.CommandContext) (*timePeriod, error) {
		nextArg, err := c.Args().Pop()
		if err != nil {
			return nil, err
		}

		switch nextArg {
		case "week":
			y, m, d := time.Now().Local().Date()
			today := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
			start := today.Add(time.Duration(today.Weekday()) * -24 * time.Hour)
			return &timePeriod{
				Start: start,
				End:   start.Add(7 * 24 * time.Hour),
			}, nil
		case "month":
			y, m, _ := time.Now().Local().Date()
			return &timePeriod{
				Start: time.Date(y, m, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(y, m+1, 1, 0, 0, 0, 0, time.Local),
			}, nil
		case "quarter":
			y, m, _ := time.Now().Local().Date()
			return &timePeriod{
				Start: time.Date(y, ((m+2)/3)*3-2, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(y, ((m+2)/3+1)*3-2, 1, 0, 0, 0, 0, time.Local),
			}, nil
		case "year":
			y, _, _ := time.Now().Local().Date()
			return &timePeriod{
				Start: time.Date(y, 1, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(y+1, 1, 1, 0, 0, 0, 0, time.Local),
			}, nil
		default:
			return nil, errors.Errorf("unrecognized time period \"%s\"", nextArg)
		}
	},
	"past": func(c types.CommandContext) (*timePeriod, error) {
		nextArg, err := c.Args().Pop()
		if err != nil {
			return nil, err
		}

		switch nextArg {
		case "week":
			y, m, d := time.Now().Local().Date()
			today := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
			return &timePeriod{
				Start: today.Add(-6 * 24 * time.Hour),
				End:   today.Add(24 * time.Hour),
			}, nil
		case "month":
			y, m, d := time.Now().Local().Date()
			today := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
			return &timePeriod{
				Start: today.Add(-30 * 24 * time.Hour),
				End:   today.Add(24 * time.Hour),
			}, nil
		case "quarter":
			y, m, d := time.Now().Local().Date()
			today := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
			return &timePeriod{
				Start: today.Add(-90 * 24 * time.Hour),
				End:   today.Add(24 * time.Hour),
			}, nil
		case "year":
			y, m, d := time.Now().Local().Date()
			today := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
			return &timePeriod{
				Start: today.Add(-365 * 24 * time.Hour),
				End:   today.Add(24 * time.Hour),
			}, nil
		default:
			return nil, errors.Errorf("unrecognized time period \"%s\"", nextArg)
		}
	},
	"today": func(c types.CommandContext) (*timePeriod, error) {
		y, m, d := time.Now().Local().Date()
		start := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
		return &timePeriod{
			Start: start,
			End:   start.Add(24 * time.Hour),
		}, nil
	},
	"yesterday": func(c types.CommandContext) (*timePeriod, error) {
		y, m, d := time.Now().Local().Date()
		start := time.Date(y, m, d-1, 0, 0, 0, 0, time.Local)
		return &timePeriod{
			Start: start,
			End:   start.Add(24 * time.Hour),
		}, nil
	},
}

func parseRequestedTime(c types.CommandContext) (*timePeriod, error) {
	firstArg, err := c.Args().Pop()
	if err != nil {
		return nil, err
	}

	if rule, ok := timePeriodParseRules[firstArg]; ok {
		return rule(c)
	}

	return nil, errors.Errorf("unrecognized time period \"%s\"", firstArg)
}
