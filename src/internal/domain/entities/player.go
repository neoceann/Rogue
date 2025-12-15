package entities

import "rogue_game/internal/config"

func NewPlayer() *Player {
	return &Player{
		Pos: Coordinates{},
		Character: config.Character{
			Health:    config.PlayerBaseHealth,
			MaxHealth: config.PlayerBaseMaxHealth,
			Agility:   config.PlayerBaseAgility,
			Strength:  config.PlayerBaseStrength,
		},
		Backpack: NewBackpack(),
		FightStatus: PlayerFightStatus{},
		Balance: PlayerGameBalance{HealthItemsUsed: 0},
	}
}

func (hero *Player) UseItem(item *Item, unequipWeapon bool, m *MapShot) {

	if item == nil && !unequipWeapon || unequipWeapon && hero.ActiveWeapon == nil {
		return
	}

	if unequipWeapon {
		hero.UnequipWeapon()
		hero.ActiveWeapon = nil
		return
	}

	if item.Effect.Duration > 0 {
		hero.Buffs = append(hero.Buffs,
			&config.ItemEffect{EffectTo: item.Effect.EffectTo, EffectValue: item.Effect.EffectValue, Duration: item.Effect.Duration})
	}

	if item.ItemType == config.ItemTypeWeapon {
		if hero.ActiveWeapon != nil {
			if hero.Backpack.DroppedItemOnFloor(hero, hero.ActiveWeapon, m) {
				changeWeapon( hero, item)
			}
		} else {
			hero.ActiveWeapon = item
			hero.IncreasePlayerStatsByItem(item)
			hero.Backpack.RemoveItem(item)
		}
	} else {
		hero.IncreasePlayerStatsByItem(item)
		hero.Backpack.ChangesAfterUsingItem(item)
	}
}

func changeWeapon(hero *Player, weapon *Item){
	hero.UnequipWeapon()
	hero.Backpack.RemoveItem(hero.ActiveWeapon)
	hero.Backpack.RemoveItem(weapon)
	hero.ActiveWeapon = weapon
	hero.IncreasePlayerStatsByItem(weapon)
}

func (hero *Player) UnequipWeapon() {
	hero.DecreasePlayerStatsByItem(hero.ActiveWeapon)
	hero.Backpack.AddItem(hero.ActiveWeapon)
}

func (hero *Player) IncreasePlayerStatsByItem(byItem *Item) {

	switch byItem.Effect.EffectTo {
	case config.Health:
		hero.Character.Health += byItem.Effect.EffectValue
		if hero.Health > hero.MaxHealth {
			hero.Health = hero.MaxHealth
		}
	case config.MaxHealth:
		hero.MaxHealth += byItem.Effect.EffectValue
		hero.Health += byItem.Effect.EffectValue
		if hero.Health > hero.MaxHealth {
			hero.Health = hero.MaxHealth
		}
	case config.Agility:
		hero.Agility += byItem.Effect.EffectValue
	case config.Strength:
		hero.Strength += byItem.Effect.EffectValue
	}
}

func (hero *Player) DecreasePlayerStatsByItem(byItem *Item) {

	switch byItem.Effect.EffectTo {
	case config.Health:
		hero.Character.Health -= byItem.Effect.EffectValue
		if hero.Health <= 0 {
			hero.Health = 1.0
		}
	case config.MaxHealth:
		hero.MaxHealth -= byItem.Effect.EffectValue
		hero.Health -= byItem.Effect.EffectValue
		if hero.Health <= 0 {
			hero.Health = 1.0
		}
	case config.Agility:
		hero.Agility -= byItem.Effect.EffectValue
	case config.Strength:
		hero.Strength -= byItem.Effect.EffectValue
	}
}

func (hero *Player) IncreasePlayerStatsByLevel(level int) {
	hero.Health += float64(level)*config.StatsIncreasePerLvl
	hero.MaxHealth += float64(level)*config.StatsIncreasePerLvl
	hero.Agility += float64(level)*config.StatsIncreasePerLvl
	hero.Strength += float64(level)*config.StatsIncreasePerLvl
}

func (hero *Player) DecreasePlayerStatsFromBuffs() {
	for _, b := range hero.Buffs{
		if b.Duration <= 0 {
			switch b.EffectTo {
			case config.Health:
				hero.Health -= b.EffectValue
				if hero.Health <= 0 {
					hero.Health = 1.0
				}
			case config.MaxHealth:
				hero.MaxHealth -= b.EffectValue
				hero.Health -= b.EffectValue
				if hero.Health <= 0 {
					hero.Health = 1.0
				}
			case config.Agility:
				hero.Agility -= b.EffectValue
			case config.Strength:
				hero.Strength -= b.EffectValue
			}
		}
	}
}

func (hero *Player) DecreaseBuffsDuration() {
	for _, b := range hero.Buffs {
		b.Duration--
	}
}

func (hero *Player) CheckExpiredBuffs() (expiredBuffExists bool) {
	for _, b := range hero.Buffs {
		if b.Duration <= 0 {
			expiredBuffExists = true
			break
		}
	}
	return
}

func (hero *Player) RemoveExpiredBuffs() {
	var buffs []*config.ItemEffect
	for _, b := range hero.Buffs {
		if b.Duration > 0 {
			buffs = append(buffs, b)
		}
	}

	hero.Buffs = buffs
}
