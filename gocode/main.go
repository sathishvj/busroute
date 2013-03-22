package main

import (
	"appengine"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	// "appengine/datastore"
	// "time"
)

var templates *template.Template

func init() {

	var err error
	templates, err = template.ParseGlob("templates/*.tmpl")
	if err != nil {
		panic("main.go: init(): error parsing templates: " + err.Error())
	}

	m := mux.NewRouter()

	m.HandleFunc("/", rootHandler)
	m.HandleFunc("/r", routeAtoBAjaxHandler).Methods("GET")
	// m.HandleFunc("/{path:.*}", pageNotFoundHandler).Methods("GET")

	http.Handle("/", m)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	c.Infof("main.go: rootHandler(): Request received to root: %s", r.URL.RequestURI())

	busStopNames := getBusStopNames()
	if len(busStopNames) == 0 {
		c.Errorf("main.go: rootHandler(): There were no bus stops.")
	}

	if err := templates.ExecuteTemplate(w, "index", struct{ BusStops []string }{busStopNames}); err != nil {
		c.Errorf("main.go: rootHandler(): Error executing template: ", err)
	}
}

func routeAtoBAjaxHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	c.Infof("main.go: fromAtoB Handler(): Request received to url: %s", r.URL.RequestURI())

	from := r.FormValue("from")
	to := r.FormValue("to")
	if len(from) == 0 || len(to) == 0 {
		http.Error(w, "Invalid input.  Both from and to values have to be present.", http.StatusBadRequest)
		return
	}

	buses := getBuses(from, to)
	writeSuccessJSONResponse(c, w, buses)

}

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	c.Infof("main.go: pageNotFoundHandler(): Request received to url: %s", r.URL.RequestURI())
	fmt.Fprintf(w, "Showing pageNotFound page!")

}

func writeSuccessJSONResponse(c appengine.Context, w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		c.Errorf("main.go: writeSuccessJSONResponse: Error marshalling data: ", err)
		return
	}
	fmt.Fprint(w, string(b))
	c.Infof("main.go: writeSuccessJSONResponse: Successfully wrote data to ResponseWriter.")
}
