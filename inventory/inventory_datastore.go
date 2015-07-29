package inventory

import (
	"encoding/json"
	"fmt"
	log "github.com/golang/glog"
	elastigo "github.com/mattbaird/elastigo/lib"
)

type BasicAsset struct {
	Version int64 `json:"version"`
}

type InventoryDatastore struct {
	*ElasticsearchDatastore
}

func NewInventoryDatastore(ds *ElasticsearchDatastore) *InventoryDatastore {
	return &InventoryDatastore{ds}
}

func (ds *InventoryDatastore) GetAsset(assetType, assetId string) (asset elastigo.BaseResponse, err error) {
	if asset, err = ds.Conn.Get(ds.Index, assetType, assetId, nil); err == nil && asset.Found {
		return
	}
	err = fmt.Errorf("Not found: %s/%s %s", assetType, assetId, err)
	return
}

func (ds *InventoryDatastore) CreateAsset(assetType, assetId string, data interface{}) (string, error) {
	//log.V(10).Infof("%v\n", data)
	_, err := ds.GetAsset(assetType, assetId)
	if err == nil {
		return "", fmt.Errorf("Asset already exists: %s", assetId)
	}

	resp, err := ds.Conn.Index(ds.Index, assetType, assetId, nil, data)
	if err != nil {
		log.Warningf("%s\n", err)
		return "", err
	}

	if !resp.Created {
		return "", fmt.Errorf("Failed: %s", resp)
	}

	return resp.Id, nil
}

func (ds *InventoryDatastore) EditAsset(assetType, assetId string, data interface{}) (string, error) {

	asset, err := ds.GetAsset(assetType, assetId)
	if err != nil {
		return "", err
	}
	//log.V(12).Infof("%#v\n", asset)

	resp, err := ds.Conn.Update(ds.Index, assetType, assetId, nil, map[string]interface{}{"doc": data})
	if err != nil {
		return "", err
	}

	nid, err := ds.CreateAssetVersion(asset)
	if err != nil {
		log.Errorf("%s", err)
	} else {
		log.V(10).Infof("Version created: %s\n", nid)
	}

	return resp.Id, nil
}

//func (ds *InventoryDatastore) ListAssets(assetType string)                           {}

func (ds *InventoryDatastore) RemoveAsset(assetType, assetId string) bool {
	// asset not found
	asset, err := ds.GetAsset(assetType, assetId)
	if err != nil {
		return false
	}

	resp, err := ds.Conn.Delete(ds.Index, assetType, assetId, nil)
	if err != nil {
		log.Errorf("%s\n", err)
		return false
	}

	nid, err := ds.CreateAssetVersion(asset)
	if err != nil {
		log.Errorf("%s", err)
	} else {
		log.V(10).Infof("Version created: %s\n", nid)
	}

	return resp.Found
}

func (e *InventoryDatastore) ListAssetTypes() (types []string, err error) {
	var (
		b []byte
	)

	if b, err = e.Conn.DoCommand("GET", "/"+e.Index+"/_mapping", nil, nil); err != nil {
		return
	}

	m := map[string]map[string]map[string]interface{}{}
	if err = json.Unmarshal(b, &m); err != nil {
		return
	}
	// Remove _default_ map
	if _, ok := m[e.Index]["mappings"]["_default_"]; ok {
		delete(m[e.Index]["mappings"], "_default_")
	}

	types = make([]string, len(m[e.Index]["mappings"]))
	i := 0
	for k, _ := range m[e.Index]["mappings"] {
		types[i] = k
		i++
	}

	return
}

func (ds *InventoryDatastore) Search(assetType string, query interface{}) (elastigo.SearchResult, error) {
	//elastigo.Search(ds.Index)
	return ds.Conn.Search(ds.Index, assetType, nil, query)
}

func (ds *InventoryDatastore) CreateAssetVersion(asset elastigo.BaseResponse) (string, error) {
	var src map[string]interface{}
	if err := json.Unmarshal(*asset.Source, &src); err != nil {
		return "", err
	}

	versionedAssets, err := ds.GetAssetVersions(asset.Type, asset.Id, 1)
	if err != nil || versionedAssets.Hits.Len() < 1 {
		log.Warning("Creating new version anyway. Error=%s", err)
		src["version"] = 1
		//asset["_timestamp"] = asset.
	} else {
		var ba BasicAsset
		err := json.Unmarshal(*versionedAssets.Hits.Hits[0].Source, &ba)
		if err != nil {
			return "", err
		}
		src["version"] = ba.Version + 1
	}

	vresp, err := ds.Conn.Index(ds.VersionIndex, asset.Type,
		fmt.Sprintf("%s.%d", asset.Id, src["version"]), nil, src)
	if err != nil {
		return "", err
	}

	log.V(12).Infof("Version created: %s/%s.%d", asset.Type, asset.Id, src["version"])
	return vresp.Id, nil
}

/* Get a single version */
func (ds *InventoryDatastore) GetAssetVersion(assetType, assetId string, version int64) (asset elastigo.BaseResponse, err error) {

	if asset, err = ds.Conn.Get(ds.VersionIndex, assetType,
		fmt.Sprintf("%s.%d", version), nil); err == nil && asset.Found {
		return
	}
	err = fmt.Errorf("Not found (%s/%s.%d): %s", assetType, assetId, version, err)
	return
}

/* Get the last `count` versions */
func (ds *InventoryDatastore) GetAssetVersions(assetType, assetId string, count int64) (elastigo.SearchResult, error) {
	query := fmt.Sprintf(
		`{"query":{"prefix":{"_id": "%s"}},"sort":{"version":"desc"},"from":0,"size": %d}`,
		assetId, count)
	return ds.Conn.Search(ds.VersionIndex, assetType, nil, query)
}
