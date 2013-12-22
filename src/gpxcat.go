package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
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

func fatal(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %d", err)
		os.Exit(-1)
	}
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
}

func (point *Point) Print(wrt io.Writer) {
	if point.Valid {
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

func main() {
	var cat Gpx
	for _, name := range os.Args[1:] {
		q, err := LoadGPX(name)
		if err != nil {
			log.Fatalf("Error loading GPX file: %d", err)
			return
		}
		cat.Trk.Segments = append(cat.Trk.Segments, q.Trk.Segments...)
	}

	// unique
	unique := make(map[string]Segment)
	for _, segment := range cat.Trk.Segments {
		unique[segment.Points[0].Time] = segment
	}
	cat.Trk.Segments = make([]Segment, 0)
	for _, segment := range unique {
		cat.Trk.Segments = append(cat.Trk.Segments, segment)
	}

	// desegmentize
	segments := make([]Segment, 1)
	for _, segment := range cat.Trk.Segments {
		segments[0].Points = append(segments[0].Points, segment.Points...)
	}
	cat.Trk.Segments = segments

	// parse dates && validate
	points := cat.Trk.Segments[0].Points
	for index := range points {
		i_time, i_err := time.Parse("2006-01-02T15:04:05Z", points[index].Time)
		points[index].GoTime = i_time
		points[index].Valid = i_err == nil
	}

	// sort
	sort.Sort(ByDate(cat.Trk.Segments[0].Points))

	// by date
	gpxMap := make(map[string]Gpx)
	for _, point := range cat.Trk.Segments[0].Points {
		p_time := point.GoTime
		p_time.Add(time.Hour * -7)
		if point.Valid {
			date_str := p_time.Format("02-01-2006")
			track := gpxMap[date_str]
			if len(track.Trk.Segments) == 0 {
				track.Trk.Segments = make([]Segment, 1)
			}
			track.Trk.Segments[0].Points = append(track.Trk.Segments[0].Points, point)
			gpxMap[date_str] = track
		}
	}

	// output to files
	for date_str, track := range gpxMap {
		file, err := os.Create(fmt.Sprintf("out/%s.gpx", date_str))
		defer file.Close()
		fatal(err)
		track.Print(file)
	}

	//cat.Print()
}
