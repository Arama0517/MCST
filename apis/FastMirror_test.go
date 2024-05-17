package apis_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/apis"
)

func TestGetFastMirrorDatas(t *testing.T) {
	data, err := apis.GetFastMirrorDatas()
	if err != nil {
		t.Error(err)
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	fmt.Print(string(jsonData))
}
