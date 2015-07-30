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

type AuthCachingConfig struct {
	TTL int64 `json:"ttl"`
}

type AuthConfig struct {
	Enabled    bool              `json:"enabled"`
	Type       string            `json:"type"`
	Config     ADAuthConfig      `json:"config"`
	Caching    AuthCachingConfig `json:"caching"`
	GroupsFile string            `json:"groups_file"`
}

type EssDatastoreConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Index       string `json:"index"`
	MappingFile string `json:"mapping_file"`
}

type EndpointsConfig struct {
	Prefix string `json:"prefix"`
}

type DatastoreConfig struct {
	Type      string             `json:"type"`
	Config    EssDatastoreConfig `json:"config"`
	BackupDir string             `json:"backup_dir"`
}

type AssetConfig struct {
	RequiredFields []string `json:"required_fields"`
}

type InventoryConfig struct {
	Auth      AuthConfig      `json:"auth"`
	Datastore DatastoreConfig `json:"datastore"`
	Endpoints EndpointsConfig `json:"endpoints"`
	AssetCfg  AssetConfig     `json:"asset"`
}

func LoadConfig(cfgfile string) (cfg *InventoryConfig, err error) {

	if !filepath.IsAbs(cfgfile) {
		cfgfile, _ = filepath.Abs(cfgfile)
	}

	var b []byte

	if b, err = ioutil.ReadFile(cfgfile); err != nil {
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

	if !filepath.IsAbs(cfg.Auth.GroupsFile) {
		cfg.Auth.GroupsFile, _ = filepath.Abs(cfg.Auth.GroupsFile)
	}

	return
}
