package inventory

import (
	"path/filepath"
	"testing"
)

var (
	testCfgfile = "../etc/infra-inventory.json.sample"
)

func Test_LoadConfig(t *testing.T) {
	cfg, err := LoadConfig(testCfgfile)
	if err != nil {
		t.Fatalf("%s", err)
	}

	if !filepath.IsAbs(cfg.Datastore.Config.MappingFile) {
		t.Fatalf("Mapping file not abs path\n")
	}
	if cfg.Auth.Caching.TTL != 7200 {
		t.Fatalf("TTL not set: %d", cfg.Auth.Caching.TTL)
	}

	if !filepath.IsAbs(cfg.Auth.GroupsFile) {
		t.Fatalf("Groups file not abs path: %s", cfg.Auth.GroupsFile)
	}
	t.Logf("%#v\n", *cfg)
}
