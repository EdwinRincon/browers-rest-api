package persistence

import (
	"fmt"
	"strings"

	"gorm.io/gorm/clause"
)

// EntityType represents the type of entity being sorted
type EntityType string

const (
	EntityArticle     EntityType = "article"
	EntityLineup      EntityType = "lineup"
	EntityMatch       EntityType = "match"
	EntityPlayer      EntityType = "player"
	EntityPlayerStats EntityType = "player_stats"
	EntityPlayerTeam  EntityType = "player_team"
	EntityRole        EntityType = "role"
	EntitySeason      EntityType = "season"
	EntityTeam        EntityType = "team"
	EntityTeamStats   EntityType = "team_stats"
	EntityUser        EntityType = "user"
)

// SortOrder indicates the direction of sorting (ascending or descending).
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// IsValid reports whether the sort order is valid.
func (s SortOrder) IsValid() bool {
	return s == SortOrderAsc || s == SortOrderDesc
}

// SortConfig represents the configuration for a sortable field.
type SortConfig struct {
	SQLFragment string
	IsRelation  bool
}

// EntitySortConfig holds all sorting configurations for all entities
var EntitySortConfig = map[EntityType]map[string]SortConfig{
	EntityArticle: {
		"id":         {SQLFragment: "articles.id", IsRelation: false},
		"title":      {SQLFragment: "articles.title", IsRelation: false},
		"content":    {SQLFragment: "articles.content", IsRelation: false},
		"img_banner": {SQLFragment: "articles.img_banner", IsRelation: false},
		"date":       {SQLFragment: "articles.date", IsRelation: false},
		"season_id":  {SQLFragment: "articles.season_id", IsRelation: false},
		"created_at": {SQLFragment: "articles.created_at", IsRelation: false},
		"updated_at": {SQLFragment: "articles.updated_at", IsRelation: false},
		"season":     {SQLFragment: "(SELECT year FROM seasons WHERE seasons.id = articles.season_id)", IsRelation: true},
	},
	EntityLineup: {
		"id":          {SQLFragment: "lineups.id", IsRelation: false},
		"position":    {SQLFragment: "lineups.position", IsRelation: false},
		"player_id":   {SQLFragment: "lineups.player_id", IsRelation: false},
		"match_id":    {SQLFragment: "lineups.match_id", IsRelation: false},
		"starting":    {SQLFragment: "lineups.starting", IsRelation: false},
		"created_at":  {SQLFragment: "lineups.created_at", IsRelation: false},
		"updated_at":  {SQLFragment: "lineups.updated_at", IsRelation: false},
		"player_name": {SQLFragment: "(SELECT nick_name FROM players WHERE players.id = lineups.player_id)", IsRelation: true},
		"match_date":  {SQLFragment: "(SELECT kickoff FROM matches WHERE matches.id = lineups.match_id)", IsRelation: true},
	},
	EntityMatch: {
		"id":            {SQLFragment: "matches.id", IsRelation: false},
		"status":        {SQLFragment: "matches.status", IsRelation: false},
		"kickoff":       {SQLFragment: "matches.kickoff", IsRelation: false},
		"location":      {SQLFragment: "matches.location", IsRelation: false},
		"home_goals":    {SQLFragment: "matches.home_goals", IsRelation: false},
		"away_goals":    {SQLFragment: "matches.away_goals", IsRelation: false},
		"home_team_id":  {SQLFragment: "matches.home_team_id", IsRelation: false},
		"away_team_id":  {SQLFragment: "matches.away_team_id", IsRelation: false},
		"season_id":     {SQLFragment: "matches.season_id", IsRelation: false},
		"mvp_player_id": {SQLFragment: "matches.mvp_player_id", IsRelation: false},
		"created_at":    {SQLFragment: "matches.created_at", IsRelation: false},
		"updated_at":    {SQLFragment: "matches.updated_at", IsRelation: false},
		"season":        {SQLFragment: "(SELECT year FROM seasons WHERE seasons.id = matches.season_id)", IsRelation: true},
		"home_team":     {SQLFragment: "(SELECT short_name FROM teams WHERE teams.id = matches.home_team_id)", IsRelation: true},
		"away_team":     {SQLFragment: "(SELECT short_name FROM teams WHERE teams.id = matches.away_team_id)", IsRelation: true},
	},
	EntityPlayer: {
		"id":                {SQLFragment: "players.id", IsRelation: false},
		"nick_name":         {SQLFragment: "players.nick_name", IsRelation: false},
		"height":            {SQLFragment: "players.height", IsRelation: false},
		"country":           {SQLFragment: "players.country", IsRelation: false},
		"secondary_country": {SQLFragment: "players.secondary_country", IsRelation: false},
		"foot":              {SQLFragment: "players.foot", IsRelation: false},
		"age":               {SQLFragment: "players.age", IsRelation: false},
		"squad_number":      {SQLFragment: "players.squad_number", IsRelation: false},
		"rating":            {SQLFragment: "players.rating", IsRelation: false},
		"matches":           {SQLFragment: "players.matches", IsRelation: false},
		"y_cards":           {SQLFragment: "players.y_cards", IsRelation: false},
		"r_cards":           {SQLFragment: "players.r_cards", IsRelation: false},
		"goals":             {SQLFragment: "players.goals", IsRelation: false},
		"assists":           {SQLFragment: "players.assists", IsRelation: false},
		"saves":             {SQLFragment: "players.saves", IsRelation: false},
		"position":          {SQLFragment: "players.position", IsRelation: false},
		"injured":           {SQLFragment: "players.injured", IsRelation: false},
		"career_summary":    {SQLFragment: "players.career_summary", IsRelation: false},
		"mvp_count":         {SQLFragment: "players.mvp_count", IsRelation: false},
		"user_id":           {SQLFragment: "players.user_id", IsRelation: false},
		"created_at":        {SQLFragment: "players.created_at", IsRelation: false},
		"updated_at":        {SQLFragment: "players.updated_at", IsRelation: false},
	},
	EntityPlayerStats: {
		"id":          {SQLFragment: "player_stats.id", IsRelation: false},
		"player_id":   {SQLFragment: "player_stats.player_id", IsRelation: false},
		"season_id":   {SQLFragment: "player_stats.season_id", IsRelation: false},
		"matches":     {SQLFragment: "player_stats.matches", IsRelation: false},
		"goals":       {SQLFragment: "player_stats.goals", IsRelation: false},
		"assists":     {SQLFragment: "player_stats.assists", IsRelation: false},
		"saves":       {SQLFragment: "player_stats.saves", IsRelation: false},
		"y_cards":     {SQLFragment: "player_stats.y_cards", IsRelation: false},
		"r_cards":     {SQLFragment: "player_stats.r_cards", IsRelation: false},
		"mvp_count":   {SQLFragment: "player_stats.mvp_count", IsRelation: false},
		"created_at":  {SQLFragment: "player_stats.created_at", IsRelation: false},
		"updated_at":  {SQLFragment: "player_stats.updated_at", IsRelation: false},
		"player_name": {SQLFragment: "(SELECT nick_name FROM players WHERE players.id = player_stats.player_id)", IsRelation: true},
		"season":      {SQLFragment: "(SELECT year FROM seasons WHERE seasons.id = player_stats.season_id)", IsRelation: true},
	},
	EntityPlayerTeam: {
		"id":          {SQLFragment: "player_teams.id", IsRelation: false},
		"player_id":   {SQLFragment: "player_teams.player_id", IsRelation: false},
		"team_id":     {SQLFragment: "player_teams.team_id", IsRelation: false},
		"start_date":  {SQLFragment: "player_teams.start_date", IsRelation: false},
		"end_date":    {SQLFragment: "player_teams.end_date", IsRelation: false},
		"created_at":  {SQLFragment: "player_teams.created_at", IsRelation: false},
		"updated_at":  {SQLFragment: "player_teams.updated_at", IsRelation: false},
		"player_name": {SQLFragment: "(SELECT nick_name FROM players WHERE players.id = player_teams.player_id)", IsRelation: true},
		"team_name":   {SQLFragment: "(SELECT short_name FROM teams WHERE teams.id = player_teams.team_id)", IsRelation: true},
	},
	EntityRole: {
		"id":          {SQLFragment: "roles.id", IsRelation: false},
		"name":        {SQLFragment: "roles.name", IsRelation: false},
		"description": {SQLFragment: "roles.description", IsRelation: false},
		"created_at":  {SQLFragment: "roles.created_at", IsRelation: false},
		"updated_at":  {SQLFragment: "roles.updated_at", IsRelation: false},
	},
	EntitySeason: {
		"id":         {SQLFragment: "seasons.id", IsRelation: false},
		"year":       {SQLFragment: "seasons.year", IsRelation: false},
		"start_date": {SQLFragment: "seasons.start_date", IsRelation: false},
		"end_date":   {SQLFragment: "seasons.end_date", IsRelation: false},
		"is_current": {SQLFragment: "seasons.is_current", IsRelation: false},
		"created_at": {SQLFragment: "seasons.created_at", IsRelation: false},
		"updated_at": {SQLFragment: "seasons.updated_at", IsRelation: false},
	},
	EntityTeam: {
		"id":              {SQLFragment: "teams.id", IsRelation: false},
		"full_name":       {SQLFragment: "teams.full_name", IsRelation: false},
		"short_name":      {SQLFragment: "teams.short_name", IsRelation: false},
		"primary_color":   {SQLFragment: "teams.primary_color", IsRelation: false},
		"secondary_color": {SQLFragment: "teams.secondary_color", IsRelation: false},
		"shield":          {SQLFragment: "teams.shield", IsRelation: false},
		"next_match_id":   {SQLFragment: "teams.next_match_id", IsRelation: false},
		"created_at":      {SQLFragment: "teams.created_at", IsRelation: false},
		"updated_at":      {SQLFragment: "teams.updated_at", IsRelation: false},
	},
	EntityTeamStats: {
		"id":              {SQLFragment: "team_stats.id", IsRelation: false},
		"team_id":         {SQLFragment: "team_stats.team_id", IsRelation: false},
		"season_id":       {SQLFragment: "team_stats.season_id", IsRelation: false},
		"matches_played":  {SQLFragment: "team_stats.matches_played", IsRelation: false},
		"wins":            {SQLFragment: "team_stats.wins", IsRelation: false},
		"draws":           {SQLFragment: "team_stats.draws", IsRelation: false},
		"losses":          {SQLFragment: "team_stats.losses", IsRelation: false},
		"goals_for":       {SQLFragment: "team_stats.goals_for", IsRelation: false},
		"goals_against":   {SQLFragment: "team_stats.goals_against", IsRelation: false},
		"goal_difference": {SQLFragment: "team_stats.goal_difference", IsRelation: false},
		"points":          {SQLFragment: "team_stats.points", IsRelation: false},
		"rank":            {SQLFragment: "team_stats.rank", IsRelation: false},
		"created_at":      {SQLFragment: "team_stats.created_at", IsRelation: false},
		"updated_at":      {SQLFragment: "team_stats.updated_at", IsRelation: false},
		"team_name":       {SQLFragment: "(SELECT short_name FROM teams WHERE teams.id = team_stats.team_id)", IsRelation: true},
		"season":          {SQLFragment: "(SELECT year FROM seasons WHERE seasons.id = team_stats.season_id)", IsRelation: true},
	},
	EntityUser: {
		"id":          {SQLFragment: "users.id", IsRelation: false},
		"name":        {SQLFragment: "users.name", IsRelation: false},
		"last_name":   {SQLFragment: "users.last_name", IsRelation: false},
		"username":    {SQLFragment: "users.username", IsRelation: false},
		"birthdate":   {SQLFragment: "users.birthdate", IsRelation: false},
		"img_profile": {SQLFragment: "users.img_profile", IsRelation: false},
		"img_banner":  {SQLFragment: "users.img_banner", IsRelation: false},
		"role_id":     {SQLFragment: "users.role_id", IsRelation: false},
		"created_at":  {SQLFragment: "users.created_at", IsRelation: false},
		"updated_at":  {SQLFragment: "users.updated_at", IsRelation: false},
		"role_name":   {SQLFragment: "(SELECT name FROM roles WHERE roles.id = users.role_id)", IsRelation: true},
	},
}

