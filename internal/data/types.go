package data

type Gearbox struct {
	ID       int64
	Serial   string
	Type     string
	Backlash float64
	Housing  string
	Polygon  string
	Gear     string
}

type Gear struct {
	ID       int64
	Serial   string
	Size     string
	Type     string
	Measure1 float64
}

type Housing struct {
	ID       int64
	Serial   string
	Size     string
	Type     string
	Measure1 float64
	Measure2 float64
}

type Polygon struct {
	ID       int64
	Serial   string
	Size     string
	Type     string
	Measure1 float64
	Measure2 float64
}
