package character

import (
	"math/rand"
	"rogue_game/internal/config"
	data "rogue_game/datalayer"
	"rogue_game/internal/domain/entities"
)

func FightProcess(dir entities.Direction, gs *entities.GameSession, m *entities.MapShot, gstat *data.GameStats) (playerDefeated bool) {
	monstersAround := MonstersAroundHero(gs, m)
	monsterToHit := FindMonsterByDirection(gs, dir, monstersAround)

	if monsterToHit == nil {
		return false
	}

	gs.Player.FightStatus.MissingHit = true
	gs.Player.FightStatus.InFight = true
	gs.Player.FightStatus.TargetMonsterName = monsterToHit.Name
	monsterToHit.FightStatus.MissingHit = true

	playerHitMonster(gs, monsterToHit, gstat)
	if checkDefeat(&monsterToHit.Character) {
		actionsIfMonsterDefeated(monsterToHit, gs, m, gstat)
	}

	gs.Player.FightStatus.Asleep = false
	monstersHitPlayer(gs, monstersAround, gstat)
	playerDefeated = checkDefeat(&gs.Player.Character)

	return
}

func actionsIfMonsterDefeated(monster *entities.Monster, gs *entities.GameSession, m *entities.MapShot, gstat *data.GameStats) {
	monster.FightStatus.InFight = false
	gs.Player.FightStatus.InFight = false
	gstat.KillMonster()
	gstat.GetTreasure(monster.Treasure)
	gs.Player.Treasure += monster.Treasure
	m[monster.Pos.X][monster.Pos.Y] = entities.RoomFloor
	removeMonsterFromGS(gs, monster)
}

func playerHitMonster(gs *entities.GameSession, monsterToHit *entities.Monster, gstat *data.GameStats) {
	if !gs.Player.FightStatus.Asleep {
		if !isHitMissed(&gs.Player.Character, &monsterToHit.Character, monsterToHit) {
			gstat.HitDealt()
			gs.Player.FightStatus.MissingHit = false
			gs.Player.FightStatus.DamageDone = DoDamage(&gs.Player.Character, &monsterToHit.Character, nil)
		}	
	}
}

func monstersHitPlayer(gs *entities.GameSession, monstersAround []*entities.Monster, gstat *data.GameStats) {
	for _, mnstr := range monstersAround {
		if (!bool(mnstr.Trick) && mnstr.Type == config.Ogre || mnstr.Type != config.Ogre) && mnstr.FightStatus.InFight {
			if !isHitMissed(&mnstr.Character, &gs.Player.Character, nil) {
				gstat.HitReceived()
				mnstr.FightStatus.MissingHit = false
				mnstr.FightStatus.DamageDone = DoDamage(&mnstr.Character, &gs.Player.Character, mnstr)
			}
		}
		MonsterTrick(gs.Player, mnstr)
	}
}

func MonstersAroundHero(gs *entities.GameSession, m *entities.MapShot) []*entities.Monster {
	var monsters []*entities.Monster

	for _, mnst := range gs.Monsters {
		mnst.FightStatus.InFight = false
	}

	for x := -1; x < 2; x++ {
        for y := -1; y < 2; y++ {
			cell := m[gs.Player.Pos.X + x][gs.Player.Pos.Y + y]
			if monster, ok := cell.(*entities.Monster); ok {
				monsters = append(monsters, monster)
				monster.FightStatus.InFight = true
			}
		}
	}

	return monsters
}

func FindMonsterByDirection(gs *entities.GameSession, dir entities.Direction, monsters []*entities.Monster) *entities.Monster {

	if monsters == nil {
		return nil
	}

	var checkAttackCoordinates entities.Coordinates
	switch dir{
	case entities.U:
		checkAttackCoordinates = entities.Coordinates{X: gs.Player.Pos.X, Y: gs.Player.Pos.Y - PlayerStep}
	case entities.D:
		checkAttackCoordinates = entities.Coordinates{X: gs.Player.Pos.X, Y: gs.Player.Pos.Y + PlayerStep}
	case entities.R:
		checkAttackCoordinates = entities.Coordinates{X: gs.Player.Pos.X + PlayerStep, Y: gs.Player.Pos.Y}
	case entities.L:
		checkAttackCoordinates = entities.Coordinates{X: gs.Player.Pos.X - PlayerStep, Y: gs.Player.Pos.Y}
	}

	for _, monster := range monsters {
		if checkAttackCoordinates == monster.Pos {
			return monster
		}
	}

	return nil
}

func isHitMissed(attackerStats, defenderStats *config.Character, monster *entities.Monster) bool {
    hitChance := float64(config.BaseHitChance + (attackerStats.Agility-defenderStats.Agility) * config.AgilityImpact)
    
	if hitChance < config.MinHitChance {
		hitChance = config.MinHitChance
	}
	if hitChance > config.MaxHitChance {
		hitChance = config.MaxHitChance
	}

	rnd := rand.Float64()

	if monster == nil {
		return rnd > hitChance
	}
	return (rnd > hitChance || bool(monster.Trick) && monster.Type == config.Vampire)
}

func DoDamage(attackerStats, defenderStats *config.Character, monster *entities.Monster) float64 {
	damage := attackerStats.Strength

	if monster != nil && monster.Type == config.Vampire {
		damage = defenderStats.MaxHealth*(config.VamprireHitMaxHealthPercent/100.0)
	}

	defenderStats.Health -= damage

	return damage
}

func checkDefeat(defenderStats *config.Character) bool {
	return defenderStats.Health <= 0.0
}

func removeMonsterFromGS(gs *entities.GameSession, monster *entities.Monster) {
	for i, it := range gs.Monsters {
		if monster == it {
			gs.Monsters = append(gs.Monsters[:i], gs.Monsters[i+1:]...)
			break
		}
	}
}

func MonsterTrick(hero *entities.Player, monster *entities.Monster){
	switch monster.Type {
	case config.Vampire, config.Mimic:
		monster.Trick = false
	case config.Ghost:
		if !monster.FightStatus.InFight {
			monster.Trick = (1 - rand.Float64()) < 0.5
		}
	case config.Ogre:
		monster.Trick = !monster.Trick
	case config.Snake:
		monster.Trick = rand.Float64() < config.SnakeSleepChance
		hero.FightStatus.Asleep = bool(monster.Trick)
	}
}

func TriggerGhostTrick(gs *entities.GameSession) {
	for _, monster := range gs.Monsters {
		if monster.Type == config.Ghost && !monster.FightStatus.InFight && !monster.IsChasing {
			MonsterTrick(gs.Player, monster)
		}
	}
}