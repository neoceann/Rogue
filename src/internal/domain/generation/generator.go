package generation

import (
	"fmt"
	"math/rand"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
	"rogue_game/internal/logger"
	"time"
)

// Generator manages level generation with defind randomness
type Generator struct {
	rng *rand.Rand
	log logger.Logger
	fog bool
}

// NewGenerator creates generator with defined seed:
// seed == 0 use current time (normal game)
// seed != 0 for game reproduction (debugging)
// log for logging
func NewGenerator(seed int64, log logger.Logger) *Generator {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	// debug log
	if log == nil {
		log = logger.NoopLogger{} // безопасный дефолт
	}
	return &Generator{
		rng: rand.New(rand.NewSource(seed)),
		log: log,
		fog: config.FogOfWar,
	}
}

// Log returns logger for other functions can log
func (g *Generator) Log() logger.Logger {
	return g.log
}

func (g *Generator) Fog() bool {
	return g.fog
}

// Generates all fields of GameSession
func (g *Generator) GenerateGameSession(gs *entities.GameSession,
) /**entities.MapShot,*/ error {
	g.log.Debug("Generating LEVEL %d", gs.Level)
	gs.Player.Pos = entities.Coordinates{}

	if err := g.generateRooms(gs); err != nil {
		return fmt.Errorf(
			"generate  Rooms in generate GameSession: %w", err)
	}
	if err := g.generateCorridors(gs); err != nil {
		return fmt.Errorf(
			"generate Corridors in generate GameSession: %w", err)
	}

	if err := g.generatePlayerPos(gs); err != nil {
		return fmt.Errorf(
			"generate PlayerPos in generate GameSession: %w", err)
	}

	g.DoorsKeysGeneration(gs)

	g.generateFinish(gs)

	if err := g.generateItems(gs); err != nil {
		return fmt.Errorf(
			"generate Items in generate GameSession: %w", err)
	}

	if err := g.generateMonsters(gs); err != nil {
		return fmt.Errorf("generate monsters in generate GameSession: %w", err)
	}

	return nil
}
