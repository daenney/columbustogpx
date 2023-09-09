package main

import (
	"encoding/xml"
	"time"
)

type GPX11 struct {
	XMLName      xml.Name `xml:"gpx"`
	XMLNs        string   `xml:"xmlns,attr"`
	XmlNsXsi     string   `xml:"xmlns:xsi,attr"`
	XmlSchemaLoc string   `xml:"xsi:schemaLocation,attr"`

	Version string `xml:"version,attr"`
	Creator string `xml:"creator,attr"`

	Waypoints []Point `xml:"wpt,omitempty"`
	Tracks    []Track `xml:"trk,omitempty"`
}

func (g *GPX11) AddTrackSegment(p Point) {
	if len(g.Tracks) == 0 {
		g.Tracks = []Track{
			{
				TrackSegment: []TrackSegment{
					{
						[]Point{p},
					},
				},
			},
		}
		return
	}
	g.Tracks[0].TrackSegment[0].TrackPoints = append(g.Tracks[0].TrackSegment[0].TrackPoints, p)
}

func (g *GPX11) AddPOI(p Point) {
	g.Waypoints = append(g.Waypoints, p)
}

func NewGPX() GPX11 {
	return GPX11{
		XMLNs:        "https://www.topografix.com/GPX/1/1",
		XmlNsXsi:     "https://www.w3.org/2001/XMLSchema-instance",
		XmlSchemaLoc: "https://www.topografix.com/GPX/1/1 https://www.topografix.com/GPX/1/1/gpx.xsd",
		Version:      "1.1",
		Creator:      "ColumbusToGPX",
	}
}

type Point struct {
	Latitude  string    `xml:"lat,attr"`
	Longitude string    `xml:"lon,attr"`
	Elevation string    `xml:"ele"`
	Timestamp time.Time `xml:"time"`
}

type Track struct {
	XMLName      xml.Name       `xml:"trk"`
	TrackSegment []TrackSegment `xml:"trkseg,omitempty"`
}

type TrackSegment struct {
	TrackPoints []Point `xml:"trkpt,omitempty"`
}
