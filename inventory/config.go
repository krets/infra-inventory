package inventory

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type ADAuthConfig struct {
	Url          string `json:"url"`
	SearchBase   string `json:"search_base"`
	BindDN       string `json:"bind_dn"`
	BindPassword string `json:"bind_password"`
}

type AuthConfig struct {
	Enabled bool         `json:"enabled"`
	Type    string       `json:"type"`
	Config  ADAuthConfig `json:"config"`
}

type EssDatastoreConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Index       string `json:"index"`
	MappingFile string `json:"mapping_file"`
}

type DatastoreConfig struct {
	Type      string             `json:"type"`
	Config    EssDatastoreConfig `json:"config"`
	BackupDir string             `json:"backup_dir"`
}

type InventoryConfig struct {
	Auth      AuthConfig      `json:"auth"`
	Datastore DatastoreConfig `json:"datastore"`
}

func LoadConfig(cfgfile string) (cfg *InventoryConfig, err error) {

	if !filepath.IsAbs(cfgfile) {
		cfgfile, _ = filepath.Abs(cfgfile)
	}

	b, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &cfg); err != nil {
		return
	}

	if !filepath.IsAbs(cfg.Datastore.Config.MappingFile) {
		cfg.Datastore.Config.MappingFile, _ = filepath.Abs(cfg.Datastore.Config.MappingFile)
	}

	if !filepath.IsAbs(cfg.Datastore.BackupDir) {
		cfg.Datastore.BackupDir, _ = filepath.Abs(cfg.Datastore.BackupDir)
	}
	return
}
