package app

import (
	"fmt"
	"os"
	"rogue_game/internal/domain/entities"
	"rogue_game/internal/domain/generation"
	"rogue_game/internal/logger"
	"rogue_game/presentation/controller"
	"time"
)

func Run() error {
	startGame := time.Now()
	verbose := true
	log := logger.NewBufferLogger(startGame, verbose)
	log.Info("-> Run: Game started")

	g := generation.NewGenerator(0, log)

	screen, err := controller.CreateScreen()
	if err != nil {
		return fmt.Errorf("screen failed in Run: %w", err)
	}
	defer screen.Fini()

	// returns newly generated or restored from download data
	gs, gStat, err := controller.StartMode(screen, g)
	if err != nil {
		return fmt.Errorf("in Run StartMode: %w", err)
	}

	var signal entities.GameSignal
	for {
		log.Debug("---===::: LEVEL %d STARTED :::===---", gs.Level)

		if signal, err = controller.RunMode(screen, g, gs, gStat); err != nil {
			return fmt.Errorf("RunMode in Run: %w", err)
		}
		if signal != entities.GsNextLevel {
			break
		}
	}

	controller.GameOverMode(screen, g, gs, gStat, signal)

	log.Debug("Signal to quit level = %d", signal)
	log.Info("********* LEVEL FINISHED *********")

	if verbose {
		if err := os.WriteFile("debug.log", log.Bytes(), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
		}
	}
	log.Info("<- Run: Game finished")

	return nil
}
