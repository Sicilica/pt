package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPopTime(t *testing.T) {
	type testCase struct {
		Desc string
		In []string
		ExpectedError string
		Out time.Time
	}

	referenceTime := time.Now()

	testCases := []testCase{
		{
			Desc: "Offset min/sec",
			In: []string{"1m44s"},
			Out: referenceTime.Add(-(time.Minute + 44 * time.Second)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Desc, func(t *testing.T) {
			res, err := NewArgsQueue(tc.In).PopTime(referenceTime, false)

			if tc.ExpectedError != "" {
				assert.EqualError(t, err, tc.ExpectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.Out, res)
			}
		})
	}
}
