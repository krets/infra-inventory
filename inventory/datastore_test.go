package inventory

import (
	"path/filepath"
	"testing"
)

var (
	testEssHost        = "localhost"
	testEssPort        = 9200
	testMappingFile, _ = filepath.Abs("../etc/mapping.json")
	testIndexM         = "test_index_with_mapping"
	testIndex          = "test_index"
)

func Test_NewElasticsearchDatastore_MappingFile(t *testing.T) {
	t.Logf("Mapping file: %s", testMappingFile)
	e, err := NewElasticsearchDatastore(testEssHost, testEssPort, testIndexM, testMappingFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	e.Conn.DeleteIndex(e.Index)
	e.Conn.DeleteIndex(e.VersionIndex)
	e.Conn.Close()
}

func Test_NewElasticsearchDatastore(t *testing.T) {
	e, err := NewElasticsearchDatastore(testEssHost, testEssPort, testIndex)
	if err != nil {
		t.Fatalf("%s", err)
	}
	e.Conn.DeleteIndex(e.Index)
	e.Conn.DeleteIndex(e.VersionIndex)
	e.Conn.Close()
}

func Test_ElasticsearchDatastore_Info(t *testing.T) {
	e, _ := NewElasticsearchDatastore(testEssHost, testEssPort, testIndex)

	info, err := e.Info()
	if err != nil {
		t.Fatalf("%s", err)
	} else {
		t.Logf("%#v\n", info)
	}
	e.Conn.DeleteIndex(e.Index)
	e.Conn.DeleteIndex(e.VersionIndex)
	e.Conn.Close()
}
