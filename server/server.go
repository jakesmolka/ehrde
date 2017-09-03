/* EHRDE's server package - server.go
v1: basic server funtion to handle calls to the webapp.
*/

//The server package implements the EHRDE's dashboard server.
package server

import (
	"html/template"
	"log"
	"net/http"
)

//page is used as data structure when rendering the index.html template
type page struct {
	TreeJSON  template.JS
	EdgesJSON template.JS
	NodesJSON template.JS
}

// pre-parsed template files
//var templates = template.Must(template.ParseFiles("assets/templates/test.html", "assets/templates/index.html"))
var templates = template.Must(template.ParseFiles("assets/templates/index.html"))

// redirector for root URL
func RootHandler(w http.ResponseWriter, r *http.Request, widgetDataJson, visNodesJson, visEdgesJson []byte) {

	log.Print("---------- Dashboard opened by IP: ", r.RemoteAddr)

	templates.ExecuteTemplate(w, "index.html", &page{TreeJSON: template.JS(widgetDataJson),
		EdgesJSON: template.JS(visEdgesJson),
		NodesJSON: template.JS(visNodesJson)})

}
