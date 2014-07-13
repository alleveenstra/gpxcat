package gpx

import (
	"fmt"
	"os"
	"time"
)

func Split(cat Gpx, out, format string) error {
	// by date
	gpxMap := make(map[string]Gpx)
	for _, point := range cat.Trk.Segments[0].Points {
		p_time := point.GoTime
		p_time = p_time.Add(time.Hour * time.Duration(int(point.Longitude/15.0)))
		if point.Valid {
			date_str := p_time.Format(format)
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
		file, err := os.Create(fmt.Sprintf("%s/%s.gpx", out, date_str))
		defer file.Close()
		if err != nil {
			return err
		}
		track.Print(file)
	}
	return nil
}
