package entities

func NewFogShot(m *MapShot) *FogShot {
	return &FogShot{
		MapShot:     m,
		RoomNow:     &Room{},
		CorridorNow: &Corridor{},
		Walls:       []Coordinates{},
		Visible:     []Coordinates{},
	}
}

func CurrentRoom(gs *GameSession) (*Room, bool) {
	for _, r := range gs.Rooms {
		if InBounds(r.Bounds, gs.Player.Pos) {
			return &r, true
		}
	}
	return nil, false
}

func InBounds(b Bounds, p Coordinates) bool {
	return p.X >= b.Pos0.X && p.X < b.Pos1.X &&
		p.Y >= b.Pos0.Y && p.Y < b.Pos1.Y
}

func (fog *FogShot) AddVisibleRoomWalls(r *Room) {
	width := r.Bounds.Pos1.X - r.Bounds.Pos0.X
	for dx := range width + 2 {
		if fog.MapShot[r.Bounds.Pos0.X+dx-1][r.Bounds.Pos0.Y-1] == Wall {
			fog.Walls = append(fog.Walls,
				Coordinates{r.Bounds.Pos0.X + dx - 1, r.Bounds.Pos0.Y - 1})
		}
		if fog.MapShot[r.Bounds.Pos1.X-dx][r.Bounds.Pos1.Y] == Wall {
			fog.Walls = append(fog.Walls,
				Coordinates{r.Bounds.Pos1.X - dx, r.Bounds.Pos1.Y})
		}
	}
	length := r.Bounds.Pos1.Y - r.Bounds.Pos0.Y
	for dy := range length {
		if fog.MapShot[r.Bounds.Pos0.X-1][r.Bounds.Pos0.Y+dy] == Wall {
			fog.Walls = append(fog.Walls,
				Coordinates{r.Bounds.Pos0.X - 1, r.Bounds.Pos0.Y + dy})
		}
		if fog.MapShot[r.Bounds.Pos1.X][r.Bounds.Pos1.Y-1-dy] == Wall {
			fog.Walls = append(fog.Walls,
				Coordinates{r.Bounds.Pos1.X, r.Bounds.Pos1.Y - 1 - dy})
		}
	}
}

func CurrentCorridor(gs *GameSession) (*Corridor, bool) {
	for _, c := range gs.Corridors {
		for _, p := range c.Points {
			if gs.Player.Pos == p {
				return &c, true
			}
		}
	}
	return nil, false
}

func (fog *FogShot) AddVisibleCorridorWalls(c *Corridor, pos Coordinates) {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if fog.MapShot[pos.X+dx][pos.Y+dy] == Wall {
				fog.Walls = append(fog.Walls,
					Coordinates{pos.X + dx, pos.Y + dy})
			}
		}
	}
}
