Inventory
=========
Inventory system


Requirements
------------

    go >= 1.4.2
    elasticsearch >= 1.4.x


Installation
------------

    # Build and install
    $ go get github.com/euforia/infra-inventory
    
    # Copy over the sample config
    $ cp $GOPATH/src/github.com/euforia/infra-inventory/etc/infra-inventory.json{.sample,}


Running the Server
------------------
Once the installation is complete, you'll need to make sure elasticsearch is running. To do so run the following as necessary:

    # Check if it running
    $ /etc/init.d/elasticsearch status

    # Start it (optional)
    $ /etc/init.d/elasticsearch start

Once elasticsearch is running, execute the following command to start the service:

    $ cd $GOPATH
    $ ./bin/infra-inventory -logtostderr -v 10 -c ./src/github.com/euforia/infra-inventory/etc/infra-inventory.json.sample

You can now start using the inventory system.

Asset Type
----------
Every asset has a asset type.  To avoid the creation of unwanted types, only users in the `admin` group of the `LocalAuthGroups` are allowed to create asset types.  The type will automatically be created upon the creation of the asset if done by an authorized `admin` user.  More information about the `LocalAuthGroups` can be found below.


Versions
--------
Versions are automatically created on each write.  When a write occurs, the existing asset is copied over to the `versions` index incrementing the version number then performing the write.  It is possible get a list of versions or specific versions of a given asset.  A version to version diff can also be obtained.


Asset
-----
Each asset must have an associated type.  Before an asset can be created, the asset type must be created.  As mentioned before, only admins are allowed to create new asset types. An asset has versions available.  These are only available after the first write operation.  A `current` asset has no version.  


### Endpoints
The following verbs and endpoints are available:

List available asset types:

    - GET /v1/

Response e.g.:
    
    ["virtualserver", "dnsrecord"]

Get asset:

    - GET /v1/<asset_type>/<asset_id>

Response e.g.:

    {
        "id": "foo.bar.org"
        "type":"virtualserver",
        "data":{
            "status":"running",
            "environment": "dev",
            "created_by": "user1",
            "updated_by": "user2"
            ....
        }
    }

Get a specific asset version:

    - GET /v1/<asset_type>/<asset_id>?version=<version>

Response e.g.:

    {
        "id": "<asset_id>"
        "type":"<asset_type>",
        "data":{
            "status":"running",
            "environment": "dev",
            "version": <version>,
            "created_by": "user1",
            "updated_by": "user2"
            ....
        }
    }

Get last 10 asset versions:

    - GET /v1/<asset_type>/<asset_id>/versions

Get a version to version incremental diff:

    - GET /v1/<asset_type>/<asset_id>/versions?diff

Response e.g.:

    [{
        version: 2,
        against_version: 1,
        diff: "<diff_data>"
    },{
        ....
    }]

Create a new asset:

    - POST /v1/<asset_type>/<asset_id>

        {
            "name": "foo.bar.com",
            "status": "running",
            "environment": "development"
            ...
        }

Response e.g.:

    { "id": "<asset_id>" }
    
Edit an existing asset:

    - PUT /v1/<asset_type>/<asset_id>

        {
            "status": "stopped"
            ...
        }

Response e.g.:

    { "id": "<asset_id>" }

Delete an asset:

    - DELETE /v1/<asset_type>/<asset_id>


Search for an asset of type `asset_type` that matches both attributes:

    - GET /v1/<asset_type>

        {
            "status": "stopped",
            "os": "ubuntu"
        }


Local Auth Groups
-----------------
Local auth groups are primarily used to create asset types.  The configuration file can be found at etc/local-groups.json. Fill in the usernames you wish to allow.  The user must match that used for 'HTTP Basic Auth'.

Example:

    {
        "admin": ["user1", "user2"]
    }


