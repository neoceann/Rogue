package controller

import (
	"fmt"
	"rogue_game/internal/config"
	data "rogue_game/datalayer"
	"rogue_game/internal/domain/entities"
	"rogue_game/internal/domain/generation"
	uiScreen "rogue_game/presentation/screen"

	"github.com/gdamore/tcell/v2"
)

func StartMode(screen tcell.Screen, g *generation.Generator) (*entities.GameSession,
	*data.GameStats, error) {
	// screen.Clear() // move to screen show
	var gs *entities.GameSession
	var gStat *data.GameStats
	x, y := X0, Y0
	x, y = uiScreen.DrawStartPicUI(screen, x, y)
	signal := GetPlayerInput(screen)
	log := g.Log()

	// GameSession and GameStats restored from downloads
	if signal == entities.GsYes {
		log.Debug("Check saved games")
		AllGamesSave, err := data.ReadGameSave(config.GameSaveFilename, g)
		data.SortSavedByTimeDesc(AllGamesSave)
		if err != nil {
			return nil, nil, fmt.Errorf("run game read saved games: %w", err)
		}
		downloadInd := choseSavedGame(screen, AllGamesSave)
		if downloadInd == -1 {
			signal = entities.GsNo
		} else {
			gameDownloaded := AllGamesSave[downloadInd]
			if err := data.ReWriteSavedGamesBack(AllGamesSave,
				downloadInd); err != nil {
				return nil, nil, fmt.Errorf("write back dsaved game: %w", err)
			}
			gStat = gameDownloaded.GameStats
			gs = gameDownloaded.GameSession
		}

	}
	// read PlayerName, create new GameSession and GameStats from scratch
	if signal == entities.GsNo || signal == entities.GsQuit {
		name := uiScreen.NameInputUI(screen, x, y+YGap)
		log.Debug("Start new game for '%s'", name)
		gStat = data.NewGameStats(name)
		gs = entities.NewGameSession()
		if err := g.GenerateGameSession(gs); err != nil {
			return nil, nil,
				fmt.Errorf("generate GameSession on level %d: %w",
					gs.Level, err)
		}
	}
	return gs, gStat, nil
}

func choseSavedGame(screen tcell.Screen, games []*data.GameSave) int {
	page := 1
	perPage := 9
	for {
		g0 := (page - 1) * perPage
		g1 := page * perPage
		if g1 > len(games) {
			g1 = len(games)
		}
		totalPages := len(games) / perPage
		if len(games)%perPage != 0 {
			totalPages++
		}
		uiScreen.DrawDownloads(screen, games[g0:g1], page, totalPages)
		signal := GetPlayerInput(screen)
		switch {
		case signal == entities.GsRight:
			if g1 < len(games) {
				page++
			}
		case signal == entities.GsLeft:
			if g0 > 0 {
				page--
			}
		case signal >= entities.GsChoice1 &&
			signal <= entities.GsChoice9:
			ind := (page-1)*perPage + int(signal-entities.GsChoice1)
			if ind < len(games) {
				return ind
			}
		case signal == entities.GsQuit || signal == entities.GsNo:
			return -1
		}
	}
}
