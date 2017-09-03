package data

import (
	"fmt"
	"testing"
)

func TestCountDatavalue(t *testing.T) {
	result, err := countDatavalue("/content[openEHR-EHR-OBSERVATION.body_weight.v1]/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value")

	if err != nil {
		fmt.Println(err)
		t.Error("error")
	} else {
		fmt.Println(result)
	}
}

func TestGetCountValues(t *testing.T) {
	dataMap := make(map[string]interface{})
	dataMap["template1"] = make(map[string]interface{})
	t1Map := dataMap["template1"].(map[string]interface{})
	t1Map["archetype1"] = make(map[string]interface{})
	a1t1Map := t1Map["archetype1"].(map[string]interface{})
	a1t1Map["datavalues"] = make(map[string]interface{})
	da1t1Map := a1t1Map["datavalues"].(map[string]interface{})
	da1t1Map["dv1"] = make(map[string]string)
	dvda1t1Map := da1t1Map["dv1"].(map[string]string)

	dvda1t1Map["aqlPath"] = "/content[openEHR-EHR-OBSERVATION.body_weight.v1]/data[at0002]/events[at0003]/data[at0001]/items[at0004]/value"

	dataMap, err := getCountValues(dataMap)
	if err != nil {
		t.Error("err getCountValues")
	} else {
		fmt.Println(dataMap)
	}
}
