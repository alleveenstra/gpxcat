package main

import (
	"flag"
	"gpx"
	"log"
	"os"
	"sort"
	"time"
)

var Sort bool
var Split bool
var SplitFormat string
var SplitOutput string
var Simplify bool
var Resample bool
var Statistics bool
var MergeSegments bool
var Keep int

func init() {
	flag.BoolVar(&Sort, "sort", false, "sort the segments (default for merge, split & resample)")
	flag.BoolVar(&Split, "split", false, "split the track according to a format")
	flag.StringVar(&SplitFormat, "split-format", "2006-01-02", "the format of tracknames")
	flag.StringVar(&SplitOutput, "split-output", "out", "the directory to output tracks (must exist)")
	flag.BoolVar(&Simplify, "simplify", false, "simplify the track")
	flag.BoolVar(&Resample, "resample", false, "resample the track")
	flag.BoolVar(&Statistics, "statistics", false, "show statistics")
	flag.BoolVar(&MergeSegments, "merge-segments", false, "merge all segments into one bigger segment")
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
	if Sort || MergeSegments || Simplify || Resample || Split {
		cat.Trk.Segments = gpx.Desegmentize(cat.Trk.Segments)
	}

	// parse dates && validate
	for segment_index := range cat.Trk.Segments {
		points := cat.Trk.Segments[segment_index].Points
		for point_index := range points {
			i_time, i_err := time.Parse("2006-01-02T15:04:05Z", points[point_index].Time)
			points[point_index].GoTime = i_time
			points[point_index].Valid = i_err == nil
			points[point_index].Removed = false
		}
	}

	// sort
	if Sort || MergeSegments || Simplify || Resample || Split {
		sort.Sort(gpx.ByDate(cat.Trk.Segments[0].Points))
	}

	if Simplify {
		gpx.Simplify(cat.Trk.Segments, Keep)
	}

	if Resample {
		gpx.Resample(cat.Trk.Segments, Keep)
	}

	if Split {
		err := gpx.Split(cat, SplitOutput, SplitFormat)
		fatal(err)
	} else {
		cat.Print(os.Stdout)
	}
}
