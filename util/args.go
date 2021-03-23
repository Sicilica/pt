package util

import (
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var userInputTimeOffsetComponents = []time.Duration{
	24 * time.Hour,
	time.Hour,
	time.Minute,
	time.Second,
}
var userInputTimeOffsetRegexp = regexp.MustCompile(`^(?:(\d+)d)?(?:(\d+)h)?(?:(\d+)m)?(?:(\d+)s)?$`)

var userInputDateRegexp = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)

var userInputTimeRegexp = regexp.MustCompile(`^(\d{2}):(\d{2})$`)

type ArgsQueue struct {
	args []string
}

func NewArgsQueue(args []string) *ArgsQueue {
	return &ArgsQueue{
		args: args,
	}
}

func (q *ArgsQueue) MustBeEmpty() error {
	if len(q.args) > 0 {
		return errors.Errorf("unexpected arg \"%s\"", q.args[0])
	}
	return nil
}

func (q *ArgsQueue) MustPop(val string) {
	actual, err := q.Pop()
	if err != nil {
		panic(errors.Wrap(err, "MustPop encountered error"))
	}
	if actual != val {
		panic(errors.New("MustPop got wrong value"))
	}
}

func (q *ArgsQueue) Peek() string {
	if len(q.args) == 0 {
		return ""
	}
	return q.args[0]
}

func (q *ArgsQueue) Pop() (string, error) {
	if len(q.args) == 0 {
		return "", errors.New("failed to parse command args")
	}
	a := q.args[0]
	q.args = q.args[1:]
	return a, nil
}

func (q *ArgsQueue) PopTime(referenceTime time.Time, positiveOffsets bool) (time.Time, error) {
	// If we return early due to an error, we want to restore all popped args
	finalArgs := q.args
	defer func() {
		q.args = finalArgs
	}()

	s, err := q.Pop()
	if err != nil {
		return time.Time{}, err
	}

	// Attempt to parse as offset
	m := userInputTimeOffsetRegexp.FindStringSubmatch(s)
	if m != nil {
		o := 0 * time.Second

		for i, scale := range userInputTimeOffsetComponents {
			s := m[i+1]
			if s != "" {
				val, err := strconv.Atoi(s)
				if err != nil {
					return time.Time{}, errors.Wrap(err, "failed to parse command args")
				}
				o += time.Duration(val) * scale
			}
		}

		finalArgs = q.args
		if positiveOffsets {
			return referenceTime.Add(o), err
		}
		return referenceTime.Add(-o), err
	}

	// Attempt to parse as absolute time
	var absTime time.Time
	// Get date from input or assume from reference time
	m = userInputDateRegexp.FindStringSubmatch(s)
	if m != nil {
		y, err := strconv.Atoi(m[1])
		if err != nil {
			return time.Time{}, errors.New("malformed date")
		}
		month, err := strconv.Atoi(m[2])
		if err != nil {
			return time.Time{}, errors.New("malformed date")
		}
		d, err := strconv.Atoi(m[3])
		if err != nil {
			return time.Time{}, errors.New("malformed date")
		}
		absTime = time.Date(y, time.Month(month), d, 0, 0, 0, 0, time.Local)

		s, err = q.Pop()
		if err != nil {
			return time.Time{}, err
		}
	} else {
		y, m, d := referenceTime.Local().Date()
		absTime = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	}
	// Get time
	m = userInputTimeRegexp.FindStringSubmatch(s)
	if m == nil {
		return time.Time{}, errors.New("failed to parse command args")
	}
	h, err := strconv.Atoi(m[1])
	if err != nil {
		return time.Time{}, errors.New("malformed time")
	}
	minute, err := strconv.Atoi(m[2])
	if err != nil {
		return time.Time{}, errors.New("malformed time")
	}
	absTime = absTime.Add(time.Hour * time.Duration(h) + time.Minute * time.Duration(minute))
	finalArgs = q.args
	return absTime, nil
}

func (q *ArgsQueue) Rest() []string {
	r := q.args
	q.args = nil
	return r
}
