package inventory

import (
	"encoding/json"
	"fmt"
	log "github.com/golang/glog"
	elastigo "github.com/mattbaird/elastigo/lib"
)

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
	//ds.Conn.Update(index, _type, id, args, data)
	fmt.Printf("%#v\n", data)
	resp, err := ds.Conn.Update(ds.Index, assetType, assetId, nil, data)
	if err != nil {
		return "", err
	}
	return resp.Id, nil

}

//func (ds *InventoryDatastore) ListAssets(assetType string)                           {}

func (e *InventoryDatastore) RemoveAsset(assetType, assetId string) bool {
	resp, err := e.Conn.Delete(e.Index, assetType, assetId, nil)
	if err != nil {
		log.Errorf("%s\n", err)
		return false
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
