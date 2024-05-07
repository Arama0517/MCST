package lib

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDetectJava(t *testing.T) {
	java := DetectJava()
	jsonJava, err := json.Marshal(java)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonJava))
}
