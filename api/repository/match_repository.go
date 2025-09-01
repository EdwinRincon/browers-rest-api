package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"gorm.io/gorm"
)

type MatchRepository interface {
	CreateMatch(ctx context.Context, match *model.Match) error
	GetMatchByID(ctx context.Context, id uint64) (*model.Match, error)
	GetPaginatedMatches(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Match, int64, error)
	GetMatchesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]model.Match, int64, error)
	GetMatchesByTeamID(ctx context.Context, teamID uint64, sort string, order string, page int, pageSize int) ([]model.Match, int64, error)
	GetNextMatchByTeamID(ctx context.Context, teamID uint64) (*model.Match, error)
	GetDetailedMatchByID(ctx context.Context, id uint64) (*model.Match, error)
	UpdateMatch(ctx context.Context, id uint64, match *model.Match) error
	DeleteMatch(ctx context.Context, id uint64) error
}

type MatchRepositoryImpl struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &MatchRepositoryImpl{db: db}
}

func (mr *MatchRepositoryImpl) CreateMatch(ctx context.Context, match *model.Match) error {
	return mr.db.WithContext(ctx).Create(match).Error
}

// GetMatchByID retrieves a match by its ID with basic preloads
func (mr *MatchRepositoryImpl) GetMatchByID(ctx context.Context, id uint64) (*model.Match, error) {
	var match model.Match
	result := mr.db.WithContext(ctx).
		Preload("Season").
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("MVPPlayer").
		Where("id = ?", id).
		First(&match)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("error getting match by ID: %w", result.Error)
	}
	return &match, nil
}

// GetDetailedMatchByID retrieves a match with all relationships loaded
func (mr *MatchRepositoryImpl) GetDetailedMatchByID(ctx context.Context, id uint64) (*model.Match, error) {
	var match model.Match
	result := mr.db.WithContext(ctx).
		Preload("Season").
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("MVPPlayer").
		Preload("Lineups").
		Preload("Lineups.Player").
		Preload("PlayerStats").
		Preload("PlayerStats.Player").
		Where("id = ?", id).
		First(&match)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("error getting detailed match by ID: %w", result.Error)
	}
	return &match, nil
}

// GetPaginatedMatches retrieves paginated matches with sorting and ordering
func (mr *MatchRepositoryImpl) GetPaginatedMatches(ctx context.Context, sort string, order string, page int, pageSize int) ([]model.Match, int64, error) {
	var matches []model.Match
	var total int64

	// Count total records
	countQuery := mr.db.WithContext(ctx).Model(&model.Match{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting total matches: %w", err)
	}

	// Build the data query with eager loading
	query := mr.db.WithContext(ctx).Model(&model.Match{}).
		Preload("Season").
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("MVPPlayer")

	// Apply sorting if provided
	if sort != "" && (order == "asc" || order == "desc") {
		query = query.Order(fmt.Sprintf("%s %s", sort, order))
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&matches).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching matches: %w", err)
	}

	return matches, total, nil
}

// GetMatchesBySeasonID retrieves matches for a specific season with pagination
func (mr *MatchRepositoryImpl) GetMatchesBySeasonID(ctx context.Context, seasonID uint64, sort string, order string, page int, pageSize int) ([]model.Match, int64, error) {
	var matches []model.Match
	var total int64

	// Count total records for this season
	countQuery := mr.db.WithContext(ctx).Model(&model.Match{}).Where("season_id = ?", seasonID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting matches for season: %w", err)
	}

	// Build the data query
	query := mr.db.WithContext(ctx).Model(&model.Match{}).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Season").
		Where("season_id = ?", seasonID)

	// Apply sorting
	if sort != "" && (order == "asc" || order == "desc") {
		query = query.Order(fmt.Sprintf("%s %s", sort, order))
	} else {
		// Default sort by date
		query = query.Order("date desc")
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&matches).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching matches by season: %w", err)
	}

	return matches, total, nil
}

// GetMatchesByTeamID retrieves matches where a specific team is home or away
func (mr *MatchRepositoryImpl) GetMatchesByTeamID(ctx context.Context, teamID uint64, sort string, order string, page int, pageSize int) ([]model.Match, int64, error) {
	var matches []model.Match
	var total int64

	// Count total records for this team
	countQuery := mr.db.WithContext(ctx).Model(&model.Match{}).
		Where("home_team_id = ? OR away_team_id = ?", teamID, teamID)

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting matches for team: %w", err)
	}

	// Build the data query
	query := mr.db.WithContext(ctx).Model(&model.Match{}).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Season").
		Where("home_team_id = ? OR away_team_id = ?", teamID, teamID)

	// Apply sorting
	if sort != "" && (order == "asc" || order == "desc") {
		query = query.Order(fmt.Sprintf("%s %s", sort, order))
	} else {
		// Default sort by date
		query = query.Order("date desc")
	}

	// Apply pagination
	offset := page * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute the query
	if err := query.Find(&matches).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching matches by team: %w", err)
	}

	return matches, total, nil
}

// GetNextMatchByTeamID retrieves the next scheduled match for a team
func (mr *MatchRepositoryImpl) GetNextMatchByTeamID(ctx context.Context, teamID uint64) (*model.Match, error) {
	var match model.Match
	result := mr.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Season").
		Where("(home_team_id = ? OR away_team_id = ?)", teamID, teamID).
		Where("date >= CURRENT_DATE").
		Where("status = ?", "scheduled").
		Order("date ASC, time ASC").
		First(&match)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("error getting next match for team: %w", result.Error)
	}
	return &match, nil
}

// UpdateMatch updates an existing match
func (mr *MatchRepositoryImpl) UpdateMatch(ctx context.Context, id uint64, match *model.Match) error {
	return mr.db.WithContext(ctx).
		Model(&model.Match{}).
		Where("id = ?", id).
		Updates(match).Error
}

func (mr *MatchRepositoryImpl) DeleteMatch(ctx context.Context, id uint64) error {
	return mr.db.WithContext(ctx).Delete(&model.Match{}, id).Error
}
