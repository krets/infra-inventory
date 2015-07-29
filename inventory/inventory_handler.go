package inventory

import (
	"encoding/json"
	"fmt"
	//"github.com/euforia/ess-go-wrapper"
	log "github.com/golang/glog"
	//"github.com/gorilla/mux"
	elastigo "github.com/mattbaird/elastigo/lib"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type AssetResponse struct {
	Id   string    `json:"id"`
	Type string    `json:"type"`
	Data AssetData `json:"data"`
}

type Inventory struct {
	//datastore *esswrapper.EssWrapper
	datastore IDatastore
	cfg       *InventoryConfig
}

func NewInventory(cfg *InventoryConfig, datastore IDatastore) (ir *Inventory) {
	ir = &Inventory{datastore, cfg}
	return
}

/* Normalize asset type input from user */
func (ir *Inventory) normalizeAssetType(assetType string) string {
	return strings.ToLower(assetType)
}

/*
	Returns:
		should also return the params as elastic search globale args/opts
*/
func (ir *Inventory) parseRequestQueryParams(r *http.Request) (err error) {

	paramsQuery := r.URL.Query()
	log.V(12).Infof("%#v\n", paramsQuery)

	// Parse global query opts.
	if vals, ok := paramsQuery["sortby"]; ok {

		sorting := make([]*elastigo.SortDsl, len(vals))

		for i, v := range vals {
			sarr := strings.Split(v, ":")
			if len(sarr) != 2 {
				err = fmt.Errorf("Invalid request: sortby=%s", v)
				return
			}
			var sDsl *elastigo.SortDsl

			switch sarr[1] {
			case "asc":
				sDsl = elastigo.Sort(sarr[0]).Asc()
				break
			case "dsc":
				sDsl = elastigo.Sort(sarr[0]).Desc()
				break
			default:
				err = fmt.Errorf("Invalid sort argument: %s", sarr[1])
				return
			}
			sorting[i] = sDsl
		}
		b, _ := json.Marshal(sorting)
		log.V(12).Infof("Query (sorting): %s\n", b)
	}
	return
}

/*
{
	"type": ["virtualserver", "physicalserver"],
	"os": "ubuntu"
}
*/
func (ir *Inventory) parseRequestBody(r *http.Request) (query interface{}, err error) {

	// check happens earlier
	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}
	defer r.Body.Close()

	var req map[string]interface{}
	if err = json.Unmarshal(body, &req); err != nil {
		return
	}

	filterOps := []interface{}{}

	for k, v := range req {
		switch v.(type) {
		case string:
			val, _ := v.(string)
			val = strings.TrimSpace(val)
			if strings.HasPrefix(val, ">") || strings.HasPrefix(val, "<") {
				// Parse number
				aVal := ""
				if strings.HasPrefix(val, ">=") || strings.HasPrefix(val, "<=") {
					aVal = strings.TrimSpace(val[2:])
				} else {
					aVal = strings.TrimSpace(val[1:])
				}
				// Parse number for comparison
				var nVal interface{}
				nVal, ierr := strconv.ParseInt(aVal, 10, 64)
				if ierr != nil {
					ierr = nil
					if nVal, ierr = strconv.ParseFloat(aVal, 64); ierr != nil {
						err = ierr
						return
					}
				}
				// Add range filterop
				if strings.HasPrefix(val, ">") {
					filterOps = append(filterOps, elastigo.Range().Field(k).Gt(nVal))
				} else {
					filterOps = append(filterOps, elastigo.Range().Field(k).Lt(nVal))
				}

			} else {
				filterOps = append(filterOps, elastigo.Filter().Terms(k, val))
			}
			break
		case int:
			//val, _ := v.(int)
			break
		case int64:
			//val, _ := v.(int64)
			break
		case float64:
			//val, _ := v.(float64)
			break
		case []interface{}:
			vals, _ := v.([]interface{})
			filterOps = append(filterOps, elastigo.Filter().Terms(k, vals...))
			break
		case interface{}:
			//val, _ := v.(interface{})
			break
		default:
			err = fmt.Errorf("invalid type: %#v", v)
			return
		}
	}

	query = elastigo.Search(ir.cfg.Datastore.Config.Index).Filter(filterOps...)

	return
}

func (ir *Inventory) executeSearchQuery(assetType string, r *http.Request) (rslt elastigo.SearchResult, err error) {
	var q interface{}

	// IN PROGRESS
	ir.parseRequestQueryParams(r)

	if q, err = ir.parseRequestBody(r); err != nil {
		return
	}

	b, _ := json.MarshalIndent(q, " ", "  ")
	log.V(15).Infof("%s ==> %s\n", r.RequestURI, b)

	rslt, err = ir.datastore.Search(assetType, q)
	return
}
