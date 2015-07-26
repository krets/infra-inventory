package inventory

import (
	"encoding/json"
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

/*
   Handle getting assets GET /<asset_type>/<asset>
*/
func (ir *Inventory) assetGetHandler(assetType, assetId string) (code int, headers map[string]string, data []byte) {
	ans, err := ir.datastore.Get(assetType, assetId)
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
		id, err = ir.datastore.AddWithId(assetType, assetId, b)
		break
	case "PUT":
		id, err = ir.datastore.Update(assetType, assetId, b)
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
   Handler for all methods to endpoint: /<asset_type>/<asset>
*/
func (ir *Inventory) AssetHandler(w http.ResponseWriter, r *http.Request) {
	var (
		headers = map[string]string{}
		code    int
		data    = make([]byte, 0)

		reqVars   = mux.Vars(r)
		assetType = ir.normalizeAssetType(reqVars["asset_type"])
		assetId   = reqVars["asset"]
	)
	log.V(15).Infof("%#v\n", reqVars)

	switch r.Method {
	case "GET":
		code, headers, data = ir.assetGetHandler(assetType, assetId)
		break
	case "POST", "PUT":
		code, headers, data = ir.assetPostPutHandler(assetType, assetId, r)
		break
	case "DELETE":
		if ir.datastore.Delete(assetType, assetId) {
			code = 200
		} else {
			code = 500
		}
		break
	}

	WriteAndLogResponse(w, r, code, headers, data)
}
