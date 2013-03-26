package main

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"io"
	"math"
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
	Subject   string
	Reference string
	Details   string
	Email     string
	At        time.Time
}

//global variables that we will have to load once
var busStops []BusStop
var buses []Bus

func initData() error {
	//Load data in the busstops file.
	fmt.Println("Info: model.go: initData(): Starting initializing data.")

	var err error

	busStopsBuf := readFile("gocode/busstops.json")
	// fmt.Println("Info: model.go: initData(): Full busstops file as bytes: ", string(busStopsBuf))
	if err = json.Unmarshal(busStopsBuf, &busStops); err != nil {
		panic("Error: model.go: initData(): Error unmarshaling busstops data: " + err.Error())
	}
	// fmt.Println("Info: model.go: initData(): Full busstops file as objects: ", busStops)

	busesBuf := readFile("gocode/buses.json")
	// fmt.Println("Info: model.go: initData(): Full busstops file as bytes: ", string(busesBuf))
	if err = json.Unmarshal(busesBuf, &buses); err != nil {
		panic("Error: model.go: initData(): Error unmarshaling buses data: " + err.Error())
	}
	// fmt.Println("Info: model.go: initData(): Full busstops file as objects: ", buses)

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

func getDirectBuses(from, to string) []Bus {
	var selBuses []Bus

	for _, bus := range buses {
		fromExists, toExists, reverse := false, false, false

		for _, stopName := range bus.BusStopsA {
			if stopName == from {
				fromExists = true
			}
			if stopName == to {
				toExists = true
				if toExists && !fromExists { //we need to reverse the order of stops
					reverse = true
				}
			}

			if fromExists && toExists {
				if reverse {
					revBus := reverseBus(bus)
					selBuses = append(selBuses, revBus)
				} else {
					selBuses = append(selBuses, bus)
				}
				break
			}
		}

	}

	return selBuses
}

func reverseBus(bus Bus) Bus {
	var revStops []string
	for j := len(bus.BusStopsA) - 1; j >= 0; j = j - 1 {
		revStops = append(revStops, bus.BusStopsA[j])
	}
	revBus := Bus{
		bus.Number,
		revStops,
		bus.BusStopsB, //for now just include this
	}
	return revBus
}

// this function should be called only after direct buses have been checked and none found
func get1HopBuses(from, to string, c appengine.Context) []Bus {
	var fromBuses, toBuses, selBuses []Bus

	for _, bus := range buses {
		for _, stopName := range bus.BusStopsA {
			if stopName == from {
				fromBuses = append(fromBuses, bus)
			}
			if stopName == to {
				toBuses = append(toBuses, bus)
			}
		}
	}
	//at this point we have all the buses that have the stop in fromBuses and that have the destination stop in toBuses

	var firstBuses, secondBuses []Bus
	for _, firstBus := range fromBuses {
		for _, secondBus := range toBuses {
			found := false
			for _, firstStop := range firstBus.BusStopsA {
				if found {
					break
				}
				for _, secondStop := range secondBus.BusStopsA {
					if firstStop == secondStop {
						firstBuses = append(firstBuses, firstBus)
						secondBuses = append(secondBuses, secondBus)
						found = true
					}
				}
			}
		}
	}
	//at this point we have all the buses pairs that have the origin and destination stops

	//find nearest common stop
	for cnt := 0; cnt < len(firstBuses); cnt++ {
		firstBus := firstBuses[cnt]
		secondBus := secondBuses[cnt]

		fromStopPos := -1
		for i := 0; i < len(firstBus.BusStopsA); i++ {
			if from == firstBus.BusStopsA[i] {
				fromStopPos = i
				break
			}
		}

		toStopPos := -1
		for i := 0; i < len(secondBus.BusStopsA); i++ {
			if to == secondBus.BusStopsA[i] {
				toStopPos = i
			}
		}

		// nearestStop := ""
		nearestDist := math.MaxFloat64
		revFrom, revTo := false, false
		for i := 0; i < len(firstBus.BusStopsA); i++ {
			for j := 0; j < len(secondBus.BusStopsA); j++ {
				if firstBus.BusStopsA[i] == secondBus.BusStopsA[j] {
					newDist := math.Abs(float64(i-fromStopPos)) + math.Abs(float64(j-toStopPos))

					if newDist < nearestDist {
						nearestDist = newDist
						// nearestStop = firstBus.BusStopsA[i]
						if i < fromStopPos {
							revFrom = true
						} else {
							revFrom = false
						}
						if j > toStopPos {
							revTo = true
						} else {
							revTo = false
						}

					}
				}
			}
		}

		// order the items correctly
		if revFrom {
			selBuses = append(selBuses, reverseBus(firstBus))
		} else {
			selBuses = append(selBuses, firstBus)
		}
		if revTo {
			selBuses = append(selBuses, reverseBus(secondBus))
		} else {
			selBuses = append(selBuses, secondBus)
		}
	}

	c.Infof("model.go: get1HopBuses(): Selected Buses: ", selBuses, "\nTotal size: ", len(selBuses))

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

func addFeedback(c appengine.Context, subject, reference, details, email string) error {
	feedback := Feedback{
		subject,
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
