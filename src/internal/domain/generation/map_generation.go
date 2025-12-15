package generation

import (
	"fmt"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
)

func (g *Generator) GenerateFullMap(gs *entities.GameSession,
) (*entities.MapShot, error) {
	m := entities.NewMapShot()
	if err := g.GeometryMap(m, gs); err != nil {
		return nil, fmt.Errorf("geometry map: %w", err)
	}
	if err := g.ItemMap(m, gs); err != nil {
		return nil, fmt.Errorf("item map: %w", err)
	}
	if err := g.CharMap(m, gs); err != nil {
		return nil, fmt.Errorf("char map: %w", err)
	}
	return m, nil
}

// Map: 1st layer, puts geometry of the level (rooms, corridors, walls)
// return success or failure and put 'E' on the map on the place of
// intersection
// Add bool return? change for error
func (g *Generator) GeometryMap(m *entities.MapShot,
	gs *entities.GameSession) error {

	// rooms, E if intersection
	for _, r := range gs.Rooms {
		for x := r.Bounds.Pos0.X; x < r.Bounds.Pos1.X; x++ {
			for y := r.Bounds.Pos0.Y; y < r.Bounds.Pos1.Y; y++ {
				if m[x][y] != entities.Outside {
					return fmt.Errorf("rooms overlap")
				}
				m[x][y] = entities.RoomFloor
			}
		}
	}
	// corridors, E if intersection
	for _, corridor := range gs.Corridors {
		for _, c := range corridor.Points {
			if m[c.X][c.Y] == entities.RoomFloor ||
				m[c.X][c.Y] == entities.CorridorFloor {
				return fmt.Errorf("corridors overlap room or corridor")
			}
			m[c.X][c.Y] = entities.CorridorFloor
		}
	}
	// walls, !!! will not return error if no wall between two
	// rooms/corridors
	for x := range m {
		for y := range m[0] {
			if (m[x][y] == entities.Outside) &&
				g.doTouchInterior(m, entities.Coordinates{X: x, Y: y}) {
				m[x][y] = entities.Wall
			}
		}
	}
	// finish (exit to next level)
	if err := putObject(m, gs.Finish, entities.Finish); err != nil {
		return fmt.Errorf("finish overlaps object: %w", err)
	}

	//doors
	for _, door := range gs.Doors {
		m[door.Pos.X][door.Pos.Y] = door
	}

	return nil
}

// Map: assist function for Wall to check if one of neighbor cells is Interior
func (g *Generator) doTouchInterior(m *entities.MapShot,
	c entities.Coordinates) bool {
	for dX := -1; dX <= 1; dX++ {
		for dY := -1; dY <= 1; dY++ {
			if dX == 0 && dY == 0 {
				continue
			}
			nX, nY := c.X+dX, c.Y+dY
			if nX >= 0 && nX < len(m) && nY >= 0 && nY < len(m[0]) {
				if m[nX][nY] == entities.RoomFloor ||
					m[nX][nY] == entities.CorridorFloor {
					return true
				}
			}
		}
	}
	return false
}

// Map: 2nd layer, puts items
// returns success or failure and put 'E' on the map
func (g *Generator) ItemMap(m *entities.MapShot,
	gs *entities.GameSession) error {
	for _, item := range gs.Items {
		if err := putObject(m, item.Pos, item); err != nil {
			return fmt.Errorf("item overlaps object, item name %s: %w", item.Name, err)
		}
	}
	return nil
}

// Map: 3rd layer, puts Monsters and Player
// returns success or failure and put 'E' on the map
func (g *Generator) CharMap(m *entities.MapShot,
	gs *entities.GameSession) error {
	// monsters
	for _, monster := range gs.Monsters {
		if err := putObject(m, monster.Pos, monster); err != nil {
			return fmt.Errorf("monster overlaps other object: %w", err)
		}
	}
	// player
	if err := putObject(m, gs.Player.Pos, gs.Player); err != nil {
		return fmt.Errorf("player overlaps other object: %w", err)
	}
	return nil
}

// Map: puts objest on the map on interior of rooms and corridors
// returns success or failure and put 'E' on the map
func putObject(m *entities.MapShot, c entities.Coordinates,
	obj any) error {
	x := c.X
	y := c.Y
	if m[x][y] == entities.RoomFloor || m[x][y] == entities.CorridorFloor {
		m[x][y] = obj
	} else {
		return fmt.Errorf("object overlaps other x=%d y=%d", x, y)
	}
	return nil
}

