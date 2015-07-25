package inventory

import (
	"encoding/json"
	//"fmt"
	"github.com/euforia/ess-go-wrapper"
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

type Inventory struct {
	datastore *esswrapper.EssWrapper
}

func NewInventory(datastore *esswrapper.EssWrapper) (ir *Inventory) {
	ir = &Inventory{datastore}
	return
}

/* Normalize asset type input from user */
func (ir *Inventory) normalizeAssetType(assetType string) string {
	return strings.ToLower(assetType)
}

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
	defer r.Body.Close()
	log.V(15).Infof("%s\n", b)

	if err != nil {
		code = 500
		data = []byte(err.Error())
		headers = map[string]string{"Content-Type": "text/plain"}
		return
	}

	if len(b) < 1 {
		code = 400
		data = []byte(`No payload`)
		headers = map[string]string{"Content-Type": "text/plain"}
		return
	}

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
	Handler for all method to endpoint: /<asset_type>/<asset>
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

/*
	TODO: Needs implementation

	Handle requests searching within an asset type
*/
func (ir *Inventory) AssetTypeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		reqVars = mux.Vars(r)
		//assetType = ir.normalizeAssetType(reqVars["asset_type"])
		code    = 200
		headers = map[string]string{"Content-Type": "application/json"}
		data    = []byte(`{"status":"To be implemented!"}`)
	)
	log.V(15).Infof("%#v\n", reqVars)

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	log.V(12).Infof("Body: %s %d\n", body, len(body))

	paramsQuery := r.URL.Query()

	if len(paramsQuery) < 1 && len(body) < 1 {
		code = 404
		data = []byte("Not found!")
		headers["Content-Type"] = "text/plain"
	} else {
		var bodyQuery map[string]string
		err := json.Unmarshal(body, &bodyQuery)
		if err == nil {
			// Process body
		}
		//ir.datastore.GetBy()

		log.V(12).Infof("%#v\n", paramsQuery)
	}

	WriteAndLogResponse(w, r, code, headers, data)
}

func (ir *Inventory) ListAssetTypeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		types []string
		err   error
		b     []byte
	)

	if types, err = ir.datastore.GetTypes(); err != nil {
		WriteAndLogResponse(w, r, 500, map[string]string{"Content-Type": "text/plain"},
			[]byte(err.Error()))
		return
	}

	if b, err = json.Marshal(types); err != nil {
		WriteAndLogResponse(w, r, 500, map[string]string{"Content-Type": "text/plain"},
			[]byte(err.Error()))
		return
	}

	WriteAndLogResponse(w, r, 200, map[string]string{"Content-Type": "application/json"}, b)
}
