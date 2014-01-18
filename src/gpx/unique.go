package gpx

func Unique(segments []Segment) []Segment {
	unique := make(map[string]Segment)
	for _, segment := range segments {
		unique[segment.Points[0].Time] = segment
	}
	newSegments := make([]Segment, 0)
	for _, segment := range unique {
		newSegments = append(newSegments, segment)
	}
	return newSegments
}
