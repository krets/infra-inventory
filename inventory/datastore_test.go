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
	//testEssDs, _       = NewElasticsearchDatastore(testEssHost, testEssPort, testIndex)
	/*
	    testData = map[string]string{
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
	*/
)

func Test_NewElasticsearchDatastore_MappingFile(t *testing.T) {
	t.Logf("Mapping file: %s", testMappingFile)
	e, err := NewElasticsearchDatastore(testEssHost, testEssPort, testIndexM, testMappingFile)
	if err != nil {
		t.Fatalf("%s", err)
	}
	e.Conn.DeleteIndex(testIndexM)
	e.Conn.Close()
}

func Test_NewElasticsearchDatastore(t *testing.T) {
	e, err := NewElasticsearchDatastore(testEssHost, testEssPort, testIndex)
	if err != nil {
		t.Fatalf("%s", err)
	}
	e.Conn.DeleteIndex(testIndex)
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
	e.Conn.DeleteIndex(testIndex)
	e.Conn.Close()
}
