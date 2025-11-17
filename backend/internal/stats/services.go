package stats

import (
	"database/sql"
	"fmt"
)

type StatsService struct {
	db *sql.DB
}

func NewStatsService(db *sql.DB) *StatsService {
	return &StatsService{db: db}
}

func (s *StatsService) GetCategoryCounts() (*CategoryCountsResponse, error) {
	ingRows, err := s.db.Query(`
		SELECT category, COUNT(*)
		FROM ingredients
		GROUP BY category
		ORDER BY category
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query ingredient categories: %w", err)
	}
	defer ingRows.Close()

	ingCounts := make(map[string]int)
	for ingRows.Next() {
		var category string
		var count int
		if err := ingRows.Scan(&category, &count); err != nil {
			return nil, fmt.Errorf("failed to scan ingredient category: %w", err)
		}
		ingCounts[category] = count
	}
	if err := ingRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ingredient categories: %w", err)
	}

	recRows, err := s.db.Query(`
		SELECT category, COUNT(*)
		FROM recipes
		GROUP BY category
		ORDER BY category
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipe categories: %w", err)
	}
	defer recRows.Close()

	recCounts := make(map[string]int)
	for recRows.Next() {
		var category string
		var count int
		if err := recRows.Scan(&category, &count); err != nil {
			return nil, fmt.Errorf("failed to scan recipe category: %w", err)
		}
		recCounts[category] = count
	}
	if err := recRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating recipe categories: %w", err)
	}

	return &CategoryCountsResponse{
		IngredientCategories: ingCounts,
		RecipeCategories:     recCounts,
	}, nil
}

func (s *StatsService) GetStats() (*Stats, error) {
	stats := &Stats{
		DifficultyDistribution: make(map[string]int),
	}

	err := s.db.QueryRow(`SELECT COUNT(*) FROM ingredients`).Scan(&stats.TotalIngredients)
	if err != nil {
		return nil, fmt.Errorf("failed to get total ingredients: %w", err)
	}

	err = s.db.QueryRow(`SELECT COUNT(*) FROM recipes`).Scan(&stats.TotalRecipes)
	if err != nil {
		return nil, fmt.Errorf("failed to get total recipes: %w", err)
	}

	err = s.db.QueryRow(`SELECT AVG(prep_time_minutes) FROM recipes`).Scan(&stats.AvgPrepTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get average prep time: %w", err)
	}

	err = s.db.QueryRow(`SELECT AVG(cook_time_minutes) FROM recipes`).Scan(&stats.AvgCookTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get average cook time: %w", err)
	}

	rows, err := s.db.Query(`
		SELECT difficulty, COUNT(*)
		FROM recipes
		GROUP BY difficulty
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query difficulty distribution: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var difficulty string
		var count int
		err = rows.Scan(&difficulty, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan difficulty distribution: %w", err)
		}
		stats.DifficultyDistribution[difficulty] = count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating difficulty distribution: %w", err)
	}

	return stats, nil
}
