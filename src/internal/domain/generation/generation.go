package generation

import (
	"fmt"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
	"rogue_game/internal/ds"
)

const (
	ConsumablesCoeff = 1
)

// generateRooms populates gs.Rooms with randomly sized rooms on the field
// a. defines random size of room, one per grid cell
// b. config.CenteredGeometry defines geometry (centered or diffused)
// c. locates room on the field in the grid cells
func (g *Generator) generateRooms(gs *entities.GameSession) error {
	rooms := make([]entities.Room, config.RoomsNum)
	for i := range rooms {
		roomWidth, roomLength, err := g.generateRoomSizeRand()
		if err != nil {
			return fmt.Errorf("failed to generate size of room %d: %w", i, err)
		}
		var pos0, pos1 entities.Coordinates
		if config.CenteredGeometry {
			pos0, pos1, err = g.centeredRooms(i, roomWidth, roomLength)
		} else {
			pos0, pos1, err = g.dancingRooms(i, roomWidth, roomLength)
		}
		if err != nil {
			return fmt.Errorf("failed to locate room %d: %w", i, err)
		}
		rooms[i] = entities.Room{Bounds: entities.Bounds{
			Pos0: pos0,
			Pos1: pos1,
		},
			RoomInd: i,
		}
	}
	gs.Rooms = rooms
	return nil
}

// generate and place Player into random room on random position inside room
func (g *Generator) generatePlayerPos(gs *entities.GameSession) error {
	// randomise choice of room where Player will start
	startRoomInd := g.rng.Intn(config.RoomsNum)
	g.log.Debug("Player room ind=%d", startRoomInd)
	// randomize position inside room
	pos, err := g.randomCoordinatesInsideRoom1(&gs.Rooms[startRoomInd])
	if err != nil {
		return fmt.Errorf("player generation on room %d:\n"+
			"x=%d y=%d: %w", startRoomInd, pos.X, pos.Y, err)
	}
	gs.Player.Pos = pos
	gs.Player.RoomInd = startRoomInd
	return nil
}

// generateFinish finds room with longest path from start point
// puts coordinates of Finish in the center of room
func (g *Generator) generateFinish(gs *entities.GameSession) {
	allTrays := gs.AllTrays()
	// choose farthest room
	RoomInd := allTrays[len(allTrays)-1][len(allTrays[len(allTrays)-1])-1]
	gs.Finish = findCenterOfBounds(gs.Rooms[RoomInd].Bounds)
}

// generateCorridors builds all corridors of current level
// a. connects all rooms with algorythm
// b. config.CenteredGeometry defines geometry (straight or curved)
// c. set up coordinates of corridors which connect rooms
// returns error if cannot build corridor between two rooms
func (g *Generator) generateCorridors(gs *entities.GameSession) error {
	g.connectRoomsMST(gs)
	var err error
	if config.CenteredGeometry {
		err = g.centeredCorridors(gs)
	} else {
		err = g.dancingCorridors(gs)
	}
	if err != nil {
		return fmt.Errorf("generate corridor: %w", err)
	}
	return nil
}

// connectRoomsMST initialise all corridors
// with FromRoomInd ToRoomInd (indexes of rooms)
func (g *Generator) connectRoomsMST(gs *entities.GameSession) {
	mst := ds.BuildMST(config.RoomsInWidth, config.RoomsInLength)
	for _, e := range mst {
		c := entities.Corridor{FromRoomInd: e.N1, ToRoomInd: e.N2}
		gs.Corridors = append(gs.Corridors, c)
	}
}

// orderRoomsFromToInc orders rooms in their index increasing order
// for further work left -> right or up -> down
func (g *Generator) orderRoomsFromToInc(gs *entities.GameSession,
	ind int) (int, int, entities.Room, entities.Room) {
	indFrom := gs.Corridors[ind].FromRoomInd
	indTo := gs.Corridors[ind].ToRoomInd
	r1 := gs.Rooms[indFrom]
	r2 := gs.Rooms[indTo]
	if indFrom > indTo {
		bTmp := r1
		r1 = r2
		r2 = bTmp
		eTmp := indFrom
		indFrom = indTo
		indTo = eTmp
		length := len(gs.Corridors[ind].Points)
		for i := range length / 2 {
			tmpCoord := gs.Corridors[ind].Points[i]
			gs.Corridors[ind].Points[i] =
				gs.Corridors[ind].Points[length-i-1]
			gs.Corridors[ind].Points[length-i-1] = tmpCoord
		}
	}

	return indFrom, indTo, r1, r2
}

// levelCoeff calculates avg per room (float64) number of objects
// (1st level: start number, 21st level: end number)
func levelCoeff(level int, start, end float64) float64 {
	inc := (end - start) / (config.LevelNum - 1)
	return start + float64(level-1)*inc
}
