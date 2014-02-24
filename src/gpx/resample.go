package gpx

import (
	"math"
)

func Resample(segments []Segment, keep int) {
	if keep < 3 {
		keep = 3
	}
	for segment_id := range segments {
		length := segments[segment_id].Length()
		every := length / float64(keep)
		accum := 0.0
		prev_point := -1
		for point_id := range segments[segment_id].Points {
			if prev_point != -1 {
				accum += segments[segment_id].Points[prev_point].distance(segments[segment_id].Points[point_id])
				if accum > every {
					segments[segment_id].Points[point_id].Removed = false
					accum = 0
				} else {
					segments[segment_id].Points[point_id].Removed = true
				}
			}
			prev_point = point_id
		}
		segments[segment_id].Points[0].Removed = false
		segments[segment_id].Points[len(segments[segment_id].Points)-1].Removed = false
	}
}

func (segment *Segment) Length() float64 {
	length := 0.0
	prev_point := -1
	for point_id := range segment.Points {
		if prev_point != -1 {
			length += segment.Points[prev_point].distance(segment.Points[point_id])
		}
		prev_point = point_id
	}
	return length
}

func (a *Point) distance(b Point) (d float64) {
	lat1 := a.Latitude
	lon1 := a.Longitude
	lat2 := b.Latitude
	lon2 := b.Longitude
	if lat1 == lat2 && lon1 == lon2 {
		return 0.0
	}
	theta := lon1 - lon2
	dist := sin(deg2rad(lat1))*sin(deg2rad(lat2)) + cos(deg2rad(lat1))*cos(deg2rad(lat2))*cos(deg2rad(theta))
	dist = acos(dist)
	dist = rad2deg(dist)
	if math.IsNaN(dist) {
		return 0.0
	}
	return dist * 111.18957696
}

func deg2rad(d float64) float64 {
	return d * math.Pi / 180.0
}

func rad2deg(d float64) float64 {
	return d * 180.0 / math.Pi
}

func sin(d float64) float64 {
	return math.Sin(d)
}

func cos(d float64) float64 {
	return math.Cos(d)
}

func acos(d float64) float64 {
	return math.Acos(d)
}
