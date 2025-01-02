package utils_test

import (
	"github.com/kpearce2430/keputils/utils"
	"github.com/segmentio/encoding/json"
	"os"
	"testing"
)

var key = "ENV_KEY"
var value = "ENV_VALUE"

func TestMain(m *testing.M) {
	if err := os.Setenv(key, value); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestGetEnv(t *testing.T) {

	var badKey = "XKDNVMCKAW"
	var goodValue = "good"

	myValue := utils.GetEnv(key, "bad")

	if myValue != value {
		t.Error("Bad Value found")
	}

	myValue = utils.GetEnv(badKey, goodValue)

	if myValue != goodValue {
		t.Error("Good Value not found")
	}
}

func TestJulDate(t *testing.T) {

	t.Logf("%s", utils.JulDate())
}

func TestExists(t *testing.T) {

	ok, err := utils.Exists("somefile.txt")
	if err != nil {
		t.Error(err.Error())
	}

	switch {
	case ok:
		t.Error("There is no file called something.txt")
	case !ok:
		t.Log("File does not exists")
	}

	ok, err = utils.Exists("utils_test.go")
	if err != nil {
		t.Error(err.Error())
	}

	switch {
	case ok:
		t.Log("There is a file called utils_test.go")
	case !ok:
		t.Error("File does not exists")
	}
}

func TestFloatParse(t *testing.T) {

	type fpTests struct {
		Description   string
		Input         string
		Output        float64
		ExpectedError string
	}

	tests := []fpTests{
		{Description: "Float 1.00", Input: "1.00", Output: 1.0, ExpectedError: ""},
		{Description: "One Dollar", Input: "$1.00", Output: 1.0, ExpectedError: ""},
		{Description: "Some Junk", Input: "Junk", Output: 0.0, ExpectedError: "strconv.ParseFloat: parsing \"Junk\": invalid syntax"},
		{Description: "Thousand 99", Input: "$1,000.99", Output: 1000.99, ExpectedError: ""},
	}

	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			result, err := utils.FloatParse(tc.Input)
			if err != nil {
				if tc.ExpectedError != err.Error() {
					t.Log(err.Error())
					t.Fail()
				}
			}

			if result != tc.Output {
				t.Log("Bad result:", result)
				t.Fail()
			}
		})
	}
}

func TestAsciiString(t *testing.T) {
	type asciiTests struct {
		Description string
		Input       string
		Output      string
	}

	tests := []asciiTests{
		{Description: "Easy Test", Input: "hello world", Output: "hello world"},
		{Description: "Some Ascii", Input: string([]byte{0x0f, 'a', 'b', 'c'}), Output: "abc"},
		{Description: "French", Input: "Accéder", Output: "Accder"},
	}

	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			result := utils.AsciiString(tc.Input)

			if result != tc.Output {
				t.Log("Bad result:", result)
				t.Fail()
			}
		})
	}
}

func TestContains(t *testing.T) {

	type testCase struct {
		Description    string
		List           []string
		Search         string
		ExpectedResult bool
	}

	testCases := []testCase{
		{
			Description:    "Happy Path",
			List:           []string{"aaa", "bbb", "ccc"},
			Search:         "aaa",
			ExpectedResult: true,
		},
		{
			Description:    "Close 1",
			List:           []string{"aaa", "bbb", "ccc"},
			Search:         "aaaa",
			ExpectedResult: false,
		},
		{
			Description:    "Close 2",
			List:           []string{"aaa", "bbb", "ccc"},
			Search:         "bb",
			ExpectedResult: false,
		},
		{
			Description:    "Not even close",
			List:           []string{"aaa", "bbb", "ccc"},
			Search:         "xyz",
			ExpectedResult: false,
		},
		{
			Description:    "Not even close",
			List:           []string{"aaa", "bbb", "ccc"},
			Search:         "",
			ExpectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			result := utils.Contains(tc.List, tc.Search)
			if result != tc.ExpectedResult {
				t.Fail()
			}
		})
	}
}

func TestPrettyPrintJson(t *testing.T) {
	type myJson struct {
		Name     string `json:"name,omitempty"`
		Age      int    `json:"age,omitempty"`
		SomeBool bool   `json:"somebool,omitempty"`
	}
	myData := myJson{
		Name:     "bob",
		Age:      63,
		SomeBool: true,
	}

	data, err := json.Marshal(myData)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(utils.PrettyPrintJson(data))
}
