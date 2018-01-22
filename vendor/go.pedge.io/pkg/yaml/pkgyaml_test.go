package pkgyaml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

const (
	numFiles = 4
)

func TestFilesEqual(t *testing.T) {
	for i := 1; i <= numFiles; i++ {
		testFileEqual(t, i)
	}
}

func TestParse(t *testing.T) {
	for i := 1; i <= numFiles; i++ {
		testParse(t, i)
	}
}

func testFileEqual(t *testing.T, i int) {
	yamlData := testReadData(t, fmt.Sprintf("testdata/%d.yml", i))
	expectedJSONData := testRemarshalJSON(t, testReadData(t, fmt.Sprintf("testdata/%d.json", i)), i)
	jsonData, err := ToJSON(yamlData, ToJSONOptions{})
	if err != nil {
		t.Fatalf("transform error for file %d: %v", i, err)
	}
	if !bytes.Equal(jsonData, expectedJSONData) {
		t.Errorf("transform mismatch for file %d: expected %s, got %s", i, string(expectedJSONData), string(jsonData))
	}
}

func testParse(t *testing.T, i int) {
	var yamlMap interface{}
	if err := ParseYAMLOrJSON(fmt.Sprintf("testdata/%d.yml", i), &yamlMap); err != nil {
		t.Fatal(err)
	}
	var jsonMap interface{}
	if err := ParseYAMLOrJSON(fmt.Sprintf("testdata/%d.json", i), &jsonMap); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(yamlMap, jsonMap) {
		t.Fatalf("expected %v and %v to be equal", yamlMap, jsonMap)
	}
}

func testReadData(t *testing.T, path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open error for file %s: %v", path, err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("readall error for file %s: %v", path, err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close error for file %s: %v", path, err)
	}
	return data
}

func testRemarshalJSON(t *testing.T, jsonData []byte, i int) []byte {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		t.Fatalf("unmarshal error for file %d: %v", i, err)
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal error for file %d: %v", i, err)
	}
	return jsonData
}
