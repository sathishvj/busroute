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
		map[string]string{
			"1": "MG Road", "2": "Trinity Circle", "3": "Brigade Road",
		},
		map[string]string{},
	},
	Bus{
		"30",
		map[string]string{
			"1": "Jayanagar 5th Block", "2": "Trinity Circle", "3": "JP Nagar",
		},
		map[string]string{},
	},
	Bus{
		"32",
		map[string]string{
			"1": "Jayanagar 5th Block", "2": "Trinity Circle", "3": "Brigade Road", "4": "KR Road", "5": "JP Nagar",
		},
		map[string]string{},
	},
}
