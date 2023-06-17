package utils_test

import (
	"github.com/kpearce2430/keputils/utils"
	"os"
	"testing"
)

var key = "ENV_KEY"
var value = "ENV_VALUE"

func TestMain(m *testing.M) {
	// log.Println("Do stuff BEFORE the tests!")
	os.Setenv(key, value)
	exitVal := m.Run()
	// log.Println("Do stuff AFTER the tests!")

	os.Exit(exitVal)
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
		Input         string
		Output        float64
		ExpectedError string
	}

	tests := []fpTests{
		{Input: "1.00", Output: 1.0, ExpectedError: ""},
		{Input: "$1.00", Output: 1.0, ExpectedError: ""},
		{Input: "Junk", Output: 0.0, ExpectedError: "strconv.ParseFloat: parsing \"Junk\": invalid syntax"},
		{Input: "$1,000.99", Output: 1000.99, ExpectedError: ""},
	}

	for _, tc := range tests {
		result, err := utils.FloatParse(tc.Input)
		if err != nil {
			if tc.ExpectedError == err.Error() {
				continue
			}
			t.Log(err.Error())
			t.Fail()
			continue
		}

		if result != tc.Output {
			t.Log("Bad result:", result)
			t.Fail()
		}

	}
}
