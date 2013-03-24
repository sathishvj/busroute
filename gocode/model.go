package main

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

type Location struct {
	Name                string
	Latitude, Longitude float32 `json:",string"`
}

type BusStop struct {
	Name      string
	Details   string     //to help identify the place.  Eg. on same side as ABC hospital
	Locations []Location //because there could be multiple bus stops with the same name

}

type Bus struct {
	Number string

	//sometimes the stops on the up direction is different from the down direction
	BusStopsA []string //an ordered collection of stops. Each string corresponds to BusStop.Name
	BusStopsB []string //an ordered collection of stops in reverse direction. Each string corresponds to BusStop.Name
}

type Feedback struct {
	Category    string
	SubCategory string
	Reference   string
	Details     string
	Email       string
	At          time.Time
}

//global variables that we will have to load once
var busStops []BusStop
var buses []Bus

func initData() error {
	//Load data in the busstops file.

	var err error

	busStopsBuf := readFile("gocode/busstops.json")
	fmt.Println("Info: model.go: initData(): Full busstops file as bytes: ", string(busStopsBuf))
	if err = json.Unmarshal(busStopsBuf, &busStops); err != nil {
		panic("Error: model.go: initData(): Error unmarshaling busstops data: " + err.Error())
	}
	fmt.Println("Info: model.go: initData(): Full busstops file as objects: ", busStops)

	busesBuf := readFile("gocode/buses.json")
	fmt.Println("Info: model.go: initData(): Full busstops file as bytes: ", string(busesBuf))
	if err = json.Unmarshal(busesBuf, &buses); err != nil {
		panic("Error: model.go: initData(): Error unmarshaling buses data: " + err.Error())
	}
	fmt.Println("Info: model.go: initData(): Full busstops file as objects: ", buses)

	return nil
}

func readFile(path string) []byte {
	// file, err := os.OpenFile("gocode/busstops.json", os.O_RDONLY, 0666)
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		// fmt.Printf("model.go: readFile(): Error opening %s: %s\n", path, err.Error())
		// return err
		panic("model.go: readFile(): Error opening " + path + ":" + err.Error())
	}

	defer func() {
		if file.Close() != nil {
			panic(err)
		}
	}()

	var fullBuf []byte
	tempBuf := make([]byte, 1024)
	n := 0
	for {
		n, err = file.Read(tempBuf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		//need to reslice the buffer using the number of bytes actually read.
		tempBuf = tempBuf[0:n]

		fullBuf = append(fullBuf, tempBuf...)
	}

	return fullBuf
}

func getBusStopNames() []string {
	var names []string
	for _, stop := range busStops {
		names = append(names, stop.Name)
	}
	sort.Strings(names)
	return names
}

func getBuses(from, to string) []Bus {
	var selBuses []Bus
	for _, bus := range buses {
		fromExists, toExists := false, false
		for _, stopName := range bus.BusStopsA {
			if stopName == from {
				fromExists = true
			}
			if stopName == to {
				toExists = true
			}

			if fromExists && toExists {
				selBuses = append(selBuses, bus)
				break
			}
		}
	}

	return selBuses
}

func getBusNumbers() []string {
	var busNumbers []string
	for _, bus := range buses {
		busNumbers = append(busNumbers, bus.Number)
	}

	sort.Strings(busNumbers)

	return busNumbers
}

func getBus(number string) *Bus {
	for _, bus := range buses {
		if bus.Number == number {
			return &bus
		}
	}

	return nil
}

func addFeedback(c appengine.Context, category, subcategory, reference, details, email string) error {
	feedback := Feedback{
		category,
		subcategory,
		reference,
		details,
		email,
		time.Now(),
	}

	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Feedback", nil), &feedback)
	if err != nil {
		c.Errorf("model.go: addFeedback: Error putting feedback: ", err.Error())
	} else {
		c.Infof("model.go: addFeedback: Successfully put feedback: ", feedback)
	}
	return err
}
