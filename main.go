package main

import (
	"flag"
	"github.com/euforia/infra-inventory/inventory"
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	listenAddr = flag.String("l", ":5454", "Address to start HTTP Server on")
	enableAuth = flag.Bool("enable-auth", false, "Enable auth on write requests")

	configFile = flag.String("c", "infra-inventory.json", "Config file")
	// global config
	cfg *inventory.InventoryConfig
)

func loadConfig() {
	flag.Parse()

	var err error
	if cfg, err = inventory.LoadConfig(*configFile); err != nil {
		log.Fatalf("%s\n", err)
	}

	if *enableAuth {
		cfg.Auth.Enabled = true
	}
	if cfg.Auth.Enabled {
		log.Infof("Auth enabled!\n")
	} else {
		log.Infof("Auth disabled!\n")
	}
}

func startServer(inv *inventory.Inventory) {
	// Register http endpoints with muxer
	rtr := mux.NewRouter()
	rtr.HandleFunc(cfg.Endpoints.Prefix+"/",
		inv.ListAssetTypesHandler).Methods("GET")

	rtr.HandleFunc(cfg.Endpoints.Prefix+"/{asset_type}",
		inv.AssetTypeHandler).Methods("GET")

	rtr.HandleFunc(cfg.Endpoints.Prefix+"/{asset_type}/{asset}",
		inv.AssetHandler).Methods("GET", "POST", "PUT", "DELETE")

	rtr.HandleFunc(cfg.Endpoints.Prefix+"/{asset_type}/{asset}/versions",
		inv.AssetVersionsHandler).Methods("GET")

	http.Handle("/", rtr)

	log.Infof("Starting server on %s%s\n", *listenAddr, cfg.Endpoints.Prefix)
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatalf("%s\n", err)
	}
}

func initializeInventory() (inv *inventory.Inventory) {
	var (
		dstore    *inventory.ElasticsearchDatastore
		invDstore *inventory.InventoryDatastore
		err       error
	)

	if dstore, err = inventory.NewElasticsearchDatastore(
		cfg.Datastore.Config.Host, cfg.Datastore.Config.Port,
		cfg.Datastore.Config.Index, cfg.Datastore.Config.MappingFile); err != nil {

		log.Fatalf("%s\n", err)
	}

	log.Infof("Elasticsearch (%s): %s:%d/%s\n", cfg.Datastore.Config.Index, cfg.Datastore.Config.Host,
		cfg.Datastore.Config.Port, cfg.Datastore.Config.Index)
	log.Infof("Elasticsearch (%s): %s:%d/%s\n", dstore.VersionIndex, cfg.Datastore.Config.Host,
		cfg.Datastore.Config.Port, dstore.VersionIndex)

	// inventory datastore
	invDstore = inventory.NewInventoryDatastore(dstore)
	// New inventory instance (api etc.)
	if inv, err = inventory.NewInventory(cfg, invDstore); err != nil {
		log.Fatalf("%s\n", err)
	}
	return
}

func main() {
	loadConfig()

	inv := initializeInventory()
	startServer(inv)
}
