package data

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"rogue_game/internal/config"
	"rogue_game/internal/domain/entities"
	"rogue_game/internal/domain/generation"
	"sort"
	"time"
)

type GameSave struct {
	SaveTime    time.Time             `json:"save_time"`
	GameSession *entities.GameSession `json:"game_session_struct"`
	GameStats   *GameStats            `json:"game_stats"`
}

type GameStats struct {
	GameStartTime     time.Time `json:"game_start_time"`
	PlayerName        string    `json:"player_name"`
	TreasureCollected int       `json:"treasure_collected"`
	DeepestLevel      int       `json:"deepest_level"`
	MonstersDefeated  int       `json:"monsters_defeated"`
	FoodConsumed      int       `json:"food_consumed"`
	PotionsDrunk      int       `json:"potions_drunk"`
	ScrollsRead       int       `json:"scrolls_read"`
	HitsDealt         int       `json:"hits_dealt"`
	HitsReceived      int       `json:"hits_received"`
	TilesTraveled     int       `json:"tiles_traveled"`
	Status            Status    `json:"status"`
}

type Status int

const (
	StatusInvalid Status = iota
	Win
	Lose
	Save
	Continue
)

func NewGameStats(name string) *GameStats {
	return &GameStats{
		GameStartTime: time.Now(),
		PlayerName:    name,
		DeepestLevel:  1,
	}
}

func NewGameSave(gs *entities.GameSession, gStat *GameStats) *GameSave {
	return &GameSave{
		GameSession: gs,
		GameStats:   gStat,
	}
}
func (gStat *GameStats) GetTreasure(inc int) {
	gStat.TreasureCollected += inc
}
func (gStat *GameStats) KillMonster() {
	gStat.MonstersDefeated++
}
func (gStat *GameStats) IncLevel() {
	gStat.DeepestLevel++
}
func (gStat *GameStats) EatFood() {
	gStat.FoodConsumed++
}
func (gStat *GameStats) DrinkPotion() {
	gStat.PotionsDrunk++
}
func (gStat *GameStats) ReadScroll() {
	gStat.ScrollsRead++
}
func (gStat *GameStats) HitDealt() {
	gStat.HitsDealt++
}
func (gStat *GameStats) HitReceived() {
	gStat.HitsReceived++
}
func (gStat *GameStats) Travel() {
	gStat.TilesTraveled++
}

func (gStat *GameStats) SetStatus(status Status) {
	gStat.Status = status
}

// for debug:
func (gStat *GameStats) PrintStats() {
	fmt.Println("GameStats:\n", gStat)
}

func ReadFromJSON[T any](filename string) ([]*T, error) {
	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil, nil // empty list - ok
	}
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", filename, err)
	}
	defer f.Close()

	var contentAll []*T
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue // pass empty lines - check if it's accute
		}

		var contentLine T
		if err := json.Unmarshal(line, &contentLine); err != nil {
			return nil, fmt.Errorf("invalid JSON at line %d: %w", lineNum, err)
		}
		contentAll = append(contentAll, &contentLine)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan %q: %w", filename, err)
	}

	return contentAll, nil
}

func ReadGameStats(filename string) ([]*GameStats, error) {
	AllGamesStats, err := ReadFromJSON[GameStats](filename)
	if err != nil {
		return nil, fmt.Errorf("readGameStats: %w", err)
	}
	return AllGamesStats, nil
}

func ReadGameSave(filename string, g *generation.Generator,
) ([]*GameSave, error) {
	AllGamesSave, err := ReadFromJSON[GameSave](filename)
	if err != nil {
		return nil, fmt.Errorf("readGameStats: %w", err)
	}
	g.Log().Debug("Saved games downloaded")
	return AllGamesSave, nil
}

func (gStat *GameStats) StatToString(n int) []any {
	timeStr := gStat.GameStartTime.Format("15:04 Jan 2")
	var status string
	switch gStat.Status {
	case Win:
		status = "win"
	case Lose:
		status = "lose"
	case Save:
		status = "save"
	case Continue:
		status = "cont"
	}
	return []any{
		fmt.Sprintf("%d", n),
		timeStr,
		gStat.PlayerName,
		status,
		fmt.Sprintf("%d", gStat.TreasureCollected),
		fmt.Sprintf("%d", gStat.DeepestLevel),
		fmt.Sprintf("%d", gStat.MonstersDefeated),
		fmt.Sprintf("%d", gStat.FoodConsumed),
		fmt.Sprintf("%d", gStat.PotionsDrunk),
		fmt.Sprintf("%d", gStat.ScrollsRead),
		fmt.Sprintf("%d", gStat.HitsDealt),
		fmt.Sprintf("%d", gStat.HitsReceived),
		fmt.Sprintf("%d", gStat.TilesTraveled),
	}
}

