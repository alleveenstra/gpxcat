package gpx

func Desegmentize(segments []Segment) []Segment {
	newSegments := make([]Segment, 1)
	for _, segment := range segments {
		newSegments[0].Points = append(newSegments[0].Points, segment.Points...)
	}
	return newSegments
}
