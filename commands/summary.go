package commands

import (
	"fmt"
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

	fmt.Printf("Tasks %s to %s\n", tp.Start.Format(time.DateTime), tp.End.Format(time.DateTime))
	fmt.Println("-----")

	var total time.Duration
	for _, t := range tasks {
		d := t.Stop.Sub(t.Start)
		total += d

		tags, err := c.PT().GetTaskTags(t)
		if err != nil {
			return err
		}

		fmt.Printf("%s - %s %s\n", formatTime(t.Start), formatDuration(d), formatTags(tags))
	}

	fmt.Println("-----")
	fmt.Println("total:", formatDuration(total))

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
