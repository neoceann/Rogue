package controller

import (
	"fmt"
	"rogue_game/internal/config"
	data "rogue_game/datalayer"
	"rogue_game/internal/domain/character"
	"rogue_game/internal/domain/entities"
	"rogue_game/internal/domain/generation"
	uiScreen "rogue_game/presentation/screen"

	"github.com/gdamore/tcell/v2"
)

func RunMode(screen tcell.Screen, g *generation.Generator,
	gs *entities.GameSession, gStat *data.GameStats,
) (entities.GameSignal, error) {
	log := g.Log()
	log.Info("-> RunMode: level=%d", gs.Level)
	defer log.Info("<- RunMode")

	m, err := g.GenerateFullMap(gs)
	if err != nil {
		return 0, fmt.Errorf(
			"generate FullMap in generate GameSession: %w", err)
	}

	fog := entities.NewFogShot(m)

	log.Debug("GameSession for level %d generated:\n\tBounds: %v,\n"+
		"\tPlayerPos: %v", gs.Level, gs.Bounds, gs.Player.Pos)

	ClearAndPrintScreen(screen, m, fog, gs, g, gStat)

	signal, err := runRogue(gs, g, fog, m, screen, gStat)

	return signal, err
}

func runRogue(gs *entities.GameSession, g *generation.Generator, fog *entities.FogShot, m *entities.MapShot, screen tcell.Screen, gStat *data.GameStats)(entities.GameSignal, error){
	for {
		userInput := GetPlayerInput(screen)
		signal, err := LevelCycle(gs, m, userInput, g, gStat, screen)
		if err != nil {
			return entities.GsError,
				fmt.Errorf("signal from controler with error: %w", err)
		}

		switch signal {
		case entities.GsQuit:
			gStat.SetStatus(data.Save)
		case entities.GsLose:
			gStat.SetStatus(data.Lose)
		case entities.GsNextLevel:
			if gs.Level == config.LevelNum {
				signal = entities.GsFinish
				gStat.SetStatus(data.Win)
			} else {
				gs.Level++
				gs.NewLevelClean()
				if err := g.GenerateGameSession(gs); err != nil {
					return entities.GsError,
						fmt.Errorf("generate GameSession on level %d: %w", gs.Level, err)
				}
				gStat.IncLevel()
				signal = entities.GsNextLevel
			}
		default:
			ClearAndPrintScreen(screen, m, fog, gs, g, gStat)
			continue
		}
		return signal, nil
	}
}

func LevelCycle(gs *entities.GameSession, m *entities.MapShot,
	userInput entities.GameSignal, g *generation.Generator,
	gStat *data.GameStats, screen tcell.Screen) (entities.GameSignal, error) {

	canMove := false
	hitMonster := false
	var moveDir entities.Direction

	signal := userInputHandler(userInput, &canMove, &hitMonster, &moveDir, m, gs, screen, gStat)
	if signal == entities.GsQuit || signal == entities.GsNextLevel {
		return signal, nil
	}

	if canMove || hitMonster {
		if canMove {
			finish := movePlayer(moveDir, m, gs, gStat)
			if finish {
				return entities.GsNextLevel, nil
			}
		}

		if hitMonster {
			playerDefeated := character.FightProcess(moveDir, gs, m, gStat)
			if playerDefeated {
				signal = entities.GsLose
			}
		} else {
			gs.Player.FightStatus.InFight = false
		}

		decreaseBuffsDuration(gs.Player)
		moveMonsters(gs, m)
	}

	return signal, nil
}

func userInputHandler(userInput entities.GameSignal, canMove, hitMonster *bool, moveDir *entities.Direction,
	m *entities.MapShot, gs *entities.GameSession, screen tcell.Screen, gStat *data.GameStats) entities.GameSignal{
	switch userInput {

	case entities.GsQuit, entities.GsNextLevel:
		return userInput

	case entities.GsUp, entities.GsDown, entities.GsRight, entities.GsLeft:
		var moveD entities.Direction
		var newX, newY int
		switch userInput {
		case entities.GsUp:
			moveD = entities.U
			newX = gs.Player.Pos.X
			newY = gs.Player.Pos.Y-character.PlayerStep
		case entities.GsDown:
			moveD = entities.D
			newX = gs.Player.Pos.X
			newY = gs.Player.Pos.Y+character.PlayerStep
		case entities.GsRight:
			moveD = entities.R
			newX = gs.Player.Pos.X+character.PlayerStep
			newY = gs.Player.Pos.Y
		case entities.GsLeft:
			moveD = entities.L
			newX = gs.Player.Pos.X-character.PlayerStep
			newY = gs.Player.Pos.Y
		}

		*canMove, *hitMonster = entities.CanMove(newX, newY, m, false, gs.Player.Backpack)
		*moveDir = moveD

	case entities.GsWeapon, entities.GsFood, entities.GsPotion, entities.GsScroll:
		item, unequipFlag := RequestItemForUse(gs, screen, userInput)
		if item != nil || unequipFlag {
			if item != nil {
				switch item.ItemType {
				case config.ItemTypeFood:
					gStat.EatFood()
					gs.Player.Balance.HealthItemsUsed++

				case config.ItemTypePotion:
					gStat.DrinkPotion()

				case config.ItemTypeScroll:
					gStat.ReadScroll()
				}
			}
			gs.Player.UseItem(item, unequipFlag, m)
		}	
	}

	return entities.GsNoAction
}

