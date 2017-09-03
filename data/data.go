/* EHRDE's data package - data.go:
v1: This file contains the basic functions of the data package.
*/

//The data package is EHRDE's data part. Getting, handling, providing, etc. all kind of internal data.
package data

import (
	"errors"
	"fmt"
	"strings"
)

// gets a map of all archtypes by going through an array of all templates and recursively searching for archetype's nodes
//TODO change name to ~ getDataByTemplates
func GetArchetypesByTemplates(templates []interface{}) (map[string]interface{}, error) {
	dataMap := make(map[string]interface{}) //should be done already
	var tempDataMap map[string]interface{}

	for i := range templates {
		template := templates[i].(map[string]interface{})
		//fmt.Println("getArchetypesByTemplates: for: ", template)

		templateId := template["templateId"]
		//GET http:.../template/id as map
		//fmt.Println("TEMPLATE: ", templateId) //NOTE debug
		templateDetails, err := HttpGetJsonThink("rest/v1/template/" + templateId.(string))
		//call searchInMap for 'nodeId' with values for openEHR archetypes

		lastArchetypeId := ""                      //reset for new template
		tempDataMap = make(map[string]interface{}) //reset temp map per template

		tempDataMap, err = searchInMap(templateDetails.(map[string]interface{}), tempDataMap, lastArchetypeId)
		if err != nil {
			return dataMap, errors.New("getArchetypesByTemplates: Error with searchInMap call")
		}

		//TODO merge tempDataMap into dataMap
		// 1st add new template node to dataMap
		// 2nd add data from tempDataMap under new template node in dataMap
		dataMap[templateId.(string)] = make(map[string]interface{})
		//dataMap[templateId].(map[string]interface{})["id"] = templateId //redundant?! is in key already
		//for every archetype found for this template
		for tempKey, tempValue := range tempDataMap {
			switch tempValueType := tempValue.(type) {
			case map[string]interface{}:
				//new archetype. value is it's definition
				//add whole archetype as node to template
				dataMap[templateId.(string)].(map[string]interface{})[tempKey] = tempValue
			default:
				fmt.Println("getArchetypesByTemplates: this should not happen. type: ", tempValueType)
			}
		}

	}
	return dataMap, nil
}

// searches through m for occurrences of archetype or datavalue nodes, extracts info and puts it into dataMap
func searchInMap(m, tempDataMap map[string]interface{}, lastArchetypeId string) (map[string]interface{}, error) {
	// 1st go through m's keys and values, check for *nodeName*
	//   and recursively call searchInMap when finding another map as value

	//TODO do on higher scope! unefficient right now!
	archetypeClasses := []string{"OBSERVATION", "EVALUATION", "INSTRUCTION", "ACTION", "SECTION"}

	datavalueClasses := []string{"DV_BOOLEAN", "DV_TEXT", "DV_CODED_TEXT", "DV_ORDINAL", "DV_QUANTITY", "DV_COUNT", "DV_PROPORTION", "DV_DATE", "DV_TIME", "DV_DATE_TIME", "DV_DURATION"}

	//1st check if current json nesting level is of interest (ie. contains archetype or datavalue)
	if _, exists := m["rmType"]; exists {
		rmType := m["rmType"].(string)
		//1.a archetype?
		for i := range archetypeClasses {
			if strings.EqualFold(rmType, archetypeClasses[i]) {
				//rmType reveals archetype

				//CONTINUE HERE - saving the actual data
				//gather relevant data. in case of archetype, nodeId is a good id
				name := m["name"].(string)
				//rmType already set
				nodeId := m["nodeId"].(string)
				aqlPath := m["aqlPath"].(string)

				//save data - where nodeId is the *key*
				tempDataMap[nodeId] = make(map[string]interface{})
				tempDataMap[nodeId].(map[string]interface{})["name"] = name
				tempDataMap[nodeId].(map[string]interface{})["rmType"] = rmType
				tempDataMap[nodeId].(map[string]interface{})["aqlPath"] = aqlPath
				//init map for entries
				if _, ok := tempDataMap[nodeId].(map[string]interface{})["datavalues"]; !ok {
					//if not already existing
					//TODO check if this shouldn't be necessary here, seems faulty
					tempDataMap[nodeId].(map[string]interface{})["datavalues"] = make(map[string]interface{})
				}

				//set reminder of id, so entries can access the right map
				lastArchetypeId = nodeId

				//fmt.Println("ARCHETYPE: ", m["name"].(string), id, rmType) //NOTE debug

				break //quit range
			}
		}

		//1.b datavalue? only possible if there was an archetype before (faulty otherwise)
		if lastArchetypeId != "" {
			for i := range datavalueClasses {
				if strings.EqualFold(rmType, datavalueClasses[i]) {
					//rmType reveals datavalue

					if _, ok := m["name"]; ok {
						//all elses will break the for and switch and continues with next node/element
						if _, ok := m["nodeId"]; ok {
							if m["nodeId"] == "" {
								//in case there is no id; special cases
								break
							}
							if _, ok := m["aqlPath"]; ok {
								//found all attributes -> just fine!
							} else {
								break
							}
						} else {
							break
						}
					} else {
						break
					}

					//CONTINUE HERE - saving the actual data
					//gather relevant data - in case of datavalue nodeId is not a good id. using 'id' instead
					id := m["id"].(string)
					name := m["name"].(string)
					//rmType already set
					//nodeId := m["nodeId"].(string) //not important right now
					aqlPath := m["aqlPath"].(string)

					//save data - first build slice
					dataSlice := make(map[string]string)
					dataSlice["name"] = name
					dataSlice["rmType"] = rmType
					dataSlice["aqlPath"] = aqlPath

					//fmt.Println("DATA_VALUE: ", nodeId, name, rmType, lastArchetypeId) //NOTE debug

					//check if "entries" map exists
					//TODO: why? should exist already!
					if _, ok := tempDataMap[lastArchetypeId]; !ok {
						//if map doesn't exist
						tempDataMap[lastArchetypeId] = make(map[string]interface{})
						tempDataMap[lastArchetypeId].(map[string]interface{})["datavalues"] = make(map[string]interface{})
					}

					tempDataMap[lastArchetypeId].(map[string]interface{})["datavalues"].(map[string]interface{})[id] = dataSlice

					//fmt.Println("DATAVALUE: ", rmType) //NOTE debug
					break //quit range
				}
			}
		}
	}

	//2nd find the next json nesting level and proceed with it
	for _, value := range m {
		switch valueType := value.(type) {
		//2.a next nesting level?
		case map[string]interface{}:
			//call searchInMap with next nesting level
			tempDataMap, err := searchInMap(value.(map[string]interface{}), tempDataMap, lastArchetypeId)
			tempDataMap = tempDataMap //TODO remove
			errorHandler(err, "searchInMap: error in recursive step")

		//2.b array? can contain a map
		case []interface{}:
			//next json node is array
			for i := range value.([]interface{}) {
				//take a look at each element of that array
				slice := value.([]interface{})[i]
				switch vType := slice.(type) {
				case map[string]interface{}:
					tempDataMap, err := searchInMap(slice.(map[string]interface{}), tempDataMap, lastArchetypeId)
					tempDataMap = tempDataMap //TODO remove
					errorHandler(err, "searchInMap: error in previous recursive step")
				default:
					vType = vType
				}
			}
		//2.c ignore everything else - most likely dead end, no further nesting
		default:
			valueType = valueType
		}
	}

	return tempDataMap, nil
}
