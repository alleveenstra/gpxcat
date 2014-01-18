package gpx

import (
	"math"
	"sort"
)

type removalInfo struct {
	id           int
	left         int
	right        int
	removalError float64
}

type byRemovalError []removalInfo

func (a byRemovalError) Len() int      { return len(a) }
func (a byRemovalError) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byRemovalError) Less(i, j int) bool {
	return a[j].removalError > a[i].removalError
}

func Simplify(segments []Segment, keep int) {
	if keep < 2 {
		keep = 2
	}
	for segment_id := range segments {
		pointsLeft := len(segments[segment_id].Points)
		for pointsLeft > keep {
			info := make([]removalInfo, 0)
			for point_id := range segments[segment_id].Points {
				if point_id != 0 && point_id < len(segments[segment_id].Points)-1 && !segments[segment_id].Points[point_id].Removed {
					left_point_id := prev_point(segments[segment_id].Points, point_id)
					left := segments[segment_id].Points[left_point_id]
					middle := segments[segment_id].Points[point_id]
					right_point_id := next_point(segments[segment_id].Points, point_id)
					right := segments[segment_id].Points[right_point_id]
					removalError := calcRemovalError(left, middle, right)
					pt := removalInfo{point_id, left_point_id, right_point_id, removalError}
					info = append(info, pt)
				}
			}
			sort.Sort(byRemovalError(info))
			removeN := pointsLeft - keep
			for _, info_point := range info[0:removeN] {
				if pointsLeft <= keep {
					break
				}
				if !segments[segment_id].Points[info_point.left].Removed && !segments[segment_id].Points[info_point.right].Removed {
					segments[segment_id].Points[info_point.id].Removed = true
					pointsLeft--
				}
			}
		}
	}
}

func prev_point(points []Point, middle int) int {
	for middle > 0 {
		middle -= 1
		if !points[middle].Removed {
			return middle
		}
	}
	return 0
}

func next_point(points []Point, middle int) int {
	for middle < len(points)-1 {
		middle += 1
		if !points[middle].Removed {
			return middle
		}
	}
	return len(points) - 1
}

func calcRemovalError(left, middle, right Point) float64 {
	lat1 := (middle.Latitude - left.Latitude) * 2
	lng1 := middle.Longitude - left.Longitude
	lat2 := (middle.Latitude - right.Latitude) * 2
	lng2 := middle.Longitude - right.Longitude
	len1 := math.Sqrt(math.Pow(lat1, 2.0) + math.Pow(lng1, 2.0))
	len2 := math.Sqrt(math.Pow(lat2, 2.0) + math.Pow(lng2, 2.0))
	// unitize
	lat1 /= len1
	lng1 /= len1
	lat2 /= len2
	lng2 /= len2
	// error = cosine
	return lat1*lat2 + lng1*lng2
}
