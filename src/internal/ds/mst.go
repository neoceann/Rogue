package ds

import (
	"math/rand"
	"sort"
)

// MST for connecting all rooms in random order
// col
// 1 2 3  r
// 4 5 6  o
// 7 8 9  w

type Edge struct {
	N1, N2 int
	Length int
}

// all edges sorted by len ascending
// for randomization length of nearest cells is in [7, 9]
// maybe not add rand just shuffle array of edges
func RandomSortedEdges(width, length int) []Edge {
	edges := []Edge{}
	for i := range length {
		for j := range width {
			n := i*width + j
			if j < width-1 {
				e := Edge{N1: n, N2: n + 1, Length: 7 + rand.Intn(3)}
				edges = append(edges, e)
			}
			if i < length-1 {
				e := Edge{N1: n, N2: n + width, Length: 7 + rand.Intn(3)}
				edges = append(edges, e)
			}
		}
	}
	sort.Slice(edges, func(k, n int) bool {
		return edges[k].Length < edges[n].Length
	})
	return edges
}

// keeps parent of the member
type UnionParent struct {
	parent []int // ancistor, highest iin the union
}

// member = parent, returns pointer
func NewUnionParent(width, length int) *UnionParent {
	parent := make([]int, width*length)
	for i := range parent {
		parent[i] = i
	}
	return &UnionParent{parent: parent}
}

// recursive call till find parent (member == parent)
func (up *UnionParent) FindParent(x int) int {
	if up.parent[x] != x {
		up.parent[x] = up.FindParent(up.parent[x])
	}
	return up.parent[x]
}

// if x and y not in the same union, then redefines
// parent of x to parent of y
func (up *UnionParent) Union(x, y int) {
	pX, pY := up.FindParent(x), up.FindParent(y)
	if pX != pY {
		up.parent[pX] = up.parent[pY]
	}

}

func BuildMST(width, length int) []Edge {
	edges := RandomSortedEdges(width, length)
	up := NewUnionParent(width, length)
	var mst []Edge
	for _, e := range edges {
		if up.FindParent(e.N1) != up.FindParent(e.N2) {
			up.Union(e.N1, e.N2)
			mst = append(mst, e)
			if len(mst) == width*length-1 {
				break
			}
		}
	}
	return mst
}
