package main

import (
	"appengine"
	"appengine/datastore"
	"sort"
	"time"
)

type BusStop struct {
	Name                string
	Detail              string //to help identify the place.  Eg. on same side as ABC hospital
	Latitude, Longitude float32
}

type Bus struct {
	Number string

	//sometimes the stops on the up direction is different from the down direction
	BusStopsA []string //an ordered collection of stops
	BusStopsB []string //an ordered collection of stops in reverse direction
}

type Feedback struct {
	Category    string
	SubCategory string
	Reference   string
	Details     string
	Email       string
	At          time.Time
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
