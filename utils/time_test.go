package tools

import (
	"testing"
	"time"
)

func TestToHumanizeDurtion(t *testing.T) {
	type testPara struct {
		input time.Duration
		want  string
	}

	test := map[string]testPara{
		"0": {time.Hour, "1:00:00"},
		"1": {time.Minute, "1:00"},
		"2": {time.Second, "0:01"},
		"3": {time.Hour + time.Minute + time.Second, "1:01:01"},
		"4": {time.Hour*10 + time.Minute*15 + time.Second*18, "10:15:18"},
		"5": {time.Hour*10 + time.Minute*15, "10:15:00"},
	}

	for name, para := range test {
		t.Run(name, func(t *testing.T) {
			got := ToHumanizeDurtion(para.input)
			testValue(t, got, para.want)
		})
	}
}
