package inventory

import (
	"path/filepath"
	"testing"
)

var testCfgfile = "../etc/infra-inventory.json.sample"

func Test_LoadConfig(t *testing.T) {
	cfg, err := LoadConfig(testCfgfile)
	if err != nil {
		t.Fatalf("%s", err)
	}

	if !filepath.IsAbs(cfg.Datastore.Config.MappingFile) {
		t.Fatalf("Mapping file not abs path\n")
	}
	t.Logf("%#v\n", *cfg)
}
