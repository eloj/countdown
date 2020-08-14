package words

import (
	"strings"
	"testing"

	test "github.com/eloj/go-eddy/pkg/test"
)

func TestFindWord(t *testing.T) {

	cd := NewCountdown(4, 8)

	dict := `
BASE
BILGE
BASIL
GABLES
BAGLESS
KISSABLE
`
	cnt, err := cd.AddDictionary(strings.NewReader(dict))

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if cnt != 6 {
		t.Errorf("Expected 6 words in dictionary, got %d", cnt)
	}

	limit := 10
	res := cd.FindWords("IBASELGSK", limit, 3)

	expected := NewFindWordsResult(limit)
	expected.Query = "ibaselgsk"
	expected.NumChecked = 5
	expected.NumHits = 3
	expected.NumDistFail = 2
	expected.Words = []WordDistResult{{"kissable", 1}, {"bagless", 2}, {"gables", 3}}

	test.CheckEqual(t, res, expected)
}
