package main

var busStops = []BusStop{
	BusStop{Name: "MG Road"},
	BusStop{Name: "Brigade Road"},
	BusStop{Name: "JP Nagar"},
	BusStop{Name: "Jayanagar 5th Block"},
}

var buses = []Bus{
	Bus{
		"292A",
		[]string{
			"MG Road", "Trinity Circle", "Brigade Road",
		},
		[]string{},
	},
	Bus{
		"30",
		[]string{
			"Jayanagar 5th Block", "Trinity Circle", "JP Nagar",
		},
		[]string{},
	},
	Bus{
		"32",
		[]string{
			"Jayanagar 5th Block", "Trinity Circle", "Brigade Road", "KR Road", "JP Nagar",
		},
		[]string{},
	},
}
