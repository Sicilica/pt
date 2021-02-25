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

	return time.Time{}, errors.New("failed to parse command args")
}

func (q *ArgsQueue) Rest() []string {
	r := q.args
	q.args = nil
	return r
}
