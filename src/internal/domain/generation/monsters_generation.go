package generation

import (
	"fmt"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
)

func (g *Generator) GetMonstersTemplates(gs *entities.GameSession,
) ([]*entities.Monster, error) {
	monsters := []*entities.Monster{}
	for t := range config.MonsterConfig {
		mNew, err := entities.NewMonster(t, gs)
		if err != nil {
			return nil, fmt.Errorf("GetMonstersTemplates: %w", err)
		}
		monsters = append(monsters, mNew)
	}
	return monsters, nil

}

func (g *Generator) generateMonsters(gs *entities.GameSession) error {
	monstersTemplates, err := g.GetMonstersTemplates(gs)
	if err != nil {
		return fmt.Errorf("generateMonsters: %w", err)
	}
	monstersCount := int(config.RoomsNum * levelCoeff(gs.Level,
		config.StartLevelMonstersPerRoom, float64(config.LastLevelMonstersPerRoom)) * gs.LevelDifficultyCoef)
	monstersInRooms, err := g.allocateObjectsInRooms(monstersCount,
		[]int{gs.Player.RoomInd})
	g.log.Info("Monsters: %v", monstersInRooms)
	if err != nil {
		return fmt.Errorf("generateMonsters: %w", err)
	}
	countDistributed, err := g.spreadMonstersInRooms(gs, monstersTemplates,
		monstersInRooms)
	if countDistributed != monstersCount || err != nil {
		return fmt.Errorf("monsters generated=%d, should be %d: %w",
			countDistributed, monstersCount, err)
	}
	g.log.Debug("Generated %d monsters from %d templates", countDistributed,
		len(monstersTemplates))
	return nil
}

func (g *Generator) spreadMonstersInRooms(gs *entities.GameSession,
	mTemplates []*entities.Monster, rooms map[int]int) (int, error) {
	mCount := 0 // checker, can be deleted
	for i, n := range rooms {
		for range n {
			if len(mTemplates) < 1 {
				return mCount, fmt.Errorf("len(all monsters) < 1, cant use rnd")
			}
			mTmp := mTemplates[g.rng.Intn(len(mTemplates))]
			mNew, err := entities.NewMonsterFromTemplate(mTmp)
			if err != nil {
				return mCount, fmt.Errorf("could not create copy of "+
					"monster template: %w", err)
			}
			pos, err := g.attemptToPlaceObject(gs, &gs.Rooms[i], MaxAttempts)
			if err != nil {
				return mCount, fmt.Errorf("spreadMonstersInRooms: %w", err)
			}
			mNew.Pos = pos
			mNew.RoomBounds = gs.Rooms[i].Bounds
			gs.Monsters = append(gs.Monsters, mNew)
			mCount++
		}
	}
	return mCount, nil
}
