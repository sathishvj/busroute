package main

import (
	"sort"
)

type BusStop struct {
	Name                string
	Detail              string //to help identify the place.  Eg. on same side as ABC hospital
	Latitude, Longitude float32
}

type Bus struct {
	Number string

	//sometimes the stops on the up direction is different from the down direction
	BusStopsA map[string]string //an ordered collection of stops
	BusStopsB map[string]string //an ordered collection of stops in reverse direction
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
