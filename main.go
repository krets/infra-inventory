package main

import (
	"flag"
	//"fmt"
	"github.com/euforia/ess-go-wrapper"
	"github.com/euforia/infra-inventory/inventory"
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	listenAddr  = flag.String("l", ":5454", "Address to start HTTP Server on")
	mappingFile = flag.String("m", "mapping.json", "Mapping file to use")
	indexName   = flag.String("index", "inventory", "Name of the index")
	essHost     = flag.String("ess-host", "localhost", "Elasticsearch host")
	essPort     = flag.Int("ess-port", 9200, "Elasticsearch port")
)

func getInventoryInstance() *inventory.Inventory {
	dstore, err := esswrapper.NewEssWrapper(*essHost, *essPort, *indexName, *mappingFile)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	return inventory.NewInventory(dstore)
}

func bootstrapServer() {
	inv := getInventoryInstance()

	rtr := mux.NewRouter()
	rtr.HandleFunc("/v1/{asset_type}/{asset}", inv.AssetHandler).Methods("GET", "POST", "PUT", "DELETE")
	rtr.HandleFunc("/v1/{asset_type}", inv.AssetTypeHandler).Methods("GET")
	rtr.HandleFunc("/v1/", inv.ListAssetTypeHandler).Methods("GET")

	http.Handle("/", rtr)
}

func main() {
	flag.Parse()

	bootstrapServer()

	log.Infof("Elasticsearch: %s:%d/%s\n", *essHost, *essPort, *indexName)
	log.Infof("Starting server on %s\n", *listenAddr)

	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatalf("%s\n", err)
	}
}
