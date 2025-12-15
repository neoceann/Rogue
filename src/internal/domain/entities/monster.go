package entities

import (
	"fmt"
	"rogue_game/internal/config"
)

func NewMonster(monster config.MonsterType, gs *GameSession) (*Monster, error) {
	if monster <= config.MonsterTypeInvalid || monster > config.Mimic {
		return nil, fmt.Errorf("invalid monster type - %d", monster)
	}

	var trick MonsterTrick = false
	switch monster {
	case config.Vampire:
		trick = firstHitToVampire
	case config.Ghost:
		trick = invisibleGhost
	case config.Ogre:
		trick = ogreAttackCooldown
	case config.Snake:
		trick = snakeSleepProc
	case config.Mimic:
		trick = mimicAsItem
	}

	m := &Monster{
		ID:   config.MonsterConfig[monster].ID,
		Type: monster,
		Name: config.MonsterConfig[monster].Name,
		Pos:  Coordinates{},
		Character: config.Character{
			Health:    config.MonsterConfig[monster].Character.Health + float64(gs.Level - 1)*config.StatsIncreasePerLvl,
			MaxHealth: config.MonsterConfig[monster].Character.MaxHealth + float64(gs.Level - 1)*config.StatsIncreasePerLvl,
			Agility:   config.MonsterConfig[monster].Character.Agility + float64(gs.Level - 1)*config.StatsIncreasePerLvl,
			Strength:  config.MonsterConfig[monster].Character.Strength + float64(gs.Level - 1)*config.StatsIncreasePerLvl,
			Treasure:  config.MonsterConfig[monster].Character.Treasure + gs.Level,
		},
		Hostility: config.MonsterConfig[monster].Hostility,
		IsChasing: false,
		Trick:	trick,
		FightStatus: BaseFightStatus{},
		LastDirection: U,
		RoomBounds: Bounds{},
	}
	return m, nil
}

func GetMonsterIndexByPointer(gs *GameSession, monster *Monster) int {
	for i, mnstr := range gs.Monsters {
		if mnstr == monster {
			return i
		}
	}
	return -1
}

func NewMonsterFromTemplate(template *Monster,
) (*Monster, error) {
	if template == nil {
		return nil, fmt.Errorf("template = nil")
	}
	m := *template
	return &m, nil
}