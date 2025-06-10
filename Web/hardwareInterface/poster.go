package hardwareInterface

import (
	"bytes"
	"com.quintindev/WebShed/config"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var CFG = config.Load()
var baseBackendAddr = "http://localhost:" + CFG.BackendPort

func PostJSON(url string, data any) string {
	url = baseBackendAddr + url
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func GetJSON(url string) any {
	url = baseBackendAddr + url
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println("error closing response body:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result any
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}

	return result
}

func GetJSONError(url string) (any, error) {
	url = baseBackendAddr + url
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if cerr := Body.Close(); cerr != nil {
			log.Println("error closing response body:", cerr)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func PostSetLock(state bool) {
	backendPayload := map[string]bool{
		"state": state,
	}
	PostJSON("/api/setlock", backendPayload)
}

func GetGetArmed() bool {
	data := GetJSON("/api/getarmed")
	obj, ok := data.(map[string]interface{})
	if !ok {
		log.Fatal("unexpected JSON structure")
	}

	newStateVal, ok := obj["state"]
	if !ok {
		log.Fatal("key 'state' not found")
	}

	newState, ok := newStateVal.(bool)
	if !ok {
		log.Fatal("state is not a boolean")
	}

	return newState
}

func GetExpiredCodes() {
	GetJSON("/api/expireoldcodes")
}
