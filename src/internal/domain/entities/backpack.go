package entities

import (
	"rogue_game/internal/config"
)

func NewBackpack() Backpack {
	return make(Backpack)
}

func (backpack Backpack) IsFull() bool {
	return len(backpack) >= config.BackpackCapacity
}

func (backpack Backpack) AddItem(item *Item) bool {

	if item == nil {
		return false
	}

	itemAdded := false
	lockedByStackSize := false

	if backpack.IsFull() && !item.Stackable {
		return itemAdded
	}

	if item.Stackable {
		if currentItem, exists := backpack[item.SubType]; exists {
			lockedByStackSize = currentItem.StackCounter >= currentItem.MaxStack
			if !lockedByStackSize {
				currentItem.StackCounter++
				itemAdded = true
			}
		}
	}

	if !itemAdded && !backpack.IsFull() && !lockedByStackSize {
		item.StackCounter = 1
		backpack[item.SubType] = item
		itemAdded = true
	}

	return itemAdded
}

func (backpack Backpack) RemoveItem(item *Item) {
	if item != nil {
		delete(backpack, item.SubType)
	}
}

func (backpack Backpack) DecreaseItemStackSize(item *Item){
	item.StackCounter--
}

func (backpack Backpack) ChangesAfterUsingItem(item *Item){
	
	backpack.DecreaseItemStackSize(item)

	if item.StackCounter <= 0 {
		backpack.RemoveItem(item)
	}
}

func (backpack Backpack) GetItemsByType(it config.ItemType) []*Item {
	var items []*Item

	for _, item := range backpack {
		if item.ItemType == it {
			items = append(items, item)
		}
	}

	if len(items) != 0 {
		return items
	}

	return nil
}

func (backpack Backpack) DroppedItemOnFloor(hero *Player, item *Item, m *MapShot) bool {
	ex, ey, correct := hero.Pos.FindFirstEmptyAround(m, hero.Backpack)

	if correct {
		m[ex][ey] = item
	}

	return correct
}

func (backpack Backpack) HaveKeyForDoor(d *Door) bool {
	switch d.Color {
	case RedDoor:
		if key, exists := backpack[config.ItemSubTypeKeyRed]; exists {
			backpack.RemoveItem(key)
			return true
		}
	case GreenDoor:
		if key, exists := backpack[config.ItemSubTypeKeyGreen]; exists {
			backpack.RemoveItem(key)
			return true
		}
	case BlueDoor:
	if key, exists := backpack[config.ItemSubTypeKeyBlue]; exists {
		backpack.RemoveItem(key)
		return true
	}
	}
	return false
}