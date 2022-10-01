package discovery

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ReturnRestResponse(w http.ResponseWriter, v interface{}, successStatus int) {

	jsonBytes, err := json.Marshal(v)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Error: %s", err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(successStatus)
		w.Write(jsonBytes)
	}
}
