package words

import (
	"testing"
)

func TestSortWord(t *testing.T) {
	var tests = map[string]string{
		"":          "",
		"SUDRELSIT": "DEILRSSTU",
		"ZZzzAAaa":  "AAZZaazz",
	}

	numErr := 0
	for input, expected := range tests {
		res := sortWord(input)
		if res != expected {
			t.Errorf("Sorting '%s' resulted in '%s', expected '%s'", input, res, expected)
			numErr++
		}
	}

	if numErr > 0 {
		t.Errorf("%d SortWord failures.", numErr)
	}
}

func TestDeriveKey32(t *testing.T) {
	var tests = map[string]uint32{
		"deilrsstu": 0x001e0918,
		"az":        0x02000001,
		"a":         0x00000001,
		"z":         0x02000000,
		"INVALIDb":  0x80000002,
	}

	numErr := 0
	for input, expected := range tests {
		res := deriveKey32(input)
		if res != expected {
			t.Errorf("DeriveKey32 '%s' resulted in '%032b' (%08x), expected '%032b' (%08x)", input, res, res, expected, expected)
			numErr++
		}
	}

	if numErr > 0 {
		t.Errorf("%d SortWord failures.", numErr)
	}
}

func TestVerifyWord(t *testing.T) {

	type TestCase struct {
		input    string
		target   string
		expected bool
	}

	tests := []TestCase{
		{"in", "abiijn", true},
		{"in", "abcdeout", false},
		{sortWord("strudels"), "deilrsstu", true},
		{sortWord("strudels"), "deilrsst", false},
	}

	numErr := 0
	for _, test := range tests {
		res := verifyWord(test.input, test.target)
		if res != test.expected {
			t.Errorf("VerifyWord '%s' against '%s' resulted in %v, expected %v", test.input, test.target, res, test.expected)
			numErr++
		}
	}

	if numErr > 0 {
		t.Fatalf("%d VerifyWord failures.", numErr)
	}
}
