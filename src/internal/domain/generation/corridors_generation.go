package generation

import (
	"fmt"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
)

// centeredCorridors build corridors between two neighbour rooms
// a. corridors are straight
// b. exits are in the middle of nearest walls
// returns error in case rooms are not neughbors horisontally or vertically
func (g *Generator) centeredCorridors(gs *entities.GameSession) error {
	for i := range gs.Corridors {
		indFrom, indTo, r1, r2 := g.orderRoomsFromToInc(gs, i)
		switch indTo - indFrom {
		// horisontal neighbours
		case 1:
			y := r1.Pos0.Y + (r1.Pos1.Y-r1.Pos0.Y)/2
			for x := r1.Pos1.X; x < r2.Pos0.X; x++ {
				gs.Corridors[i].Points = append(gs.Corridors[i].Points,
					entities.Coordinates{X: x, Y: y})
			}
		// vertical neighbours
		case config.RoomsInWidth:
			x := r1.Pos0.X + (r1.Pos1.X-r1.Pos0.X)/2
			for y := r1.Pos1.Y; y < r2.Pos0.Y; y++ {
				gs.Corridors[i].Points = append(gs.Corridors[i].Points,
					entities.Coordinates{X: x, Y: y})
			}
		// neither horisontal nor vertical neighbours
		default:
			return fmt.Errorf("corridor between rooms %d and %d,"+
				"not horisontally or vertically neighbours", indFrom, indTo)
		}
	}
	return nil
}

// dancingCorridors build corridors between two neighbour rooms
// a. corridors are curved
// b. exits are in random points of nearest walls
// returns error in case rooms are not neughbors horisontally or vertically
func (g *Generator) dancingCorridors(gs *entities.GameSession) error {
	for i := range gs.Corridors {
		indFrom, indTo, r1, r2 := g.orderRoomsFromToInc(gs, i)
		switch indTo - indFrom {
		// horisontal neighbours
		case 1:
			shift, err := g.randomShifted(r1.Pos1.Y - r1.Pos0.Y)
			if err != nil {
				return fmt.Errorf("dancivg corridors random shift "+
					"generation: %w", err)
			}
			y1 := r1.Pos0.Y + shift
			shift, err = g.randomShifted(r2.Pos1.Y - r2.Pos0.Y)
			if err != nil {
				return fmt.Errorf("dancivg corridors random shift "+
					"generation: %w", err)
			}
			y2 := r2.Pos0.Y + shift
			shift, err = g.randomShifted(r2.Pos0.X - r1.Pos1.X)
			if err != nil {
				return fmt.Errorf("dancivg corridors random shift "+
					"generation: %w", err)
			}
			xBetween := r1.Pos1.X + shift
			for x := r1.Pos1.X; x < xBetween; x++ {
				gs.Corridors[i].Points = append(gs.Corridors[i].Points,
					entities.Coordinates{X: x, Y: y1})
			}
			if y1 < y2 {
				for y := y1; y <= y2; y++ {
					gs.Corridors[i].Points = append(gs.Corridors[i].Points,
						entities.Coordinates{X: xBetween, Y: y})
				}
			} else {
				for y := y1; y >= y2; y-- {
					gs.Corridors[i].Points = append(gs.Corridors[i].Points,
						entities.Coordinates{X: xBetween, Y: y})
				}
			}
			for x := xBetween + 1; x < r2.Pos0.X; x++ {
				gs.Corridors[i].Points = append(gs.Corridors[i].Points,
					entities.Coordinates{X: x, Y: y2})
			}
		// vertical neighbours
		case config.RoomsInWidth:
			shift, err := g.randomShifted(r1.Pos1.X - r1.Pos0.X)
			if err != nil {
				return fmt.Errorf("dancivg corridors random shift "+
					"generation: %w", err)
			}
			x1 := r1.Pos0.X + shift
			shift, err = g.randomShifted(r2.Pos1.X - r2.Pos0.X)
			if err != nil {
				return fmt.Errorf("dancivg corridors random shift "+
					"generation: %w", err)
			}
			x2 := r2.Pos0.X + shift
			shift, err = g.randomShifted(r2.Pos0.Y - r1.Pos1.Y)
			if err != nil {
				return fmt.Errorf("dancivg corridors random shift "+
					"generation: %w", err)
			}
			yBetween := r1.Pos1.Y + shift
			for y := r1.Pos1.Y; y < yBetween; y++ {
				gs.Corridors[i].Points = append(gs.Corridors[i].Points,
					entities.Coordinates{X: x1, Y: y})
			}
			if x1 < x2 {
				for x := x1; x <= x2; x++ {
					gs.Corridors[i].Points = append(gs.Corridors[i].Points,
						entities.Coordinates{X: x, Y: yBetween})
				}
			} else {
				for x := x1; x >= x2; x-- {
					gs.Corridors[i].Points = append(gs.Corridors[i].Points,
						entities.Coordinates{X: x, Y: yBetween})
				}
			}
			for y := yBetween + 1; y < r2.Pos0.Y; y++ {
				gs.Corridors[i].Points = append(gs.Corridors[i].Points,
					entities.Coordinates{X: x2, Y: y})
			}
		// neither horisontal nor vertical neighbours
		default:
			return fmt.Errorf("corridor between rooms %d and %d,"+
				"not horisontally or vertically neighbours", indFrom, indTo)
		}
	}
	return nil
}
