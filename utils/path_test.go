package tools

import (
	"testing"
)

func TestPathName(t *testing.T) {
	type testPara struct {
		input1 string
		want   string
	}

	test := map[string]testPara{
		"1": {`file/name1`, "name1"},
		"2": {`file/name2.2`, "name2.2"},
		"3": {`/file\name3.3`, "name3.3"},
	}

	for name, para := range test {
		t.Run(name, func(t *testing.T) {
			got := PathName(para.input1)
			testValue(t, got, para.want)
		})
	}
}

func TestPathStem(t *testing.T) {
	type testPara struct {
		input1 string
		want   string
	}

	test := map[string]testPara{
		"1": {`file/name1`, "name1"},
		"2": {`file/name2.2`, "name2"},
		"3": {`/file\name3.3`, "name3"},
	}

	for name, para := range test {
		t.Run(name, func(t *testing.T) {
			got := PathStem(para.input1)
			testValue(t, got, para.want)
		})
	}
}

func TestPathParent(t *testing.T) {
	type testPara struct {
		input1 string
		want   string
	}

	test := map[string]testPara{
		"1": {`/file/name1`, "/file"},
		"2": {`file/name2`, "file"},
		"3": {`./file/name3`, "./file"},
		"4": {`/file`, "/"},
		"5": {`file`, ""},
		"6": {`C:/file`, "C:"},
		"7": {`C:/file/name`, "C:/file"},
	}

	for name, para := range test {
		t.Run(name, func(t *testing.T) {
			got := PathParent(para.input1)
			testValue(t, got, para.want)
		})
	}
}

func TestCorrectIllegalityPath(t *testing.T) {
	type testPara struct {
		input1 string
		want   string
	}

	test := map[string]testPara{
		"1": {`/file/name`, `/file/name`},
		"2": {`\file/name`, `/file/name`},
		"3": {`/file/fi|le/name.a`, "/file/fi_le/name.a"},
		"4": {`/file/fi|le/name|/`, "/file/fi_le/name_"},
		"5": {`/file//fi|le//name.c`, "/file/_/fi_le_/name.c"},
		"6": {`//file//fi|le/fi*/name/`, `/_/file/_/fi_le/fi_/name`},
		"7": {`C:sss//file//fi|le/fi*/name/ss*.a`, `C:sss/_/file/_/fi_le/fi_/name/ss_.a`},
		"8": {`C:||//file//fi|le/fi*/name/ss*.a`, `C:||/_/file/_/fi_le/fi_/name/ss_.a`},
		"9": {`C:||/G://file//fi|le/fi*/name/ss|.a`, `C:||/G_/_/file/_/fi_le/fi_/name/ss_.a`},
	}

	for name, para := range test {
		t.Run(name, func(t *testing.T) {
			got := CorrectIllegalityPath(para.input1)
			testValue(t, got, para.want)
		})
	}
}

func TestJoinPath(t *testing.T) {

	t.Run("JoinPath1", func(t *testing.T) {
		got := JoinPath("./a", "b", "c")
		testValue(t, got, "./a/b/c")
	})
	t.Run("JoinPath2", func(t *testing.T) {
		got := JoinPath("./a/", "/b", "/c")
		testValue(t, got, "./a/b/c")
	})

}

func TestIsInFilter(t *testing.T) {
	type testPara struct {
		input1 string
		input2 []string
		want   bool
	}

	test := map[string]testPara{
		"1":  {`11.ini`, []string{`*.ini`}, true},
		"2":  {`11.in`, []string{`*.ini`}, false},
		"3":  {`11.abc`, []string{`11.ini`}, false},
		"4":  {`11.ini`, []string{`11.ini`}, true},
		"5":  {`1.ini`, []string{`11.ini`}, false},
		"6":  {`1.ini`, []string{`11.ini`, `*.ini`}, true},
		"7":  {`1.ini`, []string{`11.ini`, `1.ini`}, true},
		"8":  {`.ini`, []string{`*.ini`}, true},
		"9":  {`ini`, []string{`*.ini`}, false},
		"10": {`1.ini`, []string{`*.INI`}, false},
	}

	for name, para := range test {
		t.Run(name, func(t *testing.T) {
			got := IsInFilter(para.input1, para.input2)
			testValue(t, got, para.want)
		})
	}
}
