/* EHRDE's widget functions
v1: supports functions to display data in jxqTree and vis.js network widgets.
*/

package data

import (
	"strings"
)

// general function to deal with jqxTree widget's data generation
func GenerateTreeData(dataMap map[string]interface{}, numNodes int) ([]interface{}, error) {
	//take dataMap and built needed treeWidgetSet

	//1st initialize with total num of nodes
	treeWidgetSet := make([]interface{}, numNodes)

	//2nd fill treeWidgetSet with data
	// [{id: ..., parentId: ..., label: ...}, {...}, ... ]
	// treeWidgetSetSlice = make(map[string]interface{}); treeWidgetSetSlice["id"] = ..id.. etc.
	//NOTE: import to use unique ids. computable and understandable concatenation like 'kk+k'
	//      seems reasonable right now
	i := 0
	for k, v := range dataMap {
		//add template
		treeWidgetSet = addToTreeWidgetSet(k, "-1", k, i, treeWidgetSet)
		i++
		for kk, vv := range v.(map[string]interface{}) {
			//add archetype
			//fmt.Println("----- DATAMAP trying to add: ", kk+k) //NOTE debug
			treeWidgetSet = addToTreeWidgetSet(kk+k, k, vv.(map[string]interface{})["name"].(string), i, treeWidgetSet)
			i++
			// TEST ------------ TODO
			//fmt.Println("----- DATAMAP - archetype added to treemap: ", vv) //NOTE debug
			for _, vvv := range vv.(map[string]interface{}) {
				switch vvvType := vvv.(type) {
				case map[string]interface{}:
					for kkkk, vvvv := range vvv.(map[string]interface{}) {
						//add entry - 1st get attribute 'name'
						name := vvvv.(map[string]string)["name"]
						treeWidgetSet = addToTreeWidgetSet(kkkk+kk+k, kk+k, name, i, treeWidgetSet)
						i++
					}
				default:
					//this ignores other attributes of archetype
					//why do I need a stupid usage of vvvType?!
					vvvType = vvvType
				}
			}
		}
	}
	return treeWidgetSet, nil
}

// general function to deal with visjs widget's data generation
//redundancy is on purpose to keep it more modular
func GenerateVisData(dataMap map[string]interface{}, numNodes int) ([]interface{}, []interface{}, error) {
	// number of nodes = number of nodes (templates, archetypes, entries)
	visNodes := make([]interface{}, numNodes)
	// number of edges = number of archetypes + number of entries
	// --> reassembled with: number of nodes - number of templates
	visEdges := make([]interface{}, numNodes-len(dataMap))

	//nodes: [{group: , id: , label: , level: }, {..}, ..]
	//       with group 0 template, 1 archetypes, 2 entries - same with level
	//edges: [{from: #id, to: #id}, {..}, ..]

	n := 0 //nodes
	e := 0 //edges
	for k, v := range dataMap {
		//add template
		visNodes = addToVisNodes(0, k, k, "", 0, n, visNodes)
		n++
		for kk, vv := range v.(map[string]interface{}) {
			//add archetype
			visNodes = addToVisNodes(1, kk+k, vv.(map[string]interface{})["name"].(string), "", 1, n, visNodes)
			n++
			visEdges = addToVisEdges(k, kk+k, e, visEdges)
			e++
			for _, vvv := range vv.(map[string]interface{}) {
				switch vvvType := vvv.(type) {
				case map[string]interface{}:
					for kkkk, vvvv := range vvv.(map[string]interface{}) {
						//add datavalue
						name := vvvv.(map[string]string)["name"]
						count := vvvv.(map[string]string)["count"]
						visNodes = addToVisNodes(2, kkkk+kk+k, name, "Count: "+count, 2, n, visNodes)
						n++
						visEdges = addToVisEdges(kk+k, kkkk+kk+k, e, visEdges)
						e++
					}
				default:
					//this ignores other attributes of archetype
					vvvType = vvvType
				}
			}
		}
	}
	return visNodes, visEdges, nil
}

//modular function to add node to visNodes array
func addToVisNodes(group int, id string, label string, title string, level int, pos int, visNodes []interface{}) []interface{} {
	slice := make(map[string]interface{})
	slice["group"] = group
	slice["id"] = id
	//check if label is too long and should have line-breaks when node is DATA_VALUE
	if group == 2 {
		label = addLinebreaks(label, 8)
	}
	slice["label"] = label
	//NOTE tooltip to display count as proof of concept
	if title != "" {
		slice["title"] = title
	}
	slice["level"] = level
	visNodes[pos] = slice

	return visNodes
}

//recursively add linebreak chars until the whole is max wide as maxLength
//TODO: add support to detect spaces (at least if next to break) and break at space
func addLinebreaks(s string, maxLength int) string {
	if len(s) > maxLength {
		temp := []string{s[0:8], s[8:len(s)]}
		temp[1] = addLinebreaks(temp[1], maxLength)
		s = strings.Join(temp, string('\n'))
	}
	return s
}

//modular function to add edge to visEdges array - 'from' and 'to' need ids of vis nodes
func addToVisEdges(from, to string, pos int, visEdges []interface{}) []interface{} {
	slice := make(map[string]interface{})
	slice["from"] = from
	slice["to"] = to
	visEdges[pos] = slice

	return visEdges
}

//modular function to add a slice to treeWidgetSet array
func addToTreeWidgetSet(id, parentId, label string, pos int, treeWidgetSet []interface{}) []interface{} {
	slice := make(map[string]interface{})
	slice["id"] = id
	slice["parentId"] = parentId
	slice["label"] = label
	treeWidgetSet[pos] = slice

	return treeWidgetSet
}

//outsourced function to count all nodes in dataMap
func CountNodes(dataMap map[string]interface{}) int {
	count := 0
	for _, v := range dataMap {
		count++ //adds each template

		//no need for type switch, dataMap is generated with map values on 1st level
		for _, vv := range v.(map[string]interface{}) {
			count++ //adds each archetype of that template
			for _, vvv := range vv.(map[string]interface{}) {
				//inside the archetypes data map
				switch vvvType := vvv.(type) {
				case map[string]interface{}:
					//found "entries"
					for _, _ = range vvv.(map[string]interface{}) {
						//at 4th level - entry's attributes at 5th
						count++ //adds each entry
					}
				default:
					//ignore other archetype attributes for now
					vvvType = vvvType
				}
			}
		}
	}
	return count
}
