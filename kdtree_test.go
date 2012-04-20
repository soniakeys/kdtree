package kdtree

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

// Wikipedia example data
func TestWP2D(t *testing.T) {
	kd := New([]Point{{2, 3}, {5, 4}, {9, 6}, {4, 7}, {8, 1}, {7, 2}},
		HyperRect{Point{0, 0}, Point{10, 10}})
	p := Point{9, 2}
	nn, ssq, nv := kd.Nearest(p)
	if p.Sqd(nn) != ssq {
		t.Error("nn, ssq results inconsistent")
	}
	if len(nn) != 2 || nn[0] != 8 || nn[1] != 1 {
		t.Error("Expected nn =", Point{8, 1}, "found", nn)
	}
	if math.Abs(ssq-2) > 1e14 {
		t.Error("Expected distance^2 =", 2, "found", ssq)
	}
	if nv != 3 {
		t.Error("Expected 3 nodes visited.  actual:", nv)
	}
}

// 1000 random 3d points
func TestRandom3D(t *testing.T) {
	rand.Seed(time.Now().Unix())
	pts := randomPts(3, 1000)
	kd := New(pts, HyperRect{Point{0, 0, 0}, Point{1, 1, 1}})
	p := randomPt(3)
	nn, ssq, nv := kd.Nearest(p)
	if p.Sqd(nn) != ssq {
		t.Error("nn, ssq results inconsistent")
	}
	if nv > 500 {
		t.Error("Expected nv << 1000, found nv =", nv)
	}
	// just check distance to random points for a while and make sure
	// none are closer that nn result.
	for a := time.After(time.Millisecond * 100); ; {
		select {
		case <-a:
			return
		default:
			pr := pts[rand.Intn(len(pts))]
			if p.Sqd(pr) < ssq {
				t.Logf("nn result (%v) not nearest to (%v).  ssq was %f",
					nn, p, ssq)
				t.Fatal("found", pr, "at sqd", p.Sqd(pr))
			}
		}
	}
}

func randomPt(dim int) Point {
	p := make(Point, dim)
	for d := range p {
		p[d] = rand.Float64()
	}
	return p
}

func randomPts(dim, n int) []Point {
	p := make([]Point, n)
	for i := range p {
		p[i] = randomPt(dim)
	}
	return p
}
