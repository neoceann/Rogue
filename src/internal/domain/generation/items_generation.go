package generation

import (
	"fmt"
	"math"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
	"slices"
)

func (g *Generator) generateItems(gs *entities.GameSession) error {
	itemsTemplates := entities.GetAllItems(gs)
	if gs.Level == 1 {
		g.log.Debug("Total %d available Items", len(itemsTemplates))
	}
	if err := g.generateItemsByType(gs, itemsTemplates, config.ItemTypeFood,
		config.StartLevelFoodPerRoom, config.LastLevelFoodPerRoom); err != nil {
		return fmt.Errorf("food items generation: %w", err)
	}
	if err := g.generateItemsByType(gs, itemsTemplates, config.ItemTypePotion,
		config.StartLevelPotionsPerRoom, config.LastLevelPotionsPerRoom,
	); err != nil {
		return fmt.Errorf("potion items generation: %w", err)
	}
	if err := g.generateItemsByType(gs, itemsTemplates, config.ItemTypeScroll,
		config.StartLevelScrollsPerRoom, config.LastLevelScrollsPerRoom,
	); err != nil {
		return fmt.Errorf("srloll items generation: %w", err)
	}
	if err := g.generateWeapons(gs, itemsTemplates); err != nil {
		return fmt.Errorf("generate weapon in generate GameSession: %w", err)
	}
	return nil
}

// returns list of Items by ItemType and take them away from available list
func selectItemsByType(items []*entities.Item,
	t config.ItemType) []*entities.Item {
	itemsByType := []*entities.Item{}
	for i := 0; i < len(items); i++ {
		if items[i].ItemType == t {
			itemsByType = append(itemsByType, items[i])
			if i < len(items)-1 {
				items = append(items[:i], items[i+1:]...)
				i--
			} else {
				items = items[:i]
			}
		}
	}
	return itemsByType
}

// generateItemsByType selects templates by ItemType and spread it evenly and
// randomly in the rooms
func (g *Generator) generateItemsByType(gs *entities.GameSession,
	allItemTemplates []*entities.Item, itemType config.ItemType,
	startLevelPerRoom, lastLevelPerRoom float64) error {
	typeSelectedItyems := selectItemsByType(allItemTemplates, itemType)
	itemsPerLevel := int(config.RoomsNum * levelCoeff(gs.Level,
		startLevelPerRoom, lastLevelPerRoom) / gs.LevelDifficultyCoef)
	itemsInRooms, err := g.allocateObjectsInRooms(itemsPerLevel, []int{})
	if err != nil {
		return fmt.Errorf("generate type %d items: %w", itemType, err)
	}
	countDistributed, err := g.spreadObjectsInRooms(gs,
		typeSelectedItyems, itemsInRooms)
	if countDistributed != itemsPerLevel || err != nil {
		return fmt.Errorf("spread type %d items: should be %d, have: %d",
			itemType, itemsPerLevel, countDistributed)
	}
	g.log.Debug("%d items of type %d generated", countDistributed, itemType)
	return nil
}

// generateWeapons puts weapon on diferent levels
func (g *Generator) generateWeapons(gs *entities.GameSession,
	allItemTemplates []*entities.Item) error {
	weaponItems := selectItemsByType(allItemTemplates, config.ItemTypeWeapon)
	weaponCount := len(weaponItems)
	levelsWithWeapon, err := g.defineLevelsWithWeapon(config.WeaponInFirstLevel,
		config.LevelNum, weaponCount)
	g.log.Debug("Levels with weapons: %v", levelsWithWeapon)
	if err != nil {
		return fmt.Errorf("generateWeapons: %w", err)
	}
	if ind := slices.Index(levelsWithWeapon, gs.Level); ind != -1 {
		// err := putWeaponInRoom(gs, weaponItems[ind], g)
		rooms, err := g.allocateObjectsInRooms(1, []int{1, len(gs.Rooms) - 1})
		if err != nil {
			return fmt.Errorf("putWeaponInRoom: %w", err)
		}
		_, err = g.spreadObjectsInRooms(gs,
			[]*entities.Item{weaponItems[ind]}, rooms)
		if err != nil {
			return fmt.Errorf("generateWeapons: %w", err)
		}
		g.log.Debug("Weapon: %s generated", weaponItems[ind].Name)
	}
	return nil
}

// defineLevelsWithWeapon defines list of Levels where weapon exists.
// startLevelOK - if true starting from 1st level, if false - starting from 2nd
// last level with weapon = LevelsNum - 1 (20 if there are 21 levels)
func (g *Generator) defineLevelsWithWeapon(startLevelOK bool, levelsNum,
	weaponCount int) ([]int, error) {
	if levelsNum < 2 {
		return nil, fmt.Errorf("levelsNum must be >= 2, got %d", levelsNum)
	}
	startLevel := 1
	if !startLevelOK {
		startLevel = 2
	}
	lastLevel := levelsNum - 1
	if startLevel > lastLevel {
		return nil, fmt.Errorf("for weapon: startLevel (%d) > lastLevel (%d)",
			startLevel, lastLevel)
	}
	levelsWithWeapon := make([]int, 0, weaponCount)
	step := float64(lastLevel-startLevel) / float64(weaponCount-1)
	for i := 0; i < weaponCount; i++ {
		level := startLevel + int(math.Round(float64(i)*step))
		if level < startLevel {
			level = startLevel
		}
		if level > lastLevel {
			level = lastLevel
		}
		levelsWithWeapon = append(levelsWithWeapon, level)
	}
	if len(levelsWithWeapon) != weaponCount {
		return []int{}, fmt.Errorf("LevelsWithWeapon %v, number should be %d,"+
			" calculated %d", levelsWithWeapon, weaponCount,
			len(levelsWithWeapon))
	}
	return levelsWithWeapon, nil
}
