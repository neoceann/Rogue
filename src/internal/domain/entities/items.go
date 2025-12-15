package entities

import (
	"fmt"
	"rogue_game/internal/config"
)

func GetAllItems(gs *GameSession) []*Item {
	var items []*Item

	for itemType, itemSubTypes := range config.ItemsToCreate {
		for _, itemSubType := range itemSubTypes {
			item, err := NewItem(itemType, itemSubType, gs)
			if err == nil {
				items = append(items, item)
			}
		}
	}
	
	return items
}

func NewItem(it config.ItemType, ist config.ItemSubType, gs *GameSession) (*Item, error) {
	item := &Item{
		SubType: ist,
	}
	item.ID = config.ItemConfig[item.SubType].ID
	item.ItemType = config.ItemConfig[item.SubType].ItemType
	item.Name = config.ItemConfig[item.SubType].Name
	item.Stackable = config.ItemConfig[item.SubType].Stackable
	item.MaxStack = config.ItemConfig[item.SubType].MaxStack
	item.Effect = config.ItemEffect{EffectTo: config.ItemConfig[item.SubType].Effect.EffectTo,
									EffectValue: config.ItemConfig[item.SubType].Effect.EffectValue + float64(gs.Level - 1)*config.ItemsStatsIncreasePerLvl,
									Duration: config.ItemConfig[item.SubType].Effect.Duration}

	item.Description = item.CreateItemDescription()

	return item, nil
}

func GetEffectName(effectType config.EffectType) string {
	switch effectType {
	case config.Strength:
		return "Strength"
	case config.Agility:
		return "Agility"
	case config.Health:
		return "Health"
	case config.MaxHealth:
		return "Max health"
	default:
		return "Undef"
	}
}

func (item *Item) CreateItemDescription() string {

	description := ": "

	switch item.ItemType {
	case config.ItemTypeFood:
		description += fmt.Sprintf("Restoring %.1f Health",
			item.Effect.EffectValue)
	case config.ItemTypePotion:
		description += fmt.Sprintf("Gives hero %.1f %s for %d moves",
			item.Effect.EffectValue, GetEffectName(item.Effect.EffectTo), item.Effect.Duration)
	case config.ItemTypeScroll:
		description += fmt.Sprintf("Gives hero %.1f %s permanently",
			item.Effect.EffectValue, GetEffectName(item.Effect.EffectTo))
	case config.ItemTypeWeapon:
		description += fmt.Sprintf("%s increase: %.1f",
		GetEffectName(item.Effect.EffectTo), item.Effect.EffectValue)
	case config.ItemTypeKey:
		description += "Can open door of the same color"
	}

	return description
}