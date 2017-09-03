/* EHRDE's rest functions
v1: rest functions are basic http GET, and wrapper with json and AQL support
*/

package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//struct to read config.json
type Configuration struct {
	BaseUrl  string
	User     string
	Password string
}

//config var
var Config Configuration

//init like func to load the config file before exec REST calls
func LoadConfig() error {
	file, err := ioutil.ReadFile("config.json")
	errorHandler(err, "")
	Config = Configuration{}
	err = json.Unmarshal(file, &Config)
	errorHandler(err, "")

	return nil
}

//wrapper for httpGetJson to execute aql queries
func execAqlQuery(query string) (interface{}, error) {
	url := "rest/v1/query/?aql=" + url.QueryEscape(query)
	data, err := HttpGetJsonThink(url)
	errorHandler(err, "execAqlQuery()")

	return data, nil
}

//wrapper for HttpGetJson for Think!EHR Platform and with url, user, pass read from config.json
func HttpGetJsonThink(requestUrl string) (interface{}, error) {
	return HttpGetJson(Config.BaseUrl+requestUrl, Config.User, Config.Password)
}

//wrapper for httpGet to get already json.marshaled result
func HttpGetJson(requestUrl, user, pass string) (interface{}, error) {
	body, err := httpGet(requestUrl, user, pass)
	if err != nil {
		fmt.Println("httpGetJson: fetching data failed!")
		//TODO: PANIC
	}

	//NOTE debug
	//err = ioutil.WriteFile("response.json", body, 0644)
	//errorHandler(err, "")

	data, err := jsonToMap(body)
	if err != nil {
		//
	}

	return data, nil
}

//basic httpGet function
func httpGet(requestUrl, user, pass string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", requestUrl, nil)
	request.SetBasicAuth(user, pass)
	resp, err := client.Do(request)
	errorHandler(err, "httpGet() - Do()")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	errorHandler(err, "httpGet() - ReadAll()")

	return body, nil
}
