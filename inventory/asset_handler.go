package inventory

import (
	"encoding/json"
	//"fmt"
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

/*
   Handle getting assets GET /<asset_type>/<asset>
*/
func (ir *Inventory) assetGetHandler(assetType, assetId string) (code int, headers map[string]string, data []byte) {
	ans, err := ir.datastore.GetAsset(assetType, assetId)
	if err != nil {
		code = 404
		headers = map[string]string{"Content-Type": "text/plain"}
		data = []byte(err.Error())
	} else {
		if data, err = json.Marshal(ans); err != nil {
			code = 500
			data = []byte(err.Error())
			headers = map[string]string{"Content-Type": "text/plain"}
		} else {
			code = 200
			headers = map[string]string{"Content-Type": "application/json"}
		}
	}
	return
}

/*
   Handle adding assets POST /<asset_type>/<asset>
   Handle editing assets PUT /<asset_type>/<asset>
*/
func (ir *Inventory) assetPostPutHandler(assetType, assetId string, r *http.Request) (code int, headers map[string]string, data []byte) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		code = 400
		data = []byte(err.Error())
		headers = map[string]string{"Content-Type": "text/plain"}
		return
	}
	defer r.Body.Close()

	var id string
	switch r.Method {
	case "POST":
		id, err = ir.datastore.CreateAsset(assetType, assetId, b)
		break
	case "PUT":
		var bmap map[string]interface{}
		err = json.Unmarshal(b, &bmap)
		if err != nil {
			break
		}
		id, err = ir.datastore.EditAsset(assetType, assetId, bmap)
		break
	}

	if err != nil {
		code = 404
		headers = map[string]string{"Content-Type": "text/plain"}
		data = []byte(err.Error())
	} else {
		code = 200
		headers = map[string]string{"Content-Type": "application/json"}
		data = []byte(`{"id": "` + id + `"}`)
	}
	return
}

/*
   Handle getting assets by version GET /<asset_type>/<asset>?version=<version>
*/
func (ir *Inventory) assetGetVersionHandler(assetType, assetId, versionStr string) (code int, headers map[string]string, data []byte) {
	var version, err = strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		code = 404
		data = []byte(err.Error())
		headers = map[string]string{"Content-Type": "text/plain"}
	} else {
		asset, err := ir.datastore.GetAssetVersion(assetType, assetId, version)
		if err != nil {
			code = 404
			data = []byte(err.Error())
			headers = map[string]string{"Content-Type": "text/plain"}
		} else {
			code = 200
			data, _ = json.Marshal(asset)
			headers = map[string]string{"Content-Type": "application/json"}
		}
	}
	return
}

/*
   Handler for all methods to endpoint: /<asset_type>/<asset>
*/
func (ir *Inventory) AssetHandler(w http.ResponseWriter, r *http.Request) {
	var (
		headers = map[string]string{}
		code    int
		data    = make([]byte, 0)

		restVars = mux.Vars(r)

		assetType = ir.normalizeAssetType(restVars["asset_type"])
		assetId   = restVars["asset"]
	)
	log.V(15).Infof("%#v\n", restVars)

	switch r.Method {
	case "GET":
		queryParams := r.URL.Query()
		if versionArr, ok := queryParams["version"]; ok {
			code, headers, data = ir.assetGetVersionHandler(assetType, assetId, versionArr[0])
		} else {
			code, headers, data = ir.assetGetHandler(assetType, assetId)
		}
		break
	case "POST", "PUT":
		code, headers, data = ir.assetPostPutHandler(assetType, assetId, r)
		break
	case "DELETE":
		if ir.datastore.RemoveAsset(assetType, assetId) {
			code = 200
		} else {
			code = 500
		}
		break
	}

	WriteAndLogResponse(w, r, code, headers, data)
}

/*
   Handle getting asset versions GET /<asset_type>/<asset>/versions
*/
func (ir *Inventory) AssetVersionsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		headers = map[string]string{}
		code    int
		data    = make([]byte, 0)

		restVars = mux.Vars(r)

		assetType = ir.normalizeAssetType(restVars["asset_type"])
		assetId   = restVars["asset"]
	)
	log.V(15).Infof("%#v\n", restVars)

	// the count should come from a query param
	assetVersions, err := ir.datastore.GetAssetVersions(assetType, assetId, 10)
	if err != nil {
		code = 404
		data = []byte(err.Error())
		headers["Content-Type"] = "text/plain"
	} else {
		code = 200
		data, _ = json.Marshal(assetVersions.Hits.Hits)
		headers["Content-Type"] = "application/json"
	}

	WriteAndLogResponse(w, r, code, headers, data)
}
