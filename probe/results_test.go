
package probe

import (
	"os"
	"reflect"
	"testing"
)

func TestAll(t *testing.T) {
	r := []Result{*CreateTestResult(), *CreateTestResult(), *CreateTestResult()}
	r[0].Name = "Test1 Name"
	r[1].Name = "Test2 Name"
	r[2].Name = "Test3 Name"

	SetResultsData(r)
	x := GetResultData("Test1 Name")
	if reflect.DeepEqual(x, r[0]) {
		t.Errorf("GetResult(\"Test1 Name\") = %v, expected %v", x, r[0])
	}

	// ensure we dont save or load from '-'
	if err := SaveDataToFile("-"); err != nil {
		t.Errorf("SaveToFile(-) error: %s", err)
	}

	if err := LoadDataFromFile("-"); err != nil {
		t.Errorf("LoadFromFile(-) error: %s", err)
	}

	filename := "/tmp/easeprobe/data.yaml"
	if err := os.MkdirAll("/tmp/easeprobe", 0755); err != nil {
		t.Errorf("Mkdirall(\"/tmp/easeprobe\") error: %v", err)
	}

	if err := SaveDataToFile(filename); err != nil {
		t.Errorf("SaveToFile(%s) error: %s", filename, err)
	}

	if err := LoadDataFromFile(filename); err != nil {
		t.Errorf("LoadFromFile(%s) error: %s", filename, err)
	}

	if reflect.DeepEqual(resultData["Test1 Name"], r[0]) {
		t.Errorf("LoadFromFile(\"%s\") = %v, expected %v", filename, resultData["Test1 Name"], r[0])
	}

	if err := os.RemoveAll("/tmp/easeprobe"); err != nil {
		t.Errorf("RemoveAll(\"/tmp/easeprobe\") = %v, expected nil", err)
	}
}
