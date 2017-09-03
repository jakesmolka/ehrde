/* EHRDE's detail functions
v1: details are just about 'count' values right now. displaying is tailored for bootstrap modals and index.html's JS
*/

package data

import (
	"errors"
	"fmt"
	"strings"
)

//gets all *count* values of each DATA_VALUE using AQL and adding them to the D_V's json node with ["count"] key
func AddCountValues(dataMap map[string]interface{}) (map[string]interface{}, error) {
	for k, v := range dataMap {
		//no need for type switch, dataMap is generated with map values on 1st level
		for _, vv := range v.(map[string]interface{}) {
			for _, vvv := range vv.(map[string]interface{}) {
				//inside the archetypes data map
				switch vvvType := vvv.(type) {
				case map[string]interface{}:
					//found "datavalues"
					for _, vvvv := range vvv.(map[string]interface{}) {
						//at 4th level - datavalue's attributes at 5th
						//getting count value through AQL for each datavalue
						count, err := countDatavalue(vvvv.(map[string]string)["aqlPath"], k)
						if err != nil {
							return dataMap, errors.New("getCountValues() -> countDatavalue()")
						}
						//NOTE: json numbers are always parsed as float64. count is always int though
						countString := fmt.Sprintf("%8.0f", count)
						vvvv.(map[string]string)["count"] = strings.TrimSpace(countString)
					}
				default:
					//ignore other archetype attributes for now
					vvvType = vvvType
				}
			}
		}
	}

	return dataMap, nil
}

// counts number of datavalue elements using AQL and given aqlPath
// returns number as float64
func countDatavalue(aqlPath, templateId string) (float64, error) {
	//1. build query 2. exec query 3. extract count

	//1
	//NOTE check if this is safe. bc the only problem could be the arbitrary length of datavalue's aqlpath
	//	and it's just split from the front, this should be fine...
	//prepare query variable
	query := "select count(a_a/"

	//1a: aql path to datavalue
	aqlDatavalue := strings.SplitN(aqlPath, "/", 3)[2]
	query += aqlDatavalue + ") from EHR e contains COMPOSITION a contains "

	//1b: rm type of archetype
	aqlArchetypeId := strings.SplitN(strings.SplitN(aqlPath, "[", 2)[1], "]", 2)[0] //needed for the next step too
	aqlArchetypeRmType := strings.SplitN(strings.SplitN(aqlArchetypeId, "-", 3)[2], ".", 2)[0]
	query += aqlArchetypeRmType + " a_a["

	//1c: whole identification of archetype
	//already calculated in last step
	query += aqlArchetypeId + "]"

	//1d: if templateId is parameter, add WHERE statement
	if templateId != "" {
		query += " where a/archetype_details/template_id/value='" + templateId + "'"
	}

	//2: exec query
	response, err := execAqlQuery(query)
	errorHandler(err, "countDatavalue()")

	//3:

	respMap := response.(map[string]interface{})
	//check for valid response
	if _, exists := respMap["resultSet"]; exists {
		if len(respMap["resultSet"].([]interface{})) != 1 {
			return 0, errors.New("countDatavalue() -> resultSet got more than 1 elements")
		} else {
			resultMap := respMap["resultSet"].([]interface{})[0].(map[string]interface{})
			//avoiding static naming of key
			var result float64
			for _, value := range resultMap {
				result = value.(float64)
			}
			return result, nil
		}
	} else {
		return 0, errors.New("countDatavalue() -> response is not valid")
	}
}

//function to get a distinct copy of dataMap with *count* values without template context.
func GetDatavaluesMap(dataMap map[string]interface{}) map[string]map[string]string {
	datavaluesMap := make(map[string]map[string]string)
	for k, v := range dataMap {
		//no need for type switch, dataMap is generated with map values on 1st level
		for kk, vv := range v.(map[string]interface{}) {
			for _, vvv := range vv.(map[string]interface{}) {
				//inside the archetypes data map
				switch vvvType := vvv.(type) {
				case map[string]interface{}:
					//found "datavalues"
					for kkkk, vvvv := range vvv.(map[string]interface{}) {
						//at 4th level - datavalue's attributes at 5th
						//fetching datavalue and saving it under the widgets id on plane json
						dv := vvvv.(map[string]string)

						//datavaluesMap[kkkk+kk+k] = dv
						//because this line is just referencing and I dont want to modify original too,
						//	element by element copy:
						datavaluesMap[kkkk+kk+k] = make(map[string]string)
						for key, value := range dv {
							datavaluesMap[kkkk+kk+k][key] = value
						}

						//---PROTOTYPE ADDITION---
						//re-query count values DWH wide (WITHOUT template context)
						//so detail modals will show count values from 'global' perspective
						//TODO should be refactored to do as less AQLs as possible
						count, err := countDatavalue(dv["aqlPath"], "")
						if err != nil {
							errors.New("getDatavaluesMap() -> countDatavalue()")
						}
						//NOTE: json numbers are always parsed as float64. count is always int though
						countString := fmt.Sprintf("%8.0f", count)
						datavaluesMap[kkkk+kk+k]["count"] = strings.TrimSpace(countString)
						//------------------------

					}
				default:
					//ignore other archetype attributes for now
					vvvType = vvvType
				}
			}
		}
	}

	return datavaluesMap
}
