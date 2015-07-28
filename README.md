Inventory
=========
Inventory system


Requirements
------------

go >= 1.4.2

Installation
------------

go get github.com/euforia/infra-inventory

cp $GOPATH/src/github.com/euforia/infra-inventory/etc/infra-inventory.json{.sample,}

$GOPATH/bin/infra-inventory -logtostderr -v 10 -c $GOPATH/src/github.com/euforia/infra-inventory/etc/infra-inventory.json


Usage
-----

List asset types

    - GET /v1/

Get asset 

    - GET /v1/<asset_type>/<asset_id>

Add asset

    - POST /v1/<asset_type>/<asset_id>

        {
            "name": "foo.bar.com",
            "status": "running"
            ...
        }

Edit asset

    - PUT /v1/<asset_type>/<asset_id>

        {
            "status": "stopped"
            ...
        }

