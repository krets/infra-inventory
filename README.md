Inventory
=========
Inventory system

Usage
-----

List asset types

    - GET /v1/

Get asset 

    - GET /v1/<asset_type>/<asset_id>

Add asset

    - POST /v1/<asset_type>/<asset_id>

        {
            "name": "foo.bar.com"
            "status": "running"
        }

Edit asset

    - PUT /v1/<asset_type>/<asset_id>

        {
            "status": "stopped"
        }

