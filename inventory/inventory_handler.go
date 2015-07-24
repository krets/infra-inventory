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

func (ir *Inventory) normalizeAssetType(assetType string) string {
	return strings.ToLower(assetType)
}

func (ir *Inventory) AssetHandler(w http.ResponseWriter, r *http.Request) {
	var (
		reqVars   = mux.Vars(r)
		assetType = ir.normalizeAssetType(reqVars["asset_type"])
		headers   = map[string]string{}
		code      int
		data      = make([]byte, 0)
	)
	log.V(15).Infof("%#v\n", reqVars)

	switch r.Method {
	case "GET":
		ans, err := ir.datastore.Get(assetType, reqVars["asset"])
		if err != nil {
			code = 404
			headers["Content-Type"] = "text/plain"
			data = []byte(err.Error())
		} else {
			if data, err = json.Marshal(ans); err != nil {
				code = 500
				data = []byte(err.Error())
				headers["Content-Type"] = "text/plain"
			} else {
				code = 200
				headers["Content-Type"] = "application/json"
			}
		}
		break
	case "POST":
		// Add
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			code = 500
			headers["Content-Type"] = "text/plain"
			data = []byte(err.Error())
			break
		}
		// TODO: check the data before adding
		id, err := ir.datastore.AddWithId(assetType, reqVars["asset"], b)
		if err != nil {
			code = 404
			headers["Content-Type"] = "text/plain"
			data = []byte(err.Error())
		} else {
			code = 200
			headers["Content-Type"] = "application/json"
			data = []byte(`{"id": "` + id + `"}`)
		}
		break
	case "PUT":
		// Update
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			code = 500
			headers["Content-Type"] = "text/plain"
			data = []byte(err.Error())
			break
		}
		// TODO: check the data before updating
		if err = ir.datastore.Update(assetType, reqVars["asset"], b); err != nil {
			code = 404
			headers["Content-Type"] = "text/plain"
			data = []byte(err.Error())
		} else {
			code = 200
		}
		break
	case "DELETE":
		if ir.datastore.Delete(assetType, reqVars["asset"]) {
			code = 200
		} else {
			code = 500
		}
		break
	}

	ir.writeResponse(w, r, code, headers, data)
}

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

	ir.writeResponse(w, r, code, headers, data)
}

func (ir *Inventory) ListAssetTypeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		types []string
		err   error
		b     []byte
	)

	if types, err = ir.datastore.GetTypes(); err != nil {
		ir.writeResponse(w, r, 500, map[string]string{"Content-Type": "text/plain"},
			[]byte(err.Error()))
		return
	}

	if b, err = json.Marshal(types); err != nil {
		ir.writeResponse(w, r, 500, map[string]string{"Content-Type": "text/plain"},
			[]byte(err.Error()))
		return
	}

	ir.writeResponse(w, r, 200, map[string]string{"Content-Type": "application/json"}, b)
}

func (ir *Inventory) writeResponse(w http.ResponseWriter, r *http.Request, code int, headers map[string]string, data []byte) {
	w.WriteHeader(code)

	if headers != nil {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}

	w.Write(data)
	log.Infof("%s %d %s %d\n", r.Method, code, r.RequestURI, len(data))
}
