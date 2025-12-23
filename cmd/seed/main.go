// разделение файла на части с помощью комментариев — инициатива Вахи
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"shumnaya/internal/models"
)

const (
	batchSize    = 10
	playersCount = 15
	seasonsCount = 3
	matchesCount = 5
)

func main() {
	// ENV
	if err := godotenv.Load(); err != nil {
		log.Println(".env файл не найден, используем переменные окружения")
	}

	var dsn string
	if url := os.Getenv("DATABASE_URL"); url != "" {
		dsn = url
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	// CLEAR DB (Postgres-safe)
	fmt.Println("Очистка базы данных...")

	truncateSQL := `
DO $$
DECLARE r RECORD;
BEGIN
	FOR r IN (
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
		  AND tablename NOT IN ('schema_migrations')
	) LOOP
		EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE';
	END LOOP;
END $$;
`
	if err := db.Exec(truncateSQL).Error; err != nil {
		log.Fatal("Ошибка очистки базы:", err)
	}

	fmt.Println("База данных очищена ✓")

	gofakeit.Seed(0)

	players := seedPlayers(db)
	seasons := seedSeasons(db)
	seedMatches(db, players, seasons)
	seedStandings(db, players, seasons)

	fmt.Println("\n=== SEED COMPLETED SUCCESSFULLY ===")
}

// ================= PLAYERS =================

func seedPlayers(db *gorm.DB) []uint {
	players := make([]models.Player, 0, batchSize)
	ids := make([]uint, 0, playersCount)

	fmt.Printf("Seeding players... 0/%d", playersCount)

	for i := 0; i < playersCount; i++ {
		p := models.Player{
			Name:         gofakeit.Name(),
			Email:        fmt.Sprintf("player_%02d@test.com", i),
			PasswordHash: gofakeit.UUID(),
			Rating:       gofakeit.Number(800, 2400),
		}

		players = append(players, p)

		if len(players) >= batchSize {
			db.Session(&gorm.Session{SkipHooks: true}).Create(&players)
			for _, p := range players {
				ids = append(ids, p.ID)
			}
			players = players[:0]
			fmt.Printf("\rSeeding players... %d/%d", i+1, playersCount)
		}
	}

	if len(players) > 0 {
		db.Session(&gorm.Session{SkipHooks: true}).Create(&players)
		for _, p := range players {
			ids = append(ids, p.ID)
		}
	}

	fmt.Println(" ✓")
	return ids
}

// ================= SEASONS =================

func seedSeasons(db *gorm.DB) []uint {
	seasons := make([]models.Season, 0, seasonsCount)
	ids := make([]uint, 0, seasonsCount)

	startYear := time.Now().Year() - seasonsCount + 1
	fmt.Printf("Seeding seasons... 0/%d", seasonsCount)

	for i := 0; i < seasonsCount; i++ {
		seasons = append(seasons, models.Season{
			Name:      fmt.Sprintf("Season %d", startYear+i),
			StartDate: time.Date(startYear+i, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(startYear+i, 12, 31, 23, 59, 59, 0, time.UTC),
			IsActive:  i == seasonsCount-1,
		})
		fmt.Printf("\rSeeding seasons... %d/%d", i+1, seasonsCount)
	}

	db.Session(&gorm.Session{SkipHooks: true}).Create(&seasons)
	for _, s := range seasons {
		ids = append(ids, s.ID)
	}

	fmt.Println(" ✓")
	return ids
}

// ================= MATCHES =================

func seedMatches(db *gorm.DB, players, seasons []uint) {
	matches := make([]models.Match, 0, batchSize)
	fmt.Printf("Seeding matches... 0/%d", matchesCount)

	for i := 0; i < matchesCount; i++ {
		var w, l uint
		for {
			w = players[gofakeit.Number(0, len(players)-1)]
			l = players[gofakeit.Number(0, len(players)-1)]
			if w != l {
				break
			}
		}

		season := seasons[gofakeit.Number(0, len(seasons)-1)]

		matches = append(matches, models.Match{
			WinnerID: w,
			LoserID:  l,
			SeasonID: season,
			Score:    "3:1",
			PlayedAt: gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
		})

		if len(matches) >= batchSize {
			db.Session(&gorm.Session{SkipHooks: true}).Create(&matches)
			matches = matches[:0]
			fmt.Printf("\rSeeding matches... %d/%d", i+1, matchesCount)
		}
	}

	if len(matches) > 0 {
		db.Session(&gorm.Session{SkipHooks: true}).Create(&matches)
	}

	fmt.Println(" ✓")
}

// ================= STANDINGS =================

func seedStandings(db *gorm.DB, players, seasons []uint) {
	fmt.Println("Seeding standings...")

	standings := make([]models.Standing, 0, batchSize)

	for _, season := range seasons {
		for _, player := range players {
			standings = append(standings, models.Standing{
				PlayerID: player,
				SeasonID: season,
				Wins:     gofakeit.Number(0, 5),
				Losses:   gofakeit.Number(0, 5),
				Points:   gofakeit.Number(0, 15),
				Rank:     gofakeit.Number(1, playersCount),
			})

			if len(standings) >= batchSize {
				db.Session(&gorm.Session{SkipHooks: true}).Create(&standings)
				standings = standings[:0]
			}
		}
	}

	if len(standings) > 0 {
		db.Session(&gorm.Session{SkipHooks: true}).Create(&standings)
	}

	fmt.Println("Standings ✓")
}
