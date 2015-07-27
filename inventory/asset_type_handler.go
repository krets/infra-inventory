package inventory

import (
	"encoding/json"
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
)

/*
   Handle requests searching within an asset type i.e GET /<asset_type>
*/
func (ir *Inventory) AssetTypeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		reqVars   = mux.Vars(r)
		assetType = ir.normalizeAssetType(reqVars["asset_type"])
		code      = 200
		headers   = map[string]string{"Content-Type": "application/json"}
		data      []byte
	)
	log.V(15).Infof("%#v\n", reqVars)

	essResp, err := ir.executeSearchQuery(assetType, r)
	if err != nil {
		data = []byte(err.Error())
		code = 400
		headers["Content-Type"] = "text/plain"
	} else {
		data, _ = json.Marshal(essResp.Hits.Hits)
	}

	WriteAndLogResponse(w, r, code, headers, data)
}

func (ir *Inventory) ListAssetTypesHandler(w http.ResponseWriter, r *http.Request) {
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
