package ofp

import (
	"encoding/xml"
	"io/ioutil"
)

type Airport struct {
	ICAOCode string `xml:"icao_code"`
	PlanRunway string `xml:"plan_rwy"`
	Name string `xml:"name"`
	Elevation int `xml:"elevation"`
	Latitude float64 `xml:"pos_lat"`
	Longitude float64 `xml:"pos_long"`
}

type Fix struct {
	Ident string `xml:"ident"`
	Name  string `xml:"name"`
	Type  string `xml:"type"`
	Latitude float64 `xml:"pos_lat"`
	Longitude float64 `xml:"pos_long"`
	Stage string `xml:"stage"`
	IsSIDorSTAR bool `xml:"is_sid_star"`
	Distance float64 `xml:"distance"`
	Altitude int `xml:"altitude_feet"`
	KIAS int `xml:"ind_airspeed"`
	TAS  int `xml:"true_airspeed"`
	MACH float64 `xml:"mach"`
	GS int `xml:"groundspeed"`
}

// Flightplan represents an OFP flight plan.
type Flightplan struct {
	// General part of the SimBrief flightplan.
	General struct {
		// Release is ... ?
		Release string `xml:"release"`

		// AirlineCode is the ICAO 3-letter code for the airline.
		AirlineCode string `xml:"icao_airline"`

		// FlightNumber is the flight number. Together with the AirlineCode
		// it forms the flight code.
		FlightNumber string `xml:"flight_number"`
	} `xml:"general"`

	Origin      Airport `xml:"origin"`
	Destination Airport `xml:"destination"`
	NavLog      []Fix   `xml:"navlog>fix"`
}


func NewFlightPlanFromFile(path string) (Flightplan, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return Flightplan{}, err
	}
	return NewFlightPlanFromBytes(contents)
}

func NewFlightPlanFromBytes(contents []byte) (Flightplan, error) {
	var fp Flightplan
	if err := xml.Unmarshal(contents, &fp); err != nil {
		return Flightplan{}, err
	}
	return fp, nil
}
