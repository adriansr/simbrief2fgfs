package fgfs

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"

	"github.com/adriansr/simbrief2fgfs/go-sb2fgfs/ofp"
)

type FlightPlan struct {
	Version int
	Departure Airport
	Destination Airport
	Route struct {
		Wp []interface{}
	}
}

type Airport struct {
	Airport string
	Runway  string
}

type Runway struct {
	Departure bool
	Arrival   bool
	ICAO string
	Type string
	Ident string
}

type Waypoint struct {
	AltitudeFt float64 `xml:"altitude-ft"`
	AltRestrict string `xml:"alt-restrict"`
	Type string
	Lat float64
	Lon float64
}

func New(fp ofp.Flightplan) FlightPlan {
	plist := FlightPlan{
		Version:     2,
		Departure:   Airport{
			Airport: fp.Origin.ICAOCode,
			Runway:  fp.Origin.PlanRunway,
		},
		Destination: Airport{
			Airport: fp.Destination.ICAOCode,
			Runway:  fp.Destination.PlanRunway,
		},
	}
	plist.Route.Wp = append(plist.Route.Wp,
		Runway{
			Departure: true,
			ICAO:      fp.Origin.ICAOCode,
			Type:      "runway",
			Ident:     fp.Origin.PlanRunway,
		})
	for _, fix := range fp.NavLog {
		plist.Route.Wp = append(plist.Route.Wp, Waypoint{
			AltitudeFt:  float64(fix.Altitude),
			AltRestrict: "at",
			Type:        "navaid",
			Lat:         fix.Latitude,
			Lon:         fix.Longitude,
		})
	}
	plist.Route.Wp = append(plist.Route.Wp,
		Runway{
			Arrival: true,
			ICAO:      fp.Destination.ICAOCode,
			Type:      "runway",
			Ident:     fp.Destination.PlanRunway,
		})
	return plist
}

func marshalValue(encoder *xml.Encoder, name string, val reflect.Value) (err error) {
	switch val.Kind() {
	case reflect.Ptr, reflect.Interface:
		if !val.IsValid() {
			return nil
		}
		return marshalValue(encoder, name, val.Elem())

	case reflect.Struct:
		if err = encoder.EncodeToken(xml.StartElement{
			Name: xml.Name{Local: name},
		}); err != nil {
			return err
		}
		for i := 0; i < val.NumField(); i++ {
			fType := val.Type().Field(i)
			fldName := strings.ToLower(fType.Name)
			if tag, ok := fType.Tag.Lookup("xml"); ok {
				fldName = tag
			}
			if err := marshalValue(encoder, fldName, val.Field(i)); err != nil {
				return err
			}
		}
		if err = encoder.EncodeToken(xml.EndElement{
			Name: xml.Name{
				Local: name,
			},
		}); err != nil {
			return err
		}

	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if err = marshalValue(encoder, name, val.Index(i)); err != nil {
				return err
			}
		}

	default:
		if val.IsZero() {
			break
		}
		if err = encoder.EncodeToken(xml.StartElement{
			Name: xml.Name{Local: name},
			Attr: []xml.Attr{
				{
					Name: xml.Name{Local: "type"},
					Value: val.Type().Name(),
				},
			},
		}); err != nil {
			return err
		}
		if err = encoder.EncodeToken(xml.CharData(fmt.Sprintf("%v", val.Interface()))); err != nil {
			return err
		}
		if err = encoder.EncodeToken(xml.EndElement{Name: xml.Name{
			Local: name,
		}}); err != nil {
			return err
		}
	}
	return nil
}

func (fp *FlightPlan) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	if err := marshalValue(encoder, "PropertyList", reflect.ValueOf(fp)); err != nil {
		return err
	}
	return nil
}