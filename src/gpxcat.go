package main
 
import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)
 
type Gpx struct {
	Trk Track `xml:"trk"`
}

func (gpx *Gpx) Print() {
	fmt.Println(`<?xml version="1.0" encoding="utf-8" standalone="no"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1"
	creator="gpxcat" version="0.1"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd">`)
	gpx.Trk.Print()
	fmt.Println(`</gpx>`)
}

type Track struct {
	Name string `xml:"name"`
	Segments []Segment `xml:"trkseg"`
}

func (track *Track) Print() {
	fmt.Println(`	<trk>`)
	for _, segment := range track.Segments {
		segment.Print()
	}
	fmt.Println(`	</trk>`)
}

type Segment struct {
	Points []Point `xml:"trkpt"`
}

func (segment *Segment) Print() {
	fmt.Println(`		<trkseg>`)
	for _, point := range segment.Points {
		point.Print()
	}
	fmt.Println(`		</trkseg>`)
}

type Point struct {
	Latitude float64 `xml:"lat,attr"`
	Longitude float64 `xml:"lon,attr"`
	Elevation float32 `xml:"ele"`
	Time string `xml:"time"`
}

func (point *Point) Print() {
	fmt.Printf(`			<trkpt lat="%.10f" lon="%.10f">
				<ele>%.2f</ele>
				<time>%s</time>
			</trkpt>
`, point.Latitude, point.Longitude, point.Elevation, point.Time)
}

func LoadGPX(name string) (gpx Gpx, err error) {
	xmlFile, err := os.Open(name)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return gpx, err
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(b, &gpx)
	return gpx, nil
}
 
func main() {
	var cat Gpx
	for _, name := range os.Args[1:] {
		q, err := LoadGPX(name)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		cat.Trk.Segments = append(cat.Trk.Segments, q.Trk.Segments...)
	}
	cat.Print()
}
