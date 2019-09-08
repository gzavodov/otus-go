package hw3

import (
	"testing"
)

type PackerTest struct {
	In  string
	Out string
}

func TestUnFold(t *testing.T) {
	tests := []PackerTest{
		PackerTest{In: "a4bc2d5e", Out: "aaaabccddddde"},
		PackerTest{In: "a11bc11", Out: "aaaaaaaaaaabccccccccccc"},
		PackerTest{In: "a1b0c", Out: "ac"},
		PackerTest{In: "a2b1c0", Out: "aab"},
		PackerTest{In: "abcd", Out: "abcd"},
		PackerTest{In: "1a1", Out: "a"},
		PackerTest{In: "45", Out: ""},
		PackerTest{In: `qwe\4\5`, Out: `qwe45`},
		PackerTest{In: `qwe\45`, Out: `qwe44444`},
		PackerTest{In: `qwe\\002`, Out: `qwe\\`},
		PackerTest{In: `qwe\\6`, Out: `qwe\\\\\\`},
		PackerTest{In: `qwe\\10`, Out: `qwe\\\\\\\\\\`},
	}

	packer := StringPacker{}
	for _, test := range tests {
		out, err := packer.Unpack(test.In)
		if err != nil {
			t.Errorf("ERROR: Unpack('%s') completed with error: '%s'\n", test.In, err.Error())
		} else {
			if out == test.Out {
				t.Logf("SUCCESS: '%s' => '%s'\n", test.In, out)
			} else {
				t.Errorf("FAILURE: '%s' => '%s' != '%s'\n", test.In, test.Out, out)
			}
		}
	}
}
