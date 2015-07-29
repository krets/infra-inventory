package main

import (
	"flag"
	//"fmt"
	//"github.com/euforia/ess-go-wrapper"
	"github.com/euforia/infra-inventory/inventory"
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	listenAddr = flag.String("l", ":5454", "Address to start HTTP Server on")
	enableAuth = flag.Bool("enable-auth", false, "Enable auth on write requests")

	configFile = flag.String("c", "infra-inventory.json", "Config file")
)

func bootstrapServer(cfg *inventory.InventoryConfig) {
	// Instantiate datastore
	dstore, err := inventory.NewElasticsearchDatastore(cfg.Datastore.Config.Host, cfg.Datastore.Config.Port,
		cfg.Datastore.Config.Index, cfg.Datastore.Config.MappingFile)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	invDs := inventory.NewInventoryDatastore(dstore)
	// New inventory instance
	inv := inventory.NewInventory(cfg, invDs)
	// Register http endpoints with muxer
	rtr := mux.NewRouter()
	rtr.HandleFunc(cfg.Endpoints.Prefix+"/", inv.ListAssetTypesHandler).Methods("GET")

	rtr.HandleFunc(cfg.Endpoints.Prefix+"/{asset_type}", inv.AssetTypeHandler).Methods("GET")

	if *enableAuth {
		// Setup handler with all pre-processors
		log.Infof("Auth enabled!\n")

		rtr.HandleFunc(cfg.Endpoints.Prefix+"/{asset_type}/{asset}",
			inventory.AuthHandler(inv.AssetHandler)).Methods("GET", "POST", "PUT", "DELETE")
	} else {
		// Setup handler with all pre-processors except auth
		log.Infof("Auth disabled!\n")

		rtr.HandleFunc(cfg.Endpoints.Prefix+"/{asset_type}/{asset}",
			inv.AssetHandler).Methods("GET", "POST", "PUT", "DELETE")
	}

	http.Handle("/", rtr)

	log.Infof("Elasticsearch (%s): %s:%d/%s\n", cfg.Datastore.Config.Index, cfg.Datastore.Config.Host,
		cfg.Datastore.Config.Port, cfg.Datastore.Config.Index)
	log.Infof("Elasticsearch (%s): %s:%d/%s\n", dstore.VersionIndex, cfg.Datastore.Config.Host,
		cfg.Datastore.Config.Port, dstore.VersionIndex)
	log.Infof("Starting server on %s%s\n", *listenAddr, cfg.Endpoints.Prefix)
}

func loadConfig() *inventory.InventoryConfig {
	flag.Parse()

	cfg, err := inventory.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	if *enableAuth {
		cfg.Auth.Enabled = true
	}
	return cfg
}

func main() {
	cfg := loadConfig()

	bootstrapServer(cfg)

	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		log.Fatalf("%s\n", err)
	}
}
