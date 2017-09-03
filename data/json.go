/* EHRDE's json functions
v1: just a helper function to get a Golang map from json
*/

package data

import (
	"encoding/json"
)

// get GO-style representation of JSON
func jsonToMap(b []byte) (interface{}, error) {
	var f interface{}
	err := json.Unmarshal(b, &f)
	return f, err
}
