package inventory

import (
	"testing"
)

var (
	testGroupsFile = "../etc/local-groups.json"
)

func Test_LoadLocalAuthGroups(t *testing.T) {
	groups, err := LoadLocalAuthGroups(testGroupsFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if len(groups) < 1 {
		t.Fatalf("no groups")
	}
}

func Test_LocalAuthGroups_UserHasGroupMembership(t *testing.T) {
	groups, _ := LoadLocalAuthGroups(testGroupsFile)

	if !groups.UserHasGroupMembership("absp", "admin") {
		t.Fatalf("Should exist")
	}
	if groups.UserHasGroupMembership("jibber", "admin") {
		t.Fatalf("Should not exist")
	}
}
