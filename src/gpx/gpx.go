package gpx

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Gpx struct {
	Trk Track `xml:"trk"`
}

func (gpx *Gpx) Print(wrt io.Writer) {
	fmt.Fprintf(wrt, `<?xml version="1.0" encoding="utf-8" standalone="no"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1"
	creator="gpxcat" version="0.1"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd">
`)
	gpx.Trk.Print(wrt)
	fmt.Fprintf(wrt, `</gpx>
`)
}

type Track struct {
	Name     string    `xml:"name"`
	Segments []Segment `xml:"trkseg"`
}

type ByDate []Point

func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool {
	if !a[i].Valid || !a[j].Valid {
		return true
	}
	return a[j].GoTime.After(a[i].GoTime)
}

func (track *Track) Print(wrt io.Writer) {
	fmt.Fprintf(wrt, `	<trk>
`)
	for _, segment := range track.Segments {
		segment.Print(wrt)
	}
	fmt.Fprintf(wrt, `	</trk>
`)
}

type Segment struct {
	Points []Point `xml:"trkpt"`
}

func (segment *Segment) Print(wrt io.Writer) {
	fmt.Fprintf(wrt, `		<trkseg>
`)
	for _, point := range segment.Points {
		point.Print(wrt)
	}
	fmt.Fprintf(wrt, `		</trkseg>
`)
}

type Point struct {
	Latitude  float64 `xml:"lat,attr"`
	Longitude float64 `xml:"lon,attr"`
	Elevation float32 `xml:"ele"`
	Time      string  `xml:"time"`
	GoTime    time.Time
	Valid     bool
	Removed   bool
}

func (point *Point) Print(wrt io.Writer) {
	if point.Valid && !point.Removed {
		fmt.Fprintf(wrt, `			<trkpt lat="%.10f" lon="%.10f">
				<ele>%.2f</ele>
				<time>%s</time>
			</trkpt>
`, point.Latitude, point.Longitude, point.Elevation, point.Time)
	}
}

func LoadGPX(name string) (gpx Gpx, err error) {
	xmlFile, err := os.Open(name)
	if err != nil {
		return gpx, err
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(b, &gpx)
	log.Printf("reading %s", name)
	return gpx, nil
}
