package inventory

import (
	"encoding/json"
	"fmt"
	"testing"
)

var (
	testPrev = map[string]interface{}{
		"attr1":   "val1",
		"attr2":   "val2",
		"attr4":   "val4",
		"attr3":   "val3",
		"version": 1,
	}
	testCurr = map[string]interface{}{
		"attr1":   "val1",
		"attr2":   "val2",
		"attr3":   "val5",
		"attr4":   "val4",
		"version": 2,
	}
	testCurr1 = map[string]interface{}{
		"attr1":   "val1",
		"attr2":   "val2",
		"attr3":   "val5",
		"attr4":   "val7",
		"version": 3,
	}
)

func Test_GenerateVersionDiffs(t *testing.T) {
	list, err := GenerateVersionDiffs(testCurr1, testCurr, testPrev)
	if err != nil {
		t.Fatalf("%s", err)
	}
	for _, v := range list {
		t.Logf("\n[ v%d ]\n%s\n", v.Version, v.Diff)
	}
}

func Test_GenerateDiff(t *testing.T) {
	bp, _ := json.MarshalIndent(testPrev, "", " ")
	bc, _ := json.MarshalIndent(testCurr, "", " ")

	diffText, err := GenerateDiff("previous", fmt.Sprintf("%s", bp), "current", fmt.Sprintf("%s", bc))

	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%s\n", diffText)
}
