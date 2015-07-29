package inventory

import (
	"encoding/json"
	"fmt"
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (ir *Inventory) checkWriteRequest(r *http.Request) (data map[string]interface{}, err error) {
	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	var bmap map[string]interface{}
	if err = json.Unmarshal(body, &bmap); err != nil {
		return
	}

	data = map[string]interface{}{}
	for k, v := range bmap {

		updated := false
		for _, rField := range ir.cfg.AssetCfg.RequiredFields {

			if strings.EqualFold(k, rField) {
				val, ok := v.(string)
				if !ok {
					err = fmt.Errorf("'%s' field must be a string!\n", rField)
					return
				}
				val = strings.TrimSpace(val)
				if len(val) < 1 {
					err = fmt.Errorf("'%s' field value required!\n", rField)
					return
				} else {
					data[rField] = val
					updated = true
				}
				break
			}
		}
		if !updated {
			data[k] = v
		}
		log.V(12).Infof("%#v\n", data)
	}

	if r.Method == "POST" {
		for _, v := range ir.cfg.AssetCfg.RequiredFields {
			if _, ok := data[v]; !ok {
				err = fmt.Errorf("'%s' field required!\n", v)
				return
			}
		}
	}
	return
}

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
		rsp, err := AssembleResponseFromBaseResponse(ans)
		if err != nil {
			code = 400
			data = []byte(err.Error())
			headers = map[string]string{"Content-Type": "text/plain"}
			return
		}

		if data, err = json.Marshal(rsp); err != nil {
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
	var (
		reqData map[string]interface{}
		err     error
		id      string
	)
	if reqData, err = ir.checkWriteRequest(r); err != nil {
		code = 400
		data = []byte(err.Error())
		headers = map[string]string{"Content-Type": "text/plain"}
		return
	}

	switch r.Method {
	case "POST":
		id, err = ir.datastore.CreateAsset(assetType, assetId, reqData, true)
		break
	case "PUT":
		id, err = ir.datastore.EditAsset(assetType, assetId, reqData)
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
			rsp, err := AssembleResponseFromBaseResponse(asset)
			if err != nil {
				code = 400
				data = []byte(err.Error())
				headers = map[string]string{"Content-Type": "text/plain"}
				return
			}
			code = 200
			data, _ = json.Marshal(rsp)
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

		restVars  = mux.Vars(r)
		assetType = ir.normalizeAssetType(restVars["asset_type"])
		assetId   = restVars["asset"]
	)

	// the count should come from a query param
	assetVersions, err := ir.datastore.GetAssetVersions(assetType, assetId, 10)
	if err != nil {
		code = 404
		data = []byte(err.Error())
		headers["Content-Type"] = "text/plain"
	} else {
		log.V(11).Infof("Found versions: %d\n", assetVersions.Hits.Len())

		if _, ok := r.URL.Query()["diff"]; ok {
			// Generates diffs for versions
			maplist := make([]map[string]interface{}, assetVersions.Hits.Len())
			for i, ver := range assetVersions.Hits.Hits {
				var m map[string]interface{}
				if err := json.Unmarshal(*ver.Source, &m); err != nil {
					log.Errorf("%s\n", err)
				}
				maplist[i] = m
			}

			diffs, err := GenerateVersionDiffs(maplist...)
			if err != nil {
				data = []byte(err.Error())
				code = 400
				headers["Content-Type"] = "text/plain"
			} else {
				code = 200
				data, _ = json.Marshal(diffs)
				headers["Content-Type"] = "application/json"
			}
		} else {
			// Return full versions
			code = 200
			rsp, _ := AssembleResponseFromHits(assetVersions.Hits.Hits)
			data, _ = json.Marshal(rsp)
			headers["Content-Type"] = "application/json"
		}
	}

	WriteAndLogResponse(w, r, code, headers, data)
}