// BuildOrderClause validates inputs and returns a safe GORM clause for ORDER BY.
// It supports both direct columns (via clause API) and pre-approved raw SQL fragments.
func BuildOrderClause(entityType EntityType, sort, order string) (clause.OrderByColumn, string, error) {
	// Default to created_at if no sort field is provided
	if sort == "" {
		sort = "created_at"
	}

	entityConfig, exists := EntitySortConfig[entityType]
	if !exists {
		return clause.OrderByColumn{}, "", fmt.Errorf("unsupported entity type: %s", entityType)
	}

	config, exists := entityConfig[sort]
	if !exists {
		return clause.OrderByColumn{}, "", fmt.Errorf("invalid sort field: %s", sort)
	}

	sortOrder := SortOrder(strings.ToLower(order))
	if !sortOrder.IsValid() {
		sortOrder = SortOrderAsc
	}

	if config.IsRelation {
		// Whitelisted raw SQL fragment for relationship sorts
		sql := fmt.Sprintf("%s %s", config.SQLFragment, strings.ToUpper(string(sortOrder)))
		return clause.OrderByColumn{}, sql, nil
	}

	// Safe structured clause for direct fields
	return clause.OrderByColumn{
		Column: clause.Column{Name: config.SQLFragment},
		Desc:   sortOrder == SortOrderDesc,
	}, "", nil
}

// GetDefaultSortField returns the default sort field for an entity
func GetDefaultSortField(entityType EntityType) string {
	return "created_at"
}

// ValidateSort validates if a sort field exists for the given entity type
func ValidateSort(entityType EntityType, sort string) error {
	if sort == "" {
		return nil
	}

	entityConfig, exists := EntitySortConfig[entityType]
	if !exists {
		return fmt.Errorf("unsupported entity type: %s", entityType)
	}

	_, exists = entityConfig[sort]
	if !exists {
		return fmt.Errorf("invalid sort field: %s", sort)
	}

	return nil
}
