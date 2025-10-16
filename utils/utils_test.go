package tools

import (
	"testing"
)

func TestCmpVersion(t *testing.T) {
	type testPara struct {
		input1 string
		input2 string
		want   int
	}

	test := map[string]testPara{
		"0": {"13.1", "13.1.0", 0},
		"1": {"13.11", "13.1.0", 1},
		"2": {"13.0.0", "13.0.1", -1},
		"3": {"13", "13.0.00", 0},
		"4": {"13.010", "13.10", 0},
		"5": {"15.010", "15.10.0.0", 0},
		"6": {"2.0.2", "2.1.0", -1},
		"7": {"2.0.2", "2.0.1", 1},
		"8": {"2.3.2", "2.3.2", 0},
		"9": {"2.003.02.1", "2.003.02", 1},
	}

	for name, para := range test {
		t.Run(name, func(t *testing.T) {
			got := CmpVersion(para.input1, para.input2)
			testValue(t, got, para.want)
		})
	}
}
