package inventory

import (
	"encoding/json"
	"testing"
)

var (
	testAssetType = "test_asset_type"
	testData      = map[string]string{
		"name": "test",
		"host": "test.foo.bar",
	}
	testData2 = map[string]string{
		"name": "test2",
		"host": "test.foo.bar",
	}
	testUpdateData = map[string]string{
		"host": "test.foo.bar.updated",
	}
	testEds, _ = NewElasticsearchDatastore(testEssHost, testEssPort, testIndex, testMappingFile)
	testIds    = NewInventoryDatastore(testEds)
)

func Test_InventoryDatastore_CreateAsset(t *testing.T) {

	id, err := testIds.CreateAsset(testAssetType, testData["name"], testData, true)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%s", id)
}

func Test_InventoryDatastore_GetAsset(t *testing.T) {

	asset, err := testIds.GetAsset(testAssetType, testData["name"])
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%#v", asset)
}

func Test_InventoryDatastore_EditAsset(t *testing.T) {

	id, err := testIds.EditAsset(testAssetType, testData["name"], testUpdateData)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%s", id)

	asset, _ := testIds.GetAsset(testAssetType, testData["name"])

	var d map[string]interface{}
	json.Unmarshal(*asset.Source, &d)
	if _, ok := d["name"]; !ok {
		t.Fatalf("Overwrote exising object")
	}
	//testIds.Close()
}

func Test_InventoryDatastore_ListAssetTypes(t *testing.T) {
	types, err := testIds.ListAssetTypes()
	if err != nil {
		t.Fatalf("%s", err)
	}

	t.Logf("%#v", types)
}

func Test_InventoryDatastore_RemoveAsset(t *testing.T) {

	if !testIds.RemoveAsset(testAssetType, testData["name"]) {
		t.Fatalf("Failed to remove asset")
	}
	_, err := testIds.GetAsset(testAssetType, testData["name"])
	if err == nil {
		t.Fatalf("Did not remove asset")
	}
	testIds.Conn.DeleteIndex(testIds.Index)
	testIds.Conn.DeleteIndex(testIds.VersionIndex)
	testIds.Close()
}
