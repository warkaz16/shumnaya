package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"shumnaya/internal/models"
)

const (
	batchSize          = 1000
	playersCount       = 200_000
	seasonsCount       = 40
	standingsPerSeason = 25_000 // total ~ 1,000,000
	matchesCount       = 2_000_000
	softDeletePercent  = 5
)

type seasonMeta struct {
	ID        uint
	StartDate time.Time
	EndDate   time.Time
}

func main() {
	// 1) Load env & connect to DB
	_ = godotenv.Load()

	db := mustConnect()

	// 2) Clear existing data
	fmt.Println("Clearing existing data...")
	if err := truncateAll(db); err != nil {
		log.Fatalf("failed to truncate tables: %v", err)
	}
	fmt.Println("Data cleared ✓")

	// 3) Seed in dependency order
	gofakeit.Seed(0)

	playerIDs := seedPlayers(db)
	seasons := seedSeasons(db)
	seedStandings(db, seasons, playerIDs)
	seedMatches(db, seasons, playerIDs)

	// 4) Summary
	fmt.Println("\n=== Seeding completed ===")
	fmt.Printf("Players:   %d\n", len(playerIDs))
	fmt.Printf("Seasons:   %d\n", len(seasons))
	fmt.Printf("Standings: ~%d\n", seasonsCount*standingsPerSeason)
	fmt.Printf("Matches:   %d\n", matchesCount)
}

