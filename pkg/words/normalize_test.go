package words

import (
	"testing"
)

func TestNormalizeLatin1(t *testing.T) {
	var tests = map[string]string{
		"ABC":      "abc",    // lowercase
		"ab\uFFFD": "ab ",    // replace invalid rune
		"ab\t\n":   "ab  ",   // replace control characters
		"ÅåÄäÖö":   "ååääöö", // retained
	}

	numErr := 0
	for input, expected := range tests {
		res := NormalizeLatin1(input)
		if res != expected {
			t.Errorf("Normalizing '%s' resulted in '%s', expected '%s'", input, res, expected)
			numErr++
		}
	}

	if numErr > 0 {
		t.Fatalf("%d Normalization failures.", numErr)
	}
}
