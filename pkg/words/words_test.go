package words

import (
	"github.com/go-test/deep"
	"strings"
	"testing"
)

func checkEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()
	diff := deep.Equal(actual, expected)
	if diff == nil {
		return
	} else if len(diff) > 0 {
		for _, d := range diff {
			t.Log("\t -- \t", d)
		}
		t.Error("checkEqual failed.")
	}
}

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

	checkEqual(t, res, expected)
}