func mustConnect() *gorm.DB {
	// Prefer DATABASE_URL if present
	if url := os.Getenv("DATABASE_URL"); url != "" {
		db, err := gorm.Open(postgres.Open(url), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}
		return db
	}

	host := getenvDefault("DB_HOST", "localhost")
	port := getenvDefault("DB_PORT", "5432")
	user := getenvDefault("DB_USER", "postgres")
	pass := getenvDefault("DB_PASSWORD", "postgres")
	name := getenvDefault("DB_NAME", "shumnaya")
	sslm := getenvDefault("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, name, sslm)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}

func getenvDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func truncateAll(db *gorm.DB) error {
	// Determine which of the expected tables actually exist in the current schema,
	// then TRUNCATE only those tables to avoid "relation does not exist" errors.
	rows, err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = current_schema()").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	existing := make(map[string]bool)
	var t string
	for rows.Next() {
		if err := rows.Scan(&t); err == nil {
			existing[t] = true
		}
	}

	expected := []string{"matches", "standings", "seasons", "players"}
	found := make([]string, 0, len(expected))
	for _, name := range expected {
		if existing[name] {
			found = append(found, name)
			continue
		}
		// also try singular form (e.g. "match")
		if strings.HasSuffix(name, "s") {
			singular := strings.TrimSuffix(name, "s")
			if existing[singular] {
				found = append(found, singular)
				continue
			}
		}
	}

	if len(found) == 0 {
		// If no tables found, try to auto-migrate models to ensure tables exist
		log.Println("No existing target tables found; running AutoMigrate to create tables...")
		if err := db.AutoMigrate(&models.Player{}, &models.Season{}, &models.Match{}, &models.Standing{}); err != nil {
			return fmt.Errorf("auto migrate failed: %w", err)
		}

		// assume default pluralized names created by GORM
		found = []string{"matches", "standings", "seasons", "players"}
	}

	sql := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(found, ", "))
	return db.Exec(sql).Error
}

func seedPlayers(db *gorm.DB) []uint {
	total := playersCount
	ids := make([]uint, 0, total)
	buf := make([]models.Player, 0, batchSize)

	fmt.Printf("Seeding players... 0/%d", total)

	// Generate in batches
	for i := 0; i < total; i++ {
		created := gofakeit.DateRange(time.Now().AddDate(-5, 0, 0), time.Now().AddDate(-1, 0, 0))
		updated := gofakeit.DateRange(created, time.Now())

		p := models.Player{
			// gorm.Model timestamps
			// Set via embedded struct literal
			// CreatedAt / UpdatedAt honored by GORM when non-zero
			// DeletedAt set for a small percentage to simulate soft delete
			// Other fields
			Name:         gofakeit.Name(),
			Email:        fmt.Sprintf("user_%06d@test.com", i),
			PasswordHash: fmt.Sprintf("hash_%s", gofakeit.UUID()),
			Rating:       gofakeit.Number(1000, 2200),
		}
		// Set embedded gorm.Model fields
		p.Model.CreatedAt = created
		p.Model.UpdatedAt = updated
		if gofakeit.Number(1, 100) <= softDeletePercent {
			del := gofakeit.DateRange(created, updated)
			p.Model.DeletedAt = gorm.DeletedAt{Time: del, Valid: true}
		}

		buf = append(buf, p)
		if len(buf) >= batchSize {
			if err := db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(&buf, batchSize).Error; err != nil {
				log.Fatalf("insert players batch failed: %v", err)
			}
			for _, r := range buf {
				ids = append(ids, r.ID)
			}
			buf = buf[:0]
			fmt.Printf("\rSeeding players... %d/%d", i+1, total)
		}
	}

	if len(buf) > 0 {
		if err := db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(&buf, batchSize).Error; err != nil {
			log.Fatalf("insert players batch failed: %v", err)
		}
		for _, r := range buf {
			ids = append(ids, r.ID)
		}
	}

	fmt.Println(" ✓")
	return ids
}

func seedSeasons(db *gorm.DB) []seasonMeta {
	total := seasonsCount
	result := make([]seasonMeta, 0, total)
	buf := make([]models.Season, 0, batchSize)

	fmt.Printf("Seeding seasons... 0/%d", total)

	baseStart := time.Now().AddDate(-10, 0, 0)

	for i := 0; i < total; i++ {
		// Create seasons of ~3 months, spread over last ~10 years
		start := gofakeit.DateRange(baseStart, time.Now().AddDate(0, -1, 0))
		// Ensure start before end by at least a week
		end := start.AddDate(0, 3, gofakeit.Number(0, 14))
		if end.After(time.Now().AddDate(1, 0, 0)) {
			end = time.Now().AddDate(1, 0, 0)
		}
		name := fmt.Sprintf("Season %02d - %d", (i%12)+1, start.Year())

		s := models.Season{
			Name:      name,
			StartDate: start,
			EndDate:   end,
			IsActive:  false,
		}
		// timestamps
		created := gofakeit.DateRange(start.AddDate(0, -1, 0), start)
		updated := gofakeit.DateRange(end, end.AddDate(0, 1, 0))
		s.Model.CreatedAt = created
		s.Model.UpdatedAt = updated
		if gofakeit.Number(1, 100) <= softDeletePercent {
			del := gofakeit.DateRange(created, updated)
			s.Model.DeletedAt = gorm.DeletedAt{Time: del, Valid: true}
		}

		buf = append(buf, s)
		if len(buf) >= batchSize {
			if err := db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(&buf, batchSize).Error; err != nil {
				log.Fatalf("insert seasons batch failed: %v", err)
			}
			for _, r := range buf {
				result = append(result, seasonMeta{ID: r.ID, StartDate: r.StartDate, EndDate: r.EndDate})
			}
			buf = buf[:0]
			fmt.Printf("\rSeeding seasons... %d/%d", len(result), total)
		}
	}

	if len(buf) > 0 {
		if err := db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(&buf, batchSize).Error; err != nil {
			log.Fatalf("insert seasons batch failed: %v", err)
		}
		for _, r := range buf {
			result = append(result, seasonMeta{ID: r.ID, StartDate: r.StartDate, EndDate: r.EndDate})
		}
	}

	// Mark a random season as active (or a few)
	if len(result) > 0 {
		activeN := int(math.Max(1, float64(len(result))/10)) // ~10%
		for i := 0; i < activeN; i++ {
			idx := gofakeit.Number(0, len(result)-1)
			_ = db.Model(&models.Season{}).Where("id = ?", result[idx].ID).Update("is_active", true).Error
		}
	}

	fmt.Println(" ✓")
	return result
}

func seedStandings(db *gorm.DB, seasons []seasonMeta, playerIDs []uint) {
	// We'll create standingsPerSeason per season
	if len(seasons) == 0 || len(playerIDs) == 0 {
		return
	}

	total := len(seasons) * standingsPerSeason
	fmt.Printf("Seeding standings... 0/%d", total)

	buf := make([]models.Standing, 0, batchSize)
	done := 0

	for si, sm := range seasons {
		// sample playersPerSeason unique players by random stepping
		step := int(math.Max(1, float64(len(playerIDs))/float64(standingsPerSeason)))
		// start offset
		offset := gofakeit.Number(0, step)

		count := 0
		for pi := offset; pi < len(playerIDs) && count < standingsPerSeason; pi += step {
			wins := gofakeit.Number(0, 60)
			losses := gofakeit.Number(0, 60)
			points := wins*3 + gofakeit.Number(0, wins) // bonus variation

			st := models.Standing{
				PlayerID: playerIDs[pi],
				SeasonID: sm.ID,
				Wins:     wins,
				Losses:   losses,
				Points:   points,
				Rank:     count + 1,
			}
			// timestamps near season end
			created := gofakeit.DateRange(sm.StartDate, sm.EndDate)
			updated := gofakeit.DateRange(created, sm.EndDate.AddDate(0, 1, 0))
			st.Model.CreatedAt = created
			st.Model.UpdatedAt = updated
			if gofakeit.Number(1, 100) <= softDeletePercent {
				del := gofakeit.DateRange(created, updated)
				st.Model.DeletedAt = gorm.DeletedAt{Time: del, Valid: true}
			}

			buf = append(buf, st)
			count++
			done++

			if len(buf) >= batchSize {
				if err := db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(&buf, batchSize).Error; err != nil {
					log.Fatalf("insert standings batch failed: %v", err)
				}
				buf = buf[:0]
				fmt.Printf("\rSeeding standings... %d/%d", done, total)
			}
		}
		_ = si // just to keep si in scope for potential future use
	}

	if len(buf) > 0 {
		if err := db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(&buf, batchSize).Error; err != nil {
			log.Fatalf("insert standings batch failed: %v", err)
		}
	}

	fmt.Println(" ✓")
}

func seedMatches(db *gorm.DB, seasons []seasonMeta, playerIDs []uint) {
	total := matchesCount
	fmt.Printf("Seeding matches... 0/%d", total)

	buf := make([]models.Match, 0, batchSize)

	for i := 0; i < total; i++ {
		// pick distinct players
		wi := gofakeit.Number(0, len(playerIDs)-1)
		li := gofakeit.Number(0, len(playerIDs)-1)
		for li == wi {
			li = gofakeit.Number(0, len(playerIDs)-1)
		}
		winnerID := playerIDs[wi]
		loserID := playerIDs[li]

		// pick season
		sidx := gofakeit.Number(0, len(seasons)-1)
		sm := seasons[sidx]

		// PlayedAt between season start and end
		played := gofakeit.DateRange(sm.StartDate, sm.EndDate)
		created := gofakeit.DateRange(sm.StartDate, played)
		updated := gofakeit.DateRange(played, sm.EndDate.AddDate(0, 1, 0))

		// score like 3:0 .. 5:4 ensuring left>right
		left := gofakeit.Number(1, 5)
		right := gofakeit.Number(0, left-1)
		score := fmt.Sprintf("%d:%d", left, right)

		delta := gofakeit.Number(1, 25)

		m := models.Match{
			WinnerID:           winnerID,
			LoserID:            loserID,
			SeasonID:           sm.ID,
			Score:              score,
			WinnerRatingChange: delta,
			LoserRatingChange:  -delta,
			PlayedAt:           played,
		}
		m.Model.CreatedAt = created
		m.Model.UpdatedAt = updated
		if gofakeit.Number(1, 100) <= softDeletePercent {
			del := gofakeit.DateRange(created, updated)
			m.Model.DeletedAt = gorm.DeletedAt{Time: del, Valid: true}
		}

		buf = append(buf, m)

		if len(buf) >= batchSize {
			if err := db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(&buf, batchSize).Error; err != nil {
				log.Fatalf("insert matches batch failed: %v", err)
			}
			buf = buf[:0]
			fmt.Printf("\rSeeding matches... %d/%d", i+1, total)
		}
	}

	if len(buf) > 0 {
		if err := db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(&buf, batchSize).Error; err != nil {
			log.Fatalf("insert matches batch failed: %v", err)
		}
	}

	fmt.Println(" ✓")
}
