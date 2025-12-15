package ds

import (
	"fmt"
	"rogue_game/internal/config"
	"testing"
)

func TestEdgesSort(t *testing.T) {
	w := config.RoomsInWidth
	l := config.RoomsInLength
	edges := RandomSortedEdges(w, l)
	for _, e := range edges {
		fmt.Printf("N1 %d, N2 %d, Length=%d\n", e.N1, e.N2, e.Length)

	}
	mst := BuildMST(w, l)
	for _, e := range mst {
		fmt.Println(e)
	}
}
