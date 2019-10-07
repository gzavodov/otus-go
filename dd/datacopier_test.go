package dd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

type TestInfo struct {
	Name       string
	SrcPath    string
	DstPath    string
	Offset     int64
	Limit      int64
	ResultPath string
}

func TestFileCopying(t *testing.T) {
	tests := []TestInfo{
		TestInfo{Name: "Patitional Text File Copying #1", SrcPath: "./test/in.txt", DstPath: "./test/out.txt", Offset: 47613, Limit: 40912, ResultPath: "./test/test1.txt"},
		TestInfo{Name: "Patitional Text File Copying #2", SrcPath: "./test/in.txt", DstPath: "./test/out.txt", Offset: 359287, Limit: 31218, ResultPath: "./test/test2.txt"},
		TestInfo{Name: "Text File Copying", SrcPath: "./test/in.txt", DstPath: "./test/out.txt", Offset: 0, Limit: 0, ResultPath: "./test/in.txt"},
	}

	copier := DataCopier{}
	for _, test := range tests {
		err := copier.Copy(test.SrcPath, test.DstPath, test.Offset, test.Limit)
		if err != nil {
			t.Fatal(err)
		}

		resultBytes, err := ioutil.ReadFile(test.DstPath)
		if err != nil {
			t.Fatal(err)
		}

		expectedBytes, err := ioutil.ReadFile(test.ResultPath)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(resultBytes, expectedBytes) {
			t.Error("result file does not contain expected data")
		}

		file, err := os.OpenFile(test.DstPath, os.O_WRONLY, 0666)
		if err != nil {
			t.Fatal(err)
		}

		file.Truncate(0)
		file.Sync()
		file.Close()
	}
}
