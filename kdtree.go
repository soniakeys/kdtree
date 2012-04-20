// Copyright 2012 Sonia Keys
// License MIT: http://www.opensource.org/licenses/MIT

// K-d tree example implementation.
//
// Implmentation follows pseudocode from "An intoductory tutorial on kd-trees"
// by Andrew W. Moore, Carnegie Mellon University, PDF accessed from
// http://www.autonlab.org/autonweb/14665
package kdtree

import (
	"math"
	"sort"
)

// Point is a k-dimensional point.
type Point []float64

// Sqd returns the square of the euclidean distance.
func (p Point) Sqd(q Point) float64 {
	var sum float64
	for dim, pCoord := range p {
		d := pCoord - q[dim]
		sum += d * d
	}
	return sum
}

// HyperRect is used to represent a k-dimensional bounding box.
type HyperRect struct {
	Min, Max Point
}

// Copy performs a deep copy, which is usually what you want.
//
// Go slices (the Point objects in a HyperRect) are reference objects.
// The data must be copied if you want to modify one without modifying
// the original.
func (hr HyperRect) Copy() HyperRect {
	return HyperRect{append(Point{}, hr.Min...), append(Point{}, hr.Max...)}
}

// KdTree represents a k-d tree and associated k-d bounding box.
type KdTree struct {
	n      *kdNode
	Bounds HyperRect
}

// kdNode following field names in the paper.
// rangeElt would be whatever data is associated with the point.
// we don't bother with it for this example.
type kdNode struct {
	domElt      Point
	split       int
	left, right *kdNode
}

// New constructs a KdTree from a list of points and a bounding box.
//
// The bounds could be computed of course, but typically you know them already.
func New(pts []Point, bounds HyperRect) KdTree {
	// algorithm is table 6.3 in the paper.
	var nk2 func([]Point, int) *kdNode
	nk2 = func(exset []Point, split int) *kdNode {
		if len(exset) == 0 {
			return nil
		}
		// pivot choosing procedure.  we find median, then find largest
		// index of points with median value.  this satisfies the
		// inequalities of steps 6 and 7 in the algorithm.
		sort.Sort(part{exset, split})
		m := len(exset) / 2
		d := exset[m]
		for m+1 < len(exset) && exset[m+1][split] == d[split] {
			m++
		}
		// next split
		s2 := split + 1
		if s2 == len(d) {
			s2 = 0
		}
		return &kdNode{d, split, nk2(exset[:m], s2), nk2(exset[m+1:], s2)}
	}
	return KdTree{nk2(pts, 0), bounds}
}

// Nearest.  find nearest neighbor.
//
// return values:
//  - nearest neighbor--the point within the tree that is nearest p.
//  - square of the distance to that point.
//  - a count of the nodes visited in the search.
func (t KdTree) Nearest(p Point) (best Point, bestSqd float64, nv int) {
	return nn(t.n, p, t.Bounds, math.Inf(1))
}

// algorithm is table 6.4 from the paper, with the addition of counting
// the number nodes visited.
func nn(kd *kdNode, target Point, hr HyperRect,
	maxDistSqd float64) (nearest Point, distSqd float64, nodesVisited int) {
	if kd == nil {
		return nil, math.Inf(1), 0
	}
	nodesVisited++
	s := kd.split
	pivot := kd.domElt
	leftHr := hr.Copy()
	rightHr := hr.Copy()
	leftHr.Max[s] = pivot[s]
	rightHr.Min[s] = pivot[s]
	targetInLeft := target[s] <= pivot[s]
	var nearerKd, furtherKd *kdNode
	var nearerHr, furtherHr HyperRect
	if targetInLeft {
		nearerKd, nearerHr = kd.left, leftHr
		furtherKd, furtherHr = kd.right, rightHr
	} else {
		nearerKd, nearerHr = kd.right, rightHr
		furtherKd, furtherHr = kd.left, leftHr
	}
	var nv int
	nearest, distSqd, nv = nn(nearerKd, target, nearerHr, maxDistSqd)
	nodesVisited += nv
	if distSqd < maxDistSqd {
		maxDistSqd = distSqd
	}
	d := pivot[s] - target[s]
	d *= d
	if d > maxDistSqd {
		return
	}
	if d = pivot.Sqd(target); d < distSqd {
		nearest = pivot
		distSqd = d
		maxDistSqd = distSqd
	}
	tempNearest, tempSqd, nv := nn(furtherKd, target, furtherHr, maxDistSqd)
	nodesVisited += nv
	if tempSqd < distSqd {
		nearest = tempNearest
		distSqd = tempSqd
	}
	return
}

// a container type used for sorting.  it holds the points to sort and
// the dimension to use for the sort key.
type part struct {
	pts   []Point
	dPart int
}

// satisfy sort.Interface
func (p part) Len() int { return len(p.pts) }
func (p part) Less(i, j int) bool {
	return p.pts[i][p.dPart] < p.pts[j][p.dPart]
}
func (p part) Swap(i, j int) { p.pts[i], p.pts[j] = p.pts[j], p.pts[i] }
