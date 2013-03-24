package main

import (
	"appengine"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var templates *template.Template

func init() {

	var err error
	templates, err = template.ParseGlob("templates/*.tmpl")
	if err != nil {
		panic("main.go: init(): error parsing templates: " + err.Error())
	}

	err = initData() //read in all the data from the json files
	if err != nil {
		panic("main.go: init(): error initializing data: " + err.Error())
	}

	m := mux.NewRouter()

	m.HandleFunc("/", rootHandler)
	m.HandleFunc("/r", routeAtoBAjaxHandler).Methods("GET")
	m.HandleFunc("/b", busNumberAjaxHandler).Methods("GET")
	m.HandleFunc("/f", feedbackAjaxHandler).Methods("POST")
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

	busNumbers := getBusNumbers()
	if len(busNumbers) == 0 {
		c.Errorf("main.go: rootHandler(): There were no bus numbers.")
	}

	if err := templates.ExecuteTemplate(w, "index", struct {
		BusStops   []string
		BusNumbers []string
	}{
		busStopNames,
		busNumbers,
	}); err != nil {
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

func busNumberAjaxHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	c.Infof("main.go: busNumberAjaxHandler(): Request received to url: %s", r.URL.RequestURI())

	number := r.FormValue("number")
	if len(number) == 0 {
		c.Errorf("main.go: busNumberAjaxHandler(): Invalid input.  No bus number specified.")
		http.Error(w, "Invalid input.  No bus number specified.", http.StatusBadRequest)
		return
	}

	bus := getBus(number)
	writeSuccessJSONResponse(c, w, bus)

}

func feedbackAjaxHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	c.Infof("main.go: feedbackAjaxHandler(): Request received to url: %s", r.URL.RequestURI())

	subject := r.FormValue("feedbackSubject")
	reference := r.FormValue("feedbackReference")
	details := r.FormValue("feedbackDetails")
	email := r.FormValue("feedbackEmail")

	err := addFeedback(c, subject, reference, details, email)
	if err != nil {
		http.Error(w, "Error saving feedback: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeSuccessJSONResponse(c, w, "Successfully added issue/feedback.")
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
