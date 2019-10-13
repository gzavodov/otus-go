package envdir

import (
	"io/ioutil"
	"regexp"
	"runtime"
	"testing"
)

func TestEnvironmentModifier(t *testing.T) {

	var dirPath string
	var executablePath string
	var expectedFilePath string

	if runtime.GOOS == "windows" {
		dirPath = ".\\test\\env"
		executablePath = ".\\test\\windows\\test.cmd"
		expectedFilePath = ".\\test\\test.txt"

	} else { //linux || darwin
		dirPath = "./test/env"
		executablePath = "./test/linux/test.sh"
		expectedFilePath = "./test/test.txt"
	}

	envmodifier := EnvModifier{}
	resultBytes, err := envmodifier.Run(dirPath, executablePath)
	if err != nil {
		t.Fatal(err)
	}

	//Reading of reference file
	expectedBytes, err := ioutil.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatal(err)
	}

	//Remove all spaces from string before comparison
	re := regexp.MustCompile(`\s+`)
	expected := re.ReplaceAllString(string(expectedBytes), "")
	result := re.ReplaceAllString(string(resultBytes), "")

	if expected != result {
		t.Error("result file does not contain expected data")
	}
}
