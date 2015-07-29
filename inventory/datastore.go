package inventory

import (
	"encoding/json"
	"fmt"
	log "github.com/golang/glog"
	elastigo "github.com/mattbaird/elastigo/lib"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type IDatastore interface {
	GetAsset(assetType, assetId string) (elastigo.BaseResponse, error)
	GetAssetVersion(assetType, assetId string, version int64) (elastigo.BaseResponse, error)
	GetAssetVersions(assetType, assetId string, count int64) (elastigo.SearchResult, error)

	CreateAsset(assetType, assetId string, data interface{}, createType bool) (string, error)
	EditAsset(assetType, assetId string, data interface{}) (string, error)
	RemoveAsset(assetType, assetId string) bool
	//ListAssets(assetType string)
	ListAssetTypes() ([]string, error)
	Search(assetType string, query interface{}) (elastigo.SearchResult, error)
}

type ElasticsearchVersion struct {
	Number         string `json:"number"`
	BuildHash      string `json:"build_hash"`
	BuildTimestamp string `json:"build_timestamp"`
	BuildSnapshot  bool   `json:"build_snapshot"`
	LuceneVersion  string `json:"lucene_version"`
}

type EssInfo struct {
	Status      int64                `json:"status"`
	Name        string               `json:"name"`
	ClusterName string               `json:"cluster_name"`
	Version     ElasticsearchVersion `json:"version"`
	Tagline     string               `json:"tagline"`
}

/*
   Generic wrapper to elasticsearch
*/
type EssMapping struct {
	Meta             map[string]interface{} `json:"_meta"`
	DynamicTemplates []interface{}          `json:"dynamic_templates"`
}

type ElasticsearchDatastore struct {
	Conn         *elastigo.Conn
	Index        string
	VersionIndex string
}

/*
   Create the index if it does not exist.
   Optionally apply a mapping if mapping file is supplied.
*/
func NewElasticsearchDatastore(esshost string, essport int, index string, mappingfile ...string) (*ElasticsearchDatastore, error) {

	ed := ElasticsearchDatastore{
		Conn:         elastigo.NewConn(),
		Index:        index,
		VersionIndex: index + "_versions",
	}

	ed.Conn.Domain = esshost
	ed.Conn.Port = fmt.Sprintf("%d", essport)

	if !ed.IndexExists() {
		if len(mappingfile) > 0 {
			log.V(9).Infof("Initializing with mapping file: %#v\n", mappingfile[0])
			return &ed, ed.initializeIndex(mappingfile[0])
		} else {
			return &ed, ed.initializeIndex("")
		}
	}
	return &ed, nil
}

func (e *ElasticsearchDatastore) IndexExists() bool {
	_, err := e.Conn.DoCommand("GET", "/"+e.Index, nil, nil)
	if err != nil {
		return false
	}
	return true
}

/* Used to determine if the mapping file can be applied with the given version */
func (e *ElasticsearchDatastore) IsVersionSupported() (supported bool) {
	supported = false

	info, err := e.Info()
	if err != nil {
		log.V(0).Infof("Could not get version: %s\n", err)
		return
	}

	versionStr := strings.Join(strings.Split(info.Version.Number, ".")[:2], ".")
	verNum, err := strconv.ParseFloat(versionStr, 64)
	if err != nil {
		log.V(0).Infof("Could not get version: %s\n", err)
		return
	}

	if verNum >= 1.4 {
		supported = true
	}
	return
}

/* Elasticsearch instance information.  e.g. version */
func (e *ElasticsearchDatastore) Info() (info EssInfo, err error) {
	var b []byte
	b, err = e.Conn.DoCommand("GET", "", nil, nil)
	err = json.Unmarshal(b, &info)
	return
}

func (e *ElasticsearchDatastore) Close() {
	e.Conn.Close()
}

func (e *ElasticsearchDatastore) applyMappingFile(mapfile string) (err error) {
	if !e.IsVersionSupported() {
		err = fmt.Errorf("Not creating mapping. ESS version not supported. Must be > 1.4.")
		return
	}

	if _, err = os.Stat(mapfile); err != nil {
		err = fmt.Errorf("Not creating mapping. Mapping file not found (%s): %s", mapfile, err)
		return
	}
	// Read mapping file to get map name
	var mdb []byte
	mdb, err = ioutil.ReadFile(mapfile)
	if err != nil {
		return
	}

	var mapData map[string]interface{}
	if err = json.Unmarshal(mdb, &mapData); err != nil {
		return
	}
	// Get map name from first key
	var (
		normMap  = map[string]interface{}{}
		mapname  string
		mapbytes []byte
	)
	for k, _ := range mapData {
		normMap[k] = mapData[k]
		mapname = k
		break
	}

	if mapbytes, err = json.Marshal(normMap); err != nil {
		return
	}
	log.V(10).Infof("Mapping (%s): %s\n", mapname, mapbytes)

	if err = e.Conn.PutMappingFromJSON(e.Index, mapname, mapbytes); err != nil {
		return
	} else {
		log.Infof("Updated '%s' mapping for index '%s'\n", mapname, e.Index)
	}
	// Versioning index
	if err = e.Conn.PutMappingFromJSON(e.VersionIndex, mapname, mapbytes); err != nil {
		return
	} else {
		log.Infof("Updated '%s' mapping for index '%s'\n", mapname, e.VersionIndex)
	}

	return
}

func (e *ElasticsearchDatastore) initializeIndex(mappingFile string) error {
	resp, err := e.Conn.CreateIndex(e.Index)
	if err != nil {
		return err
	}
	log.V(3).Infof("Index created: %s %s\n", e.Index, resp)
	// Versioning index
	resp, err = e.Conn.CreateIndex(e.VersionIndex)
	if err != nil {
		return err
	}
	log.V(3).Infof("Version index created: %s %s\n", e.Index, resp)

	if len(mappingFile) > 1 {
		log.V(6).Infof("Applying mapping file: %s\n", mappingFile)
		return e.applyMappingFile(mappingFile)
	}

	return nil
}
