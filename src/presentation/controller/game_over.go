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

const (
	X0   = 2
	Y0   = 2
	YGap = 2
)

func GameOverMode(screen tcell.Screen, g *generation.Generator,
	gs *entities.GameSession, gStat *data.GameStats,
	signal entities.GameSignal) error {
	screen.Clear()
	x, y := X0, Y0
	switch signal {
	case entities.GsFinish, entities.GsLose:
		if signal == entities.GsFinish {
			x, y = uiScreen.DrawWinBanner(screen, x, y)
		} else {
			x, y = uiScreen.DrawLoseBanner(screen, x, y)
		}
	case entities.GsQuit:
		x, y = uiScreen.DrawSavedBanner(screen, x, y)
	}
	if err := ShowGameOverScreen(screen, g, gStat, signal, x, y); err != nil {
		return fmt.Errorf("show screen in GameOverMode: %w", err)
	}
	switch signal {
	case entities.GsFinish, entities.GsLose:
		if err := data.AppendGameStats(gStat,
			config.GameStatsFilename); err != nil {
			return fmt.Errorf("stats save in run after GsFinish or "+
				"GsLose signal: %w", err)
		}
	case entities.GsQuit:
		if err := data.AppendGameSession(gs, gStat,
			config.GameSaveFilename); err != nil {
			return fmt.Errorf("game save in run after GsQuit signal: %w", err)
		}
	}
	return nil
}

func ShowGameOverScreen(screen tcell.Screen, g *generation.Generator,
	gStat *data.GameStats, signal entities.GameSignal, x, y int) error {
	log := g.Log()
	ShowGameResult(screen, gStat, signal, x, y)
	screen.Show()
	if GetPlayerInput(screen) == entities.GsYes {
		log.Debug("ShowStatisticView chosen")
		if signal == entities.GsQuit {
			gStat = nil
		}
		if err := ShowStatisticView(screen, gStat, 2, 2); err != nil {
			return fmt.Errorf("in ShowGameOverScreen: %w", err)
		}
	}
	log.Debug("Quit gameOver screen")
	return nil
}

func ShowGameResult(screen tcell.Screen, gStat *data.GameStats,
	signal entities.GameSignal, x, y int) (int, int) {
	// game result
	style := tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	x, y = uiScreen.DrawLine(screen, "Game result:", x, y, style)
	x, y = uiScreen.DrawStatHeader(screen, x, y, style)
	x, y = uiScreen.DrawStatData(screen, gStat, x, y, 1, style)
	prompt := []string{
		"Press:",
		"	'Y' if you want proceed to statistic view",
		"	any other key for quit",
	}
	x, y = uiScreen.DrawLines(screen, prompt, x, y+YGap, style)

	return x, y
}

func ShowStatisticView(screen tcell.Screen, gStat *data.GameStats,
	x, y int) error {
	screen.Clear()
	AllStats, err := data.ReadGameStats(config.GameStatsFilename)
	if err != nil {
		return fmt.Errorf("reading stats from file in "+
			"ShowGameOverScreen: %w", err)
	}
	sortedStats, summaryStat, place :=
		data.EvaluateAndAddStats(AllStats, gStat)
	scrollAllStats(screen, sortedStats, summaryStat, place)
	screen.Show()
	return nil
}

func scrollAllStats(screen tcell.Screen, stats []*data.GameStats,
	summaryStat *data.GameStats, place int) {
	page := 1
	perPage := 9
	if len(stats) == 0 {
		summaryStat.DeepestLevel = 0
	}
StatLoop:
	for {
		x, y := X0, Y0
		screen.Clear()
		x, y = uiScreen.DrawStaticStatisticView(screen, summaryStat, place,
			len(stats), x, y)
		g0 := (page - 1) * perPage
		g1 := page * perPage
		if g1 > len(stats) {
			g1 = len(stats)
		}
		totalPages := len(stats) / perPage
		if len(stats)%perPage != 0 {
			totalPages++
		}
		uiScreen.DrawAllStatsPages(screen, stats[g0:g1], page, perPage, totalPages, x, y)
		signal := GetPlayerInput(screen)
		switch signal {
		case entities.GsRight:
			if g1 < len(stats) {
				page++
			}
		case entities.GsLeft:
			if g0 > 0 {
				page--
			}
		case entities.GsQuit, entities.GsNo:
			break StatLoop
		}
		screen.Show()
	}
}
