package generation

import (
	"fmt"
	"rogue_game/internal/domain/entities"
)

// findCenterOfBounds calculates Coordinates of center of Bounds
func findCenterOfBounds(b entities.Bounds) entities.Coordinates {
	x := (b.Pos0.X + b.Pos1.X) / 2
	y := (b.Pos0.Y + b.Pos1.Y) / 2
	return entities.Coordinates{X: x, Y: y}
}

// ake it away later when change key
// returns random position within Bounds inside Room by RoomInd
func (g *Generator) randomCoordinatesInsideRoom(gs *entities.GameSession,
	roomInd int) (entities.Coordinates, error) {
	x0 := gs.Rooms[roomInd].Bounds.Pos0.X
	x1 := gs.Rooms[roomInd].Bounds.Pos1.X
	xRange := x1 - x0
	y0 := gs.Rooms[roomInd].Bounds.Pos0.Y
	y1 := gs.Rooms[roomInd].Bounds.Pos1.Y
	yRange := y1 - y0
	x := g.rng.Intn(xRange) + x0
	y := g.rng.Intn(yRange) + y0
	if x < x0 || x >= x1 || y < y0 || y >= y1 {
		return entities.Coordinates{},
			fmt.Errorf("coordinates are outside of room ind=%d "+
				"(randomCoordinatesInsideRoom)", roomInd)
	}
	return entities.Coordinates{X: x, Y: y}, nil
}

func (g *Generator) randomCoordinatesInsideRoom1(r *entities.Room,
) (entities.Coordinates, error) {
	x0, x1 := r.Bounds.Pos0.X, r.Bounds.Pos1.X
	y0, y1 := r.Bounds.Pos0.Y, r.Bounds.Pos1.Y
	xRange := x1 - x0
	yRange := y1 - y0
	x := g.rng.Intn(xRange) + x0
	y := g.rng.Intn(yRange) + y0
	if x < x0 || x >= x1 || y < y0 || y >= y1 {
		return entities.Coordinates{},
			fmt.Errorf("coordinates are outside of room" +
				"(randomCoordinatesInsideRoom)")
	}
	return entities.Coordinates{X: x, Y: y}, nil
}

// randomShifted returns random point of line within its borders
// with shioft from edges
// return -1
func (g *Generator) randomShifted(len int) (int, error) {
	shift := 2 // can be adjusted
	if n := len - shift*2; n > 0 {
		return g.rng.Intn(len-shift*2) + shift, nil
	}
	return 0, fmt.Errorf("invalid shift=%dfor dancing corridors build", shift)
}

// check if coordinates are not busy with other already generated objects
func isPosFree(gs *entities.GameSession, pos entities.Coordinates) bool {
	for _, i := range gs.Items {
		if i.Pos == pos {
			return false
		}
	}
	for _, m := range gs.Monsters {
		if m.Pos == pos {
			return false
		}
	}
	for _, k := range gs.Keys {
		if k.Pos.X == pos.X && k.Pos.Y == pos.Y {
			return false
		}
	}
	if gs.Player.Pos == pos || gs.Finish == pos {
		return false
	}
	return true
}