func Unfog(fog *entities.FogShot, gs *entities.GameSession,
	g *Generator) *entities.MapShot {
	fogMS := entities.NewMapShot()
	room, roomOk := entities.CurrentRoom(gs)
	if roomOk {
		unfogRoom(fog, room, fogMS)
		fog.AddVisibleRoomWalls(room)
	}
	corridor, corridorOk := entities.CurrentCorridor(gs)
	if corridorOk {
		fog.AddVisibleCorridorWalls(corridor, gs.Player.Pos)
		unfogAngleView(fog, corridor, gs, fogMS)
	}
	unfogNeighbours(fog, gs, fogMS)
	for _, w := range fog.Walls {
		fogMS[w.X][w.Y] = entities.Wall
	}
	return fogMS
}

func unfogRoom(fog *entities.FogShot, r *entities.Room,
	fogMS *entities.MapShot) {
	for x := r.Pos0.X; x < r.Pos1.X; x++ {
		for y := r.Pos0.Y; y < r.Pos1.Y; y++ {
			fogMS[x][y] = fog.MapShot[x][y]
		}
	}
}
func unfogNeighbours(fog *entities.FogShot, gs *entities.GameSession,
	fogMS *entities.MapShot) {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			x := gs.Player.Pos.X + dx
			y := gs.Player.Pos.Y + dy
			fogMS[x][y] = fog.MapShot[x][y]
		}
	}

}
func unfogAngleView(fog *entities.FogShot, c *entities.Corridor,
	gs *entities.GameSession, fogMS *entities.MapShot) {
	// tmpMS is empty
	tmpMS := entities.NewMapShot()
	tmpFog := entities.NewFogShot(fog.MapShot)

	var room entities.Room
	pp := gs.Player.Pos
	switch pp {
	case c.Points[0]:
		room = gs.Rooms[c.FromRoomInd]
	case c.Points[len(c.Points)-1]:
		room = gs.Rooms[c.ToRoomInd]
	default:
		return
	}
	// tmpMS has interiors with obgects in bounds of room
	unfogRoom(fog, &room, tmpMS)
	// tmpFog.Walls have coordinates of all points arround room (no exits)
	tmpFog.AddVisibleRoomWalls(&room)
	// tmpMS adds walls only where walls are in original MapShot
	for _, p := range tmpFog.Walls {
		if fog.MapShot[p.X][p.Y] == entities.Wall {
			tmpMS[p.X][p.Y] = entities.Wall
		}
	}
	// upper and bottom walls
	r0 := room.Bounds.Pos0
	r1 := room.Bounds.Pos1
	if pp.Y == r0.Y-1 || pp.Y == r1.Y {
		for dx := 0; dx < config.GameFieldWidth/config.RoomsInWidth; dx++ {
			for dy := -dx + 1; dy < dx-1; dy++ {
				if pp.Y+dy >= 0 && pp.Y+dy < config.GameFieldLength &&
					pp.X-dx >= 0 && pp.X-dx < config.GameFieldWidth {
					tmpMS[pp.X-dx][pp.Y+dy] = entities.Outside
				}
				if pp.Y+dy >= 0 && pp.Y+dy < config.GameFieldLength &&
					pp.X+dx >= 0 && pp.X+dx < config.GameFieldWidth {
					tmpMS[pp.X+dx][pp.Y+dy] = entities.Outside
				}
			}
		}
	}
	// left and right walls
	if pp.X == r0.X-1 || pp.X == r1.X {
		for dy := 0; dy < config.GameFieldWidth/config.RoomsInWidth; dy++ {
			for dx := -dy + 1; dx < dy-1; dx++ {
				if pp.X+dx >= 0 && pp.X+dx < config.GameFieldWidth &&
					pp.Y-dy >= 0 && pp.Y-dy < config.GameFieldLength {
					tmpMS[pp.X+dx][pp.Y-dy] = entities.Outside
				}
				if pp.X+dx >= 0 && pp.X+dx < config.GameFieldWidth &&
					pp.Y+dy >= 0 && pp.Y+dy < config.GameFieldLength {
					tmpMS[pp.X+dx][pp.Y+dy] = entities.Outside
				}
			}
		}
	}
	// add only visible walls to FogShot
	for x := 0; x < len(tmpMS); x++ {
		for y := 0; y < len(tmpMS[0]); y++ {
			if tmpMS[x][y] != entities.Outside && tmpMS[x][y] != entities.Wall {
				(*fogMS)[x][y] = tmpMS[x][y]
			}
			if tmpMS[x][y] == entities.Wall {
				fog.Walls = append(fog.Walls, entities.Coordinates{X: x, Y: y})
			}
		}
	}
}
