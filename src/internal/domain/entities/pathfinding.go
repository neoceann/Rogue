package entities

// AllTrays return list of all possible paths from starting point
// a. no cycles
// b. paths are partial (A, AB, AC, ABD, ABF...)
func (gs *GameSession) AllTrays() [][]int {
	edges := make([][]int, 0, 9)
	trays := [][]int{}
	for _, c := range gs.Corridors {
		e := []int{c.FromRoomInd, c.ToRoomInd}
		edges = append(edges, e)
	}
	trays = append(trays, []int{gs.Player.RoomInd})
	for j := 0; j < len(trays); j++ {
		t_end := trays[j][len(trays[j])-1]
		for i := 0; i < len(edges); i++ {
			a := edges[i][0]
			b := edges[i][1]
			if a == t_end || b == t_end {
				trays = append(trays, trays[j])
				if a == t_end {
					a = b
				}
				trays[len(trays)-1] = append(trays[len(trays)-1], a)
				edges = append(edges[:i], edges[i+1:]...)
				i--
			}
		}

		if len(edges) == 0 {
			break
		}
	}
	return trays
}
