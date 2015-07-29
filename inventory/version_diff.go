package inventory

import (
	"encoding/json"
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
)

type VersionDiff struct {
	Version        int64  `json:"version"`
	AgainstVersion int64  `json:"against_version"`
	Diff           string `json:"diff"`
}

func GenerateDiff(prevName, prevStr, currName, currStr string) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(prevStr),
		B:        difflib.SplitLines(currStr),
		FromFile: prevName,
		ToFile:   currName,
		Context:  0,
	}
	return difflib.GetUnifiedDiffString(diff)
}

func parseVersion(ver interface{}) (verInt int64, err error) {
	switch ver.(type) {
	case float64:
		verF, _ := ver.(float64)
		verInt = int64(verF)
		return
	case int64:
		verInt, _ = ver.(int64)
		return
	case int:
		verI, _ := ver.(int)
		verInt = int64(verI)
		return
	default:
		err = fmt.Errorf("Could not parse version: %v", ver)
		return
	}
}

func GenerateVersionDiffs(versions ...map[string]interface{}) (list []VersionDiff, err error) {
	list = make([]VersionDiff, len(versions)-1)

	for i, version := range versions {
		if i+1 >= len(versions) {
			break
		}

		var (
			bi, bi1 []byte
			text    string
		)

		//fmt.Printf("%#v\n", versions[i])
		// Store and remove version before diff'ing
		var verInt int64
		if verInt, err = parseVersion(version["version"]); err != nil {
			return
		}
		delete(versions[i], "version")

		var verInt1 int64
		if verInt1, err = parseVersion(versions[i+1]["version"]); err != nil {
			return
		}
		delete(versions[i+1], "version")

		if bi, err = json.MarshalIndent(version, "", " "); err != nil {
			return
		}
		if bi1, err = json.MarshalIndent(versions[i+1], "", " "); err != nil {
			return
		}

		if text, err = GenerateDiff(
			fmt.Sprintf("v%d", verInt1), fmt.Sprintf("%s", bi1),
			fmt.Sprintf("v%d", verInt), fmt.Sprintf("%s", bi)); err != nil {
			return
		}
		//list[fmt.Sprintf("v%d", i+1)] = text
		list[i] = VersionDiff{Version: verInt, AgainstVersion: verInt1, Diff: text}
		// put next version back
		versions[i+1]["version"] = verInt1
	}
	return
}
