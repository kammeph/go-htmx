package data

var Gearboxes = []Gearbox{
	{
		ID:       1,
		Serial:   "1234567",
		Type:     "RG 75",
		Backlash: 0.1,
		Housing:  "H123",
		Polygon:  "P125",
		Gear:     "G123",
	},
	{
		ID:       2,
		Serial:   "7412589",
		Type:     "RG 175",
		Backlash: 0.3,
		Housing:  "H126",
		Polygon:  "P127",
		Gear:     "G126",
	},
	{
		ID:       3,
		Serial:   "9632587",
		Type:     "RG 135",
		Backlash: -0.1,
		Housing:  "H125",
		Polygon:  "P124",
		Gear:     "G127",
	},
	{
		ID:       4,
		Serial:   "7532674",
		Type:     "RG 135",
		Backlash: -0.2,
		Housing:  "H127",
		Polygon:  "P123",
		Gear:     "G124",
	},
}

var Housings = []Housing{
	{
		ID:       1,
		Serial:   "H123",
		Size:     "Large",
		Type:     "Plastic",
		Measure1: 20.0,
		Measure2: 30.0,
	},
	{
		ID:       2,
		Serial:   "H124",
		Size:     "Medium",
		Type:     "Metal",
		Measure1: 21.0,
		Measure2: 31.0,
	},
	{
		ID:       3,
		Serial:   "H125",
		Size:     "Small",
		Type:     "Wood",
		Measure1: 22.0,
		Measure2: 32.0,
	},
	{
		ID:       4,
		Serial:   "H126",
		Size:     "Large",
		Type:     "Plastic",
		Measure1: 23.0,
		Measure2: 33.0,
	},
	{
		ID:       5,
		Serial:   "H127",
		Size:     "Medium",
		Type:     "Metal",
		Measure1: 24.0,
		Measure2: 34.0,
	},
}

var Polygons = []Polygon{
	{
		ID:       1,
		Serial:   "P123",
		Size:     "Small",
		Type:     "Hexagon",
		Measure1: 10.0,
		Measure2: 20.0,
	},
	{
		ID:       2,
		Serial:   "P124",
		Size:     "Medium",
		Type:     "Octagon",
		Measure1: 11.0,
		Measure2: 21.0,
	},
	{
		ID:       3,
		Serial:   "P125",
		Size:     "Large",
		Type:     "Decagon",
		Measure1: 12.0,
		Measure2: 22.0,
	},
	{
		ID:       4,
		Serial:   "P126",
		Size:     "Small",
		Type:     "Hexagon",
		Measure1: 13.0,
		Measure2: 23.0,
	},
	{
		ID:       5,
		Serial:   "P127",
		Size:     "Medium",
		Type:     "Octagon",
		Measure1: 14.0,
		Measure2: 24.0,
	},
}

var Gears = []Gear{
	{
		ID:       1,
		Serial:   "G123",
		Size:     "Medium",
		Type:     "Spur",
		Measure1: 15.0,
	},
	{
		ID:       2,
		Serial:   "G124",
		Size:     "Large",
		Type:     "Helical",
		Measure1: 16.0,
	},
	{
		ID:       3,
		Serial:   "G125",
		Size:     "Small",
		Type:     "Bevel",
		Measure1: 17.0,
	},
	{
		ID:       4,
		Serial:   "G126",
		Size:     "Medium",
		Type:     "Spur",
		Measure1: 18.0,
	},
	{
		ID:       5,
		Serial:   "G127",
		Size:     "Large",
		Type:     "Helical",
		Measure1: 19.0,
	},
}
