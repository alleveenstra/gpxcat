package main

import (
	"flag"
	"gpx"
	"log"
	"os"
	"sort"
	"time"
)

var Keep int
var Simplify bool
var Split bool
var Resample bool
var Out string
var Format string

func init() {
	flag.StringVar(&Format, "format", "2006-01-02", "the format of tracknames")
	flag.StringVar(&Out, "out", "out", "the directory to output tracks (must exist)")
	flag.BoolVar(&Split, "split", false, "split the track according to a format")
	flag.BoolVar(&Simplify, "simplify", false, "simplify the track")
	flag.BoolVar(&Resample, "resample", false, "resample the track")
	flag.IntVar(&Keep, "keep", 200, "number of points to keep")
}

func fatal(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %d", err)
		os.Exit(-1)
	}
}

func main() {
	flag.Parse()

	var cat gpx.Gpx
	for _, name := range flag.Args() {
		q, err := gpx.LoadGPX(name)
		if err != nil {
			log.Fatalf("Error loading GPX file: %d", err)
			return
		}
		cat.Trk.Segments = append(cat.Trk.Segments, q.Trk.Segments...)
	}

	// unique
	cat.Trk.Segments = gpx.Unique(cat.Trk.Segments)

	// desegmentize
	cat.Trk.Segments = gpx.Desegmentize(cat.Trk.Segments)

	// parse dates && validate
	points := cat.Trk.Segments[0].Points
	for index := range points {
		i_time, i_err := time.Parse("2006-01-02T15:04:05Z", points[index].Time)
		points[index].GoTime = i_time
		points[index].Valid = i_err == nil
		points[index].Removed = false
	}

	// sort
	sort.Sort(gpx.ByDate(cat.Trk.Segments[0].Points))

	if Simplify {
		gpx.Simplify(cat.Trk.Segments, Keep)
	}

	if Resample {
		gpx.Resample(cat.Trk.Segments, Keep)
	}

	if Split {
		err := gpx.Split(cat, Out, Format)
		fatal(err)
	} else {
		cat.Print(os.Stdout)
	}
}
