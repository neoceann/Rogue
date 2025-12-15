package entities

import "rogue_game/internal/config"

func NewGameSession() *GameSession {
	return &GameSession{
		Level: 1,
		Bounds: Bounds{
			Pos0: Coordinates{X: 0, Y: 0},
			Pos1: Coordinates{X: config.GameFieldWidth,
				Y: config.GameFieldLength},
		},
		Rooms:     []Room{},
		Corridors: []Corridor{},
		Items:     []*Item{},
		Player:    NewPlayer(),
		Monsters:  []*Monster{},
		Doors:     []*Door{},
		Keys:      []*Key{},
		Finish:    Coordinates{},
		LevelDifficultyCoef: 1.0,
	}
}

func (gs *GameSession) NewLevelClean() {
	gs.Rooms = []Room{}
	gs.Corridors = []Corridor{}
	gs.Items = []*Item{}
	gs.Monsters = []*Monster{}
	gs.Doors = []*Door{}
	gs.Keys = []*Key{}
	gs.Finish = Coordinates{}
	gs.Player.Pos = Coordinates{}
}