func AppendGameStats(gStat *GameStats, filename string) error {
	if err := AppendToJSONLine(gStat, filename); err != nil {
		return fmt.Errorf("saveGame: %w", err)
	}
	return nil
}

func AppendGameSession(gs *entities.GameSession, gStat *GameStats,
	filename string) error {
	gStat.DeepestLevel = gs.Level
	gSave := NewGameSave(gs, gStat)
	gSave.SaveTime = time.Now()

	if err := AppendToJSONLine(gSave, filename); err != nil {
		return fmt.Errorf("saveGame: %w", err)
	}
	return nil
}

func AppendToJSONLine(content any, filename string) error {
	if content == nil {
		return errors.New("cannot save nil content")
	}
	data, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	// open file for append, create if there is no file
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open %q for append: %w", filename, err)
	}
	defer f.Close()

	// append to JSON + '\n'
	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write to %q: %w", filename, err)
	}
	return nil
}

// check twice
func SaveJSONArray(saves []*GameSave, filename string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open %q: %w", filename, err)
	}
	defer f.Close()
	for i, save := range saves {
		data, err := json.Marshal(save)
		if err != nil {
			return fmt.Errorf("marshal save #%d: %w", i, err)
		}
		if _, err := f.Write(append(data, '\n')); err != nil {
			return fmt.Errorf("write save #%d: %w", i, err)
		}
	}
	return nil
}

// sorts slice of saved games read from file in descending order by save time
func SortSavedByTimeDesc(saved []*GameSave) {
	sort.Slice(saved, func(i, j int) bool {
		return saved[i].SaveTime.After(saved[j].SaveTime)
	})
}

// sorts slice of saved games read from file in descending order by save time
func SortStatsByTreasureDesc(stats []*GameStats) {
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TreasureCollected > stats[j].TreasureCollected
	})
}

// only if game is finished
func EvaluateAndAddStats(stats []*GameStats,
	gStat *GameStats) ([]*GameStats, *GameStats, int) {

	SortStatsByTreasureDesc(stats)
	// make summary of all games stat
	summaryStat := NewGameStats("TOTAL")
	summaryStat.GameStartTime = time.Now()
	place := -1 //index, starting fom 0

	resultStats := make([]*GameStats, 0)
	if gStat != nil && len(stats) == 0 {
		addToSummaryStat(summaryStat, gStat)
		resultStats = append(resultStats, gStat)
	}
	for count, stat := range stats {
		if gStat != nil && gStat.TreasureCollected >= stat.TreasureCollected {
			place = count
			addToSummaryStat(summaryStat, gStat)
			resultStats = append(resultStats, gStat)
		}
		addToSummaryStat(summaryStat, stat)
		resultStats = append(resultStats, stat)
	}

	return resultStats, summaryStat, place
}

func addToSummaryStat(summaryStat *GameStats, gStat *GameStats) {
	if summaryStat.DeepestLevel < gStat.DeepestLevel {
		summaryStat.DeepestLevel = gStat.DeepestLevel
	}
	summaryStat.FoodConsumed += gStat.FoodConsumed
	if summaryStat.GameStartTime.After(gStat.GameStartTime) {
		summaryStat.GameStartTime = gStat.GameStartTime
	}
	summaryStat.HitsDealt += gStat.HitsDealt
	summaryStat.HitsReceived += gStat.HitsReceived
	summaryStat.MonstersDefeated += gStat.MonstersDefeated
	summaryStat.PotionsDrunk += gStat.PotionsDrunk
	summaryStat.ScrollsRead += gStat.ScrollsRead
	summaryStat.TilesTraveled += gStat.TilesTraveled
	summaryStat.TreasureCollected += gStat.TreasureCollected

}

func ReWriteSavedGamesBack(savedGames []*GameSave, exceptInd int) error {
	right := []*GameSave{}
	if exceptInd < len(savedGames)-1 {
		right = savedGames[exceptInd+1:]
	}
	savedGames = append(savedGames[:exceptInd], right...)
	if err := SaveJSONArray(savedGames, config.GameSaveFilename); err != nil {
		return fmt.Errorf("write array of saved games except selected with"+
			"ind=%d: %w", exceptInd, err)
	}
	return nil
}
