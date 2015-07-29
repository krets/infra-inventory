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


Running the Server
------------------

Make sure elasticsearch is runnin or start it.

    - /etc/init.d/elasticsearch status
    - /etc/init.d/elasticsearch start

Run infra-inventory

    - $GOPATH/bin/infra-inventory -logtostderr -v 10 -c $GOPATH/src/github.com/euforia/infra-inventory/etc/infra-inventory.json


Usage
-----
The following query parameters are available depending on the endpoint:

    - version   int
    - sortby    <attr:[asc|desc]>
    - attrs     []

### Endpoints
The following verbs and endpoints are available:

List asset types:

    - GET /v1/

Get asset:

    - GET /v1/<asset_type>/<asset_id>

Get a specific asset version:

    - GET /v1/<asset_type>/<asset_id>?version=<version>

Get last 10 asset versions:

    - GET /v1/<asset_type>/<asset_id>/versions

Add asset:

    - POST /v1/<asset_type>/<asset_id>

        {
            "name": "foo.bar.com",
            "status": "running"
            ...
        }

Edit asset:

    - PUT /v1/<asset_type>/<asset_id>

        {
            "status": "stopped"
            ...
        }

Delete asset:

    - DELETE /v1/<asset_type>/<asset_id>


Search for an asset of type `asset_type` that matches both attributes:

    - GET /v1/<asset_type>

        {
            "status": "stopped",
            "os": "ubuntu"
        }


