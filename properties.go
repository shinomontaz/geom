package geomop

import (
	"github.com/twpayne/gogeom/geom"
	"math"
)

const tolerance = 1.e-9

// Function Area returns the area of a polygon, or the combined area of a
// MultiPolygon, assuming that none of the polygons in the MultiPolygon
// overlap and that nested polygons have alternating winding directions.
func Area(g geom.T) float64 {
	a := 0.
	switch g.(type) {
	case geom.Polygon:
		for _, r := range g.(geom.Polygon).Rings {
			a += area(r)
		}
	case geom.MultiPolygon:
		for _, p := range g.(geom.MultiPolygon).Polygons {
			a += Area(p)
		}
	case geom.GeometryCollection:
		for _, g := range g.(geom.GeometryCollection).Geoms {
			a += Area(g)
		}
	}
	return math.Abs(a)
}

// Function Length returns the length of a LineString, or the combined
// length of a MultiLineString.
func Length(g geom.T) float64 {
	l := 0.
	switch g.(type) {
	case geom.LineString:
		l = length(g.(geom.LineString).Points)
	case geom.MultiLineString:
		for _, line := range g.(geom.MultiLineString).LineStrings {
			l += Length(line)
		}
	case geom.GeometryCollection:
		for _, g := range g.(geom.GeometryCollection).Geoms {
			l += Length(g)
		}
	}
	return l
}

// see http://www.mathopenref.com/coordpolygonarea2.html
func area(polygon []geom.Point) float64 {
	highI := len(polygon) - 1
	A := (polygon[highI].X +
		polygon[0].X) * (polygon[0].Y - polygon[highI].Y)
	for i := 0; i < highI; i++ {
		A += (polygon[i].X +
			polygon[i+1].X) * (polygon[i+1].Y - polygon[i].Y)
	}
	return A / 2.
}

func length(line []geom.Point) float64 {
	l := 0.
	for i := 0; i < len(line)-1; i++ {
		p1 := line[i]
		p2 := line[i+1]
		l += math.Hypot(p2.X-p1.X, p2.Y-p1.Y)
	}
	return l
}

// orientation2D_Polygon(): test the orientation of a simple 2D polygon
//  Input:  Point* V = an array of n+1 vertex points with V[n]=V[0]
//  Return: >0 for counterclockwise
//          =0 for none (degenerate)
//          <0 for clockwise
//  Note: this algorithm is faster than computing the signed area.
//  From http://geomalgorithms.com/a01-_area.html#orientation2D_Polygon()
func orientation(V geom.Polygon) []float64 {
	// first find rightmost lowest vertex of the polygon
	out := make([]float64, len(V.Rings))
	for j, r := range V.Rings {
		rmin := 0
		xmin := r[0].X
		ymin := r[0].Y
		for i, p := range r {
			if p.Y > ymin {
				continue
			} else if p.Y == ymin { // just as low
				if p.X < xmin { // and to left
					continue
				}
			}
			rmin = i // a new rightmost lowest vertex
			xmin = p.X
			ymin = p.Y
		}

		// test orientation at the rmin vertex
		// ccw <=> the edge leaving V[rmin] is left of the entering edge
		if rmin == 0 || rmin == len(r)-1 {
			out[j] = isLeft(r[len(r)-2], r[0], r[1])
		} else {
			out[j] = isLeft(r[rmin-1], r[rmin], r[rmin+1])
		}
	}
	return out
}

// isLeft(): test if a point is Left|On|Right of an infinite 2D line.
//    Input:  three points P0, P1, and P2
//    Return: >0 for P2 left of the line through P0 to P1
//          =0 for P2 on the line
//          <0 for P2 right of the line
//    From http://geomalgorithms.com/a01-_area.html#isLeft()
func isLeft(P0, P1, P2 geom.Point) float64 {
	return ((P1.X-P0.X)*(P2.Y-P0.Y) -
		(P2.X-P0.X)*(P1.Y-P0.Y))
}

// Change the winding direction of the outer and inner
// rings so the outer ring is counter-clockwise and
// nesting rings alternate directions.
func FixOrientation(g geom.T) {
	p := g.(geom.Polygon)
	o := orientation(p)
	for i, inner := range p.Rings {
		numInside := 0
		for j, outer := range p.Rings {
			if i != j {
				if polyInPoly(Contour(outer), Contour(inner)) {
					numInside++
				}
			}
		}
		if numInside%2 == 1 && o[i] > 0. {
			reversePolygon(inner)
		} else if numInside%2 == 0 && o[i] < 0. {
			reversePolygon(inner)
		}
	}
}

func reversePolygon(s []geom.Point) []geom.Point {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func polyInPoly(outer, inner Contour) bool {
	for _, p := range inner {
		if !outer.Contains(p) {
			return false
		}
	}
	return true
}

func Within(inner, outer geom.T) bool {
	switch outer.(type) {
	case geom.Polygon:
		op := outer.(geom.Polygon)
		switch inner.(type) {
		case geom.Polygon:
			ip := inner.(geom.Polygon)
			for _, r := range ip.Rings {
				for _, p := range r {
					if !PointInPolygon(p, op) {
						return false
					}
				}
			}
			return true
		default:
			panic(NewError(outer))
			return false
		}
	default:
		panic(NewError(outer))
		return false
	}
}

// Function PointInPolygon determines whether "point" is
// within "polygon". If "polygon" is not actually a polygon,
// return false.
func PointInPolygon(point geom.Point, polygon geom.T) bool {
	inCount := 0
	switch polygon.(type) {
	case geom.Polygon:
		o := orientation(polygon.(geom.Polygon))
		for i, r := range polygon.(geom.Polygon).Rings {
			if Contour(r).Contains(point) {
				if o[i] > 0. {
					inCount++
				} else if o[i] < 0. {
					inCount--
				}
			}
		}
		return inCount > 0
	case geom.MultiPolygon:
		for _, pp := range polygon.(geom.MultiPolygon).Polygons {
			if PointInPolygon(point, geom.T(pp)) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// dot product
func dot(u, v geom.Point) float64 { return u.X*v.X + u.Y*v.Y }

// norm = length of  vector
func norm(v geom.Point) float64 { return math.Sqrt(dot(v, v)) }

// distance = norm of difference
func d(u, v geom.Point) float64 { return norm(pointSubtract(u, v)) }

// dist_Point_to_Segment(): get the distance of a point to a segment
//     Input:  a Point P and a Segment S (in any dimension)
//     Return: the shortest distance from P to S
// from http://geomalgorithms.com/a02-_lines.html
func distPointToSegment(p, segStart, segEnd geom.Point) float64 {
	v := pointSubtract(segEnd, segStart)
	w := pointSubtract(p, segStart)

	c1 := dot(w, v)
	if c1 <= 0. {
		return d(p, segStart)
	}

	c2 := dot(v, v)
	if c2 <= c1 {
		return d(p, segEnd)
	}

	b := c1 / c2
	pb := geom.Point{segStart.X + b*v.X, segStart.Y + b*v.Y}
	return d(p, pb)
}

func pointOnSegment(p, segStart, segEnd geom.Point) bool {
	return distPointToSegment(p, segStart, segEnd) < tolerance
}
