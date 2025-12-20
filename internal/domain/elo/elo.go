package elo

import "math"

func ExpectedScore(playerElo, opponentElo int) float64 {
	return 1.0 / (1.0 + math.Pow(10, float64(opponentElo-playerElo)/400))
}

func CalculateK(games, wins int) int {
	if games < 10 {
		return 40
	}

	winrate := float64(wins) / float64(games)

	switch {
	case winrate > 0.7:
		return 20
	case winrate < 0.4:
		return 30
	default:
		return 25
	}
}

func NewRating(playerElo, opponentElo int, isWin bool, games, wins int) int {
	expected := ExpectedScore(playerElo, opponentElo)

	score := 0.0
	if isWin {
		score = 1
	}

	k := CalculateK(games, wins)

	return int(math.Round(
		float64(playerElo) + float64(k)*(score-expected),
	))
}