func movePlayer(moveDir entities.Direction, m *entities.MapShot,
	gs *entities.GameSession, gStat *data.GameStats) (playerOnFinish bool) {
		
	movementChanges := &entities.MovementChanges{
		OldPos:     entities.Coordinates{},
		NewPos:     entities.Coordinates{},
		OldElement: nil,
		NewElement: nil,
	}
		entities.MoveUnit(moveDir, 0, 0, m, gs, gs.Player, character.PlayerStep, movementChanges)
		switch v := m[gs.Player.Pos.X][gs.Player.Pos.Y].(type) {
		case *entities.Item:
			itemAdded := gs.Player.Backpack.AddItem(v)
			if itemAdded {
				RemoveItemFromGS(gs, v)
			} else {
				movementChanges.OldElement = v
			}

		case entities.Field:
			switch v {
			case entities.Finish:
				for _,key := range config.ItemsToCreate[config.ItemTypeKey]{
					gs.Backpack.RemoveItem(gs.Backpack[key])
				}
				balanceAdjustment(gs)
				gs.Player.IncreasePlayerStatsByLevel(gs.Level)
				playerOnFinish = true
			}
		}
		UpdateMapShot(m, movementChanges)
		gStat.Travel()
		
		return
}

func balanceAdjustment(gs *entities.GameSession) {
	if gs.Level % config.BalanceCheckPerLevel == 0 {
		increaseDifficulty := gs.Player.Balance.HealthItemsUsed < 
							config.BalanceAvgHealthItemsUsedPerLevel*config.BalanceCheckPerLevel - config.BalanceAvgHealthItemsUsedTolerance
		decreaseDifficulty := gs.Player.Balance.HealthItemsUsed > 
					config.BalanceAvgHealthItemsUsedPerLevel*config.BalanceCheckPerLevel + config.BalanceAvgHealthItemsUsedTolerance

		if increaseDifficulty {
			gs.LevelDifficultyCoef += config.BalanceLevelDifficultyCoefStep
		} else if decreaseDifficulty {
			gs.LevelDifficultyCoef -= config.BalanceLevelDifficultyCoefStep
			if gs.LevelDifficultyCoef <= 0 {
				gs.LevelDifficultyCoef = 0.1
			}
		}

		gs.Player.Balance.HealthItemsUsed = 0
	}
}

func moveMonsters(gs *entities.GameSession, m *entities.MapShot){

	movementChanges := &entities.MovementChanges{
		OldPos:     entities.Coordinates{},
		NewPos:     entities.Coordinates{},
		OldElement: nil,
		NewElement: nil,
	}

	character.TriggerGhostTrick(gs)

	for _, monst := range gs.Monsters {
		monsterstep := character.MonsterBaseStep
		if monst.Type == config.Ogre {
			monsterstep = character.OgreStep
		}

		monst.FightStatus.InFight = monst.Pos.DistTo(gs.Player.Pos) <= 1
		if !monst.FightStatus.InFight {
			monst.IsChasing = monst.Pos.DistTo(gs.Player.Pos) <= int(monst.Hostility)
			if monst.IsChasing {
				if monst.Type == config.Ghost {
					monst.Trick = false
				}
				newpos := character.CalcMonsterPosTowardPlayer(monst, gs.Player, m)
				if newpos != &monst.Pos {
					entities.MoveUnit(0, newpos.X, newpos.Y, m, gs, monst, monsterstep, movementChanges)
					UpdateMapShot(m, movementChanges)
				}
			}

			if !monst.IsChasing && monst.Type != config.Mimic {
				canm := false
				newDir, newX, newY := character.MonsterDirectionHandler(monst, monsterstep)
				if monst.Type == config.Ghost {
					canm, _ = entities.CanMove(newX, newY, m, true, gs.Player.Backpack)
				} else {
					canm, _ = entities.CanMove(monst.Pos.X+newX,
						monst.Pos.Y+newY, m, true, gs.Player.Backpack)
				}

				if canm && !monst.FightStatus.InFight {
					entities.MoveUnit(newDir, newX, newY, m, gs, monst, monsterstep, movementChanges)
					UpdateMapShot(m, movementChanges)
				}
			}
		}
	}
}

func decreaseBuffsDuration(player *entities.Player){
	player.DecreaseBuffsDuration()
	if player.CheckExpiredBuffs() {
		player.DecreasePlayerStatsFromBuffs()
		player.RemoveExpiredBuffs()
	}
}

func ClearAndPrintScreen(screen tcell.Screen, m *entities.MapShot,
	fog *entities.FogShot, gs *entities.GameSession, g *generation.Generator,
	gStat *data.GameStats) {
	screen.Clear()
	if g.Fog() {
		uiScreen.PrintMapUI(screen, *generation.Unfog(fog, gs, g))
	} else {
		uiScreen.PrintMapUI(screen, *m)
	}
	uiScreen.PrintInfoUI(screen, *gs, gStat)
	uiScreen.PrintMonstersAroundHero(screen, gs, m)
	uiScreen.PrintFightInfo(screen, gs, m)
	screen.Show()
}