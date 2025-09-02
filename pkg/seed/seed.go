package seed

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/EdwinRincon/browersfc-api/api/model"
	"github.com/EdwinRincon/browersfc-api/pkg/orm"
	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
)

// SeedDatabase populates the database with fake data
func SeedDatabase() error {
	// Get database connection
	db, err := orm.GetDBInstance()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Initialize gofakeit with a seed for reproducibility
	gofakeit.Seed(42)

	// Create seed data in the correct order to respect dependencies
	roles, err := seedRoles(db)
	if err != nil {
		return fmt.Errorf("failed to seed roles: %w", err)
	}

	users, err := seedUsers(db, roles)
	if err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	teams, err := seedTeams(db)
	if err != nil {
		return fmt.Errorf("failed to seed teams: %w", err)
	}

	seasons, err := seedSeasons(db)
	if err != nil {
		return fmt.Errorf("failed to seed seasons: %w", err)
	}

	players, err := seedPlayers(db, users)
	if err != nil {
		return fmt.Errorf("failed to seed players: %w", err)
	}

	err = seedPlayerTeams(db, players, teams, seasons)
	if err != nil {
		return fmt.Errorf("failed to seed player teams: %w", err)
	}

	err = seedTeamStats(db, teams, seasons)
	if err != nil {
		return fmt.Errorf("failed to seed team stats: %w", err)
	}

	matches, err := seedMatches(db, teams, seasons, players)
	if err != nil {
		return fmt.Errorf("failed to seed matches: %w", err)
	}

	err = seedLineups(db, players, matches)
	if err != nil {
		return fmt.Errorf("failed to seed lineups: %w", err)
	}

	err = seedArticles(db, seasons)
	if err != nil {
		return fmt.Errorf("failed to seed articles: %w", err)
	}

	err = seedPlayerStats(db, players, matches, seasons)
	if err != nil {
		return fmt.Errorf("failed to seed player stats: %w", err)
	}

	err = updateTeamsNextMatch(db, teams, matches)
	if err != nil {
		return fmt.Errorf("failed to update teams with next match: %w", err)
	}

	return nil
}

// seedRoles creates 5 roles in the database
func seedRoles(db *gorm.DB) ([]model.Role, error) {
	roles := []model.Role{
		{Name: "admin", Description: "Administrator with full access"},
		{Name: "coach", Description: "Team coach with management permissions"},
		{Name: "player", Description: "Registered player"},
		{Name: "referee", Description: "Match official"},
		{Name: "fan", Description: "Regular user with limited access"},
	}

	if err := db.Create(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

// seedUsers creates 5 users in the database
func seedUsers(db *gorm.DB, roles []model.Role) ([]model.User, error) {
	users := make([]model.User, 5)

	for i := range users {
		birthdate := gofakeit.Date()

		users[i] = model.User{
			ID:         gofakeit.UUID(),
			Name:       gofakeit.FirstName(),
			LastName:   gofakeit.LastName(),
			Username:   gofakeit.Email(),
			Birthdate:  &birthdate,
			ImgProfile: gofakeit.URL(),
			ImgBanner:  gofakeit.URL(),
			RoleID:     roles[i].ID,
		}
	}

	if err := db.Create(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// seedTeams creates 5 teams in the database
func seedTeams(db *gorm.DB) ([]model.Team, error) {
	teams := make([]model.Team, 5)

	teamColors := []string{"#FF0000", "#0000FF", "#00FF00", "#FFFF00", "#FF00FF"}
	teamSecondColors := []string{"#FFFFFF", "#000000", "#CCCCCC", "#333333", "#666666"}

	for i := range teams {
		teams[i] = model.Team{
			FullName:  gofakeit.Company(),
			ShortName: gofakeit.LetterN(3),
			Color:     teamColors[i],
			Color2:    teamSecondColors[i],
			Shield:    gofakeit.URL(),
			// NextMatchID will be set after matches are created
		}
	}

	if err := db.Create(&teams).Error; err != nil {
		return nil, err
	}

	return teams, nil
}

// seedSeasons creates 5 seasons in the database
func seedSeasons(db *gorm.DB) ([]model.Season, error) {
	seasons := make([]model.Season, 5)

	currentYear := time.Now().Year()

	for i := range seasons {
		year := uint16(currentYear - 4 + i)
		startDate := time.Date(int(year), 8, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(int(year)+1, 5, 31, 0, 0, 0, 0, time.UTC)

		seasons[i] = model.Season{
			Year:      year,
			StartDate: startDate,
			EndDate:   endDate,
			IsCurrent: i == 4, // The latest season is current
		}
	}

	if err := db.Create(&seasons).Error; err != nil {
		return nil, err
	}

	return seasons, nil
}

// seedPlayers creates 5 players in the database
func seedPlayers(db *gorm.DB, users []model.User) ([]model.Player, error) {
	players := make([]model.Player, 5)
	positions := []string{"por", "ceni", "cend", "lati", "med"}
	feet := []string{"L", "R"}

	for i := range players {
		// Assign a user to each player (optional relationship)
		var userID *string
		if i < len(users) {
			userID = &users[i].ID
		}

		players[i] = model.Player{
			NickName:      gofakeit.Username(),
			Height:        uint16(160 + rand.Intn(91)), // Height between 160-250 cm
			Country:       gofakeit.CountryAbr(),
			Country2:      gofakeit.CountryAbr(),
			Foot:          feet[rand.Intn(len(feet))],
			Age:           uint8(18 + rand.Intn(20)),
			SquadNumber:   uint8(1 + rand.Intn(99)),
			Rating:        uint8(50 + rand.Intn(51)),
			Matches:       uint16(rand.Intn(100)),
			YCards:        uint8(rand.Intn(10)),
			RCards:        uint8(rand.Intn(3)),
			Goals:         uint16(rand.Intn(50)),
			Assists:       uint16(rand.Intn(30)),
			Saves:         uint16(rand.Intn(100)),
			Position:      positions[i%len(positions)],
			Injured:       rand.Float32() < 0.2, // 20% chance of being injured
			CareerSummary: gofakeit.Paragraph(3, 5, 10, " "),
			MVPCount:      uint8(rand.Intn(10)),
			UserID:        userID,
		}
	}

	if err := db.Create(&players).Error; err != nil {
		return nil, err
	}

	return players, nil
}

// seedPlayerTeams creates player-team associations
func seedPlayerTeams(db *gorm.DB, players []model.Player, teams []model.Team, seasons []model.Season) error {
	for i := 0; i < 5; i++ {
		playerIndex := i % len(players)
		teamIndex := i % len(teams)
		seasonIndex := i % len(seasons)

		player := players[playerIndex]
		team := teams[teamIndex]
		season := seasons[seasonIndex]

		playerTeam := model.PlayerTeam{
			PlayerID:  player.ID,
			TeamID:    team.ID,
			SeasonID:  season.ID,
			StartDate: time.Now(),
			EndDate:   nil, // Explicitly set EndDate to NULL
		}

		if err := db.Create(&playerTeam).Error; err != nil {
			return fmt.Errorf("failed to create player_team for player %d, team %d, season %d: %w",
				player.ID, team.ID, season.ID, err)
		}
	}

	return nil
}

// seedTeamStats creates team statistics for each team in each season
func seedTeamStats(db *gorm.DB, teams []model.Team, seasons []model.Season) error {
	teamStats := make([]model.TeamStat, 0)

	for _, team := range teams {
		for _, season := range seasons {
			// Generate random stats for each team in each season
			wins := uint8(rand.Intn(15))
			draws := uint8(rand.Intn(10))
			losses := uint8(rand.Intn(15))
			goalsFor := uint16(rand.Intn(60))
			goalsAgainst := uint16(rand.Intn(40))

			teamStats = append(teamStats, model.TeamStat{
				TeamID:       team.ID,
				SeasonID:     season.ID,
				Wins:         wins,
				Draws:        draws,
				Losses:       losses,
				GoalsFor:     goalsFor,
				GoalsAgainst: goalsAgainst,
				Points:       uint16(wins*3 + draws),
				Rank:         uint8(1 + rand.Intn(10)),
			})
		}
	}

	if err := db.Create(&teamStats).Error; err != nil {
		return err
	}

	return nil
}

// seedMatches creates 5 matches in the database
func seedMatches(db *gorm.DB, teams []model.Team, seasons []model.Season, players []model.Player) ([]model.Match, error) {
	matches := make([]model.Match, 5)
	statuses := []string{"scheduled", "in_progress", "completed", "postponed", "cancelled"}
	locations := []string{"Home Stadium", "Away Field", "Central Arena", "Main Stadium", "City Field"}

	currentSeason := seasons[len(seasons)-1]

	for i := range matches {
		// Make sure home and away teams are different
		homeTeamIndex := i % len(teams)
		awayTeamIndex := (i + 1) % len(teams)

		// Random date within the current season
		matchDate := gofakeit.DateRange(currentSeason.StartDate, currentSeason.EndDate)

		// Random kickoff time (HH:MM)
		hour := rand.Intn(12) + 10  // Between 10:00 and 21:00
		minute := rand.Intn(4) * 15 // 00, 15, 30, or 45

		// Combine into full kickoff datetime
		kickoff := time.Date(
			matchDate.Year(),
			matchDate.Month(),
			matchDate.Day(),
			hour,
			minute,
			0, 0,
			time.UTC,
		)

		// For completed matches, set goals
		var homeGoals, awayGoals uint8
		if statuses[i%len(statuses)] == "completed" {
			homeGoals = uint8(rand.Intn(5))
			awayGoals = uint8(rand.Intn(5))
		}

		// Set MVP for completed matches
		var mvpPlayerID *uint64
		if statuses[i%len(statuses)] == "completed" {
			playerID := players[rand.Intn(len(players))].ID
			mvpPlayerID = &playerID
		}

		matches[i] = model.Match{
			Status:      statuses[i%len(statuses)],
			Kickoff:     kickoff,
			Location:    locations[i%len(locations)],
			HomeGoals:   homeGoals,
			AwayGoals:   awayGoals,
			HomeTeamID:  teams[homeTeamIndex].ID,
			AwayTeamID:  teams[awayTeamIndex].ID,
			SeasonID:    currentSeason.ID,
			MVPPlayerID: mvpPlayerID,
		}
	}

	if err := db.Create(&matches).Error; err != nil {
		return nil, err
	}

	return matches, nil
}

// seedLineups creates lineups for matches
func seedLineups(db *gorm.DB, players []model.Player, matches []model.Match) error {
	lineups := make([]model.Lineup, 0)
	positions := []string{"por", "ceni", "cend", "lati", "med", "latd", "del", "deli", "deld"}

	// Create 5 lineups
	for i := 0; i < 5; i++ {
		matchIndex := i % len(matches)
		playerIndex := i % len(players)

		lineups = append(lineups, model.Lineup{
			Position: positions[i%len(positions)],
			PlayerID: players[playerIndex].ID,
			MatchID:  matches[matchIndex].ID,
			Starting: rand.Float32() > 0.3, // 70% chance of being a starting player
		})
	}

	if err := db.Create(&lineups).Error; err != nil {
		return err
	}

	return nil
}

// seedArticles creates 5 articles in the database
func seedArticles(db *gorm.DB, seasons []model.Season) error {
	articles := make([]model.Article, 5)

	for i := range articles {
		seasonIndex := i % len(seasons)

		articles[i] = model.Article{
			Title:     gofakeit.Sentence(5),
			Content:   gofakeit.Paragraph(5, 10, 15, " "),
			ImgBanner: gofakeit.URL(),
			Date:      gofakeit.DateRange(seasons[seasonIndex].StartDate, seasons[seasonIndex].EndDate),
			SeasonID:  seasons[seasonIndex].ID,
		}
	}

	if err := db.Create(&articles).Error; err != nil {
		return err
	}

	return nil
}

// getCompletedMatches filters for completed matches or returns the first match if none are completed
func getCompletedMatches(matches []model.Match) []model.Match {
	completedMatches := make([]model.Match, 0)

	for _, match := range matches {
		if match.Status == "completed" {
			completedMatches = append(completedMatches, match)
		}
	}

	if len(completedMatches) == 0 && len(matches) > 0 {
		// If no completed matches, use the first match for demonstration purposes
		completedMatches = append(completedMatches, matches[0])
	}

	return completedMatches
}

// getShuffledPlayers returns a randomly shuffled copy of the players slice
func getShuffledPlayers(players []model.Player) []model.Player {
	shuffledPlayers := make([]model.Player, len(players))
	copy(shuffledPlayers, players)

	rand.Shuffle(len(shuffledPlayers), func(i, j int) {
		shuffledPlayers[i], shuffledPlayers[j] = shuffledPlayers[j], shuffledPlayers[i]
	})

	return shuffledPlayers
}

// generatePlayerStat creates a player stat with random but realistic values
func generatePlayerStat(player model.Player, match model.Match, seasonID uint64, isMVP bool) model.PlayerStat {
	positions := []string{"por", "ceni", "cend", "lati", "med", "latd", "del", "deli", "deld"}

	// Generate random stats
	goals := uint8(rand.Intn(3))   // 0-2 goals
	assists := uint8(rand.Intn(3)) // 0-2 assists
	saves := uint8(0)              // Default 0

	// Adjust stats based on position
	if player.Position == "por" { // If goalkeeper
		saves = uint8(rand.Intn(6)) // 0-5 saves
		goals = 0                   // No goals for goalkeepers typically
		assists = 0                 // No assists for goalkeepers typically
	}

	// Handle cards
	yellowCards := uint8(0)
	if rand.Float32() < 0.2 { // 20% chance of yellow card
		yellowCards = 1
	}

	redCards := uint8(0)
	if rand.Float32() < 0.05 { // 5% chance of red card
		redCards = 1
		yellowCards = 0 // If red card, no yellow (direct red)
	}

	// Calculate rating
	rating := uint8(50 + rand.Intn(51)) // 50-100 rating

	// Players with goals/assists/saves tend to have higher ratings
	if goals > 0 || assists > 0 || saves > 3 {
		rating = uint8(75 + rand.Intn(26)) // 75-100 rating
	}

	// MVP should have high rating
	if isMVP {
		rating = uint8(85 + rand.Intn(16)) // 85-100 rating
	}

	// Determine if starting and minutes played
	starting := rand.Float32() > 0.2 // 80% chance of starting

	var minutesPlayed uint8
	if starting {
		// Starters typically play between 60-90 minutes
		minutesPlayed = uint8(60 + rand.Intn(31))

		// If red card, likely played fewer minutes
		if redCards > 0 {
			minutesPlayed = uint8(rand.Intn(70) + 1) // 1-70 minutes
		}
	} else {
		// Substitutes typically play between 1-30 minutes
		minutesPlayed = uint8(1 + rand.Intn(30))
	}

	// Randomly decide if we want to set a TeamID (for demonstration of the feature)
	var teamID *uint64
	if rand.Float32() > 0.5 { // 50% chance of setting team ID
		// Randomly assign either home or away team
		if rand.Float32() > 0.5 {
			teamID = &match.HomeTeamID
		} else {
			teamID = &match.AwayTeamID
		}
	}

	return model.PlayerStat{
		PlayerID:      player.ID,
		MatchID:       match.ID,
		SeasonID:      seasonID,
		TeamID:        teamID, // Added TeamID field
		Goals:         goals,
		Assists:       assists,
		Saves:         saves,
		YellowCards:   yellowCards,
		RedCards:      redCards,
		Rating:        rating,
		IsStarting:    starting,
		MinutesPlayed: minutesPlayed,
		IsMVP:         isMVP,
		Position:      positions[rand.Intn(len(positions))], // Random position
	}
}

// updateTeamsNextMatch updates each team with a reference to an upcoming match
func updateTeamsNextMatch(db *gorm.DB, teams []model.Team, matches []model.Match) error {
	// Find matches that are scheduled
	scheduledMatches := make([]model.Match, 0)
	for _, match := range matches {
		if match.Status == "scheduled" {
			scheduledMatches = append(scheduledMatches, match)
		}
	}

	// If no scheduled matches, use any match
	if len(scheduledMatches) == 0 {
		scheduledMatches = matches
	}

	// For each team, assign a NextMatchID
	for i, team := range teams {
		// Pick a match for this team (could be home or away)
		matchIndex := i % len(scheduledMatches)
		match := scheduledMatches[matchIndex]

		// Set NextMatchID
		matchID := match.ID

		// Update the team in the database
		if err := db.Model(&model.Team{}).Where("id = ?", team.ID).Update("next_match_id", matchID).Error; err != nil {
			return fmt.Errorf("failed to update NextMatchID for team %d: %w", team.ID, err)
		}
	}

	return nil
}

// seedPlayerStats creates player statistics for matches
func seedPlayerStats(db *gorm.DB, players []model.Player, matches []model.Match, seasons []model.Season) error {
	// Get completed matches
	completedMatches := getCompletedMatches(matches)
	if len(completedMatches) == 0 {
		return nil // No matches to generate stats for
	}

	playerStats := make([]model.PlayerStat, 0)

	// For each completed match, create stats for some players
	for _, match := range completedMatches {
		// Get shuffled players
		shuffledPlayers := getShuffledPlayers(players)

		// Find the season for this match
		var seasonID uint64
		for _, season := range seasons {
			if season.ID == match.SeasonID {
				seasonID = season.ID
				break
			}
		}

		// Create stats for the selected players
		numPlayersToGenerate := min(5, len(shuffledPlayers))
		for i := 0; i < numPlayersToGenerate; i++ {
			player := shuffledPlayers[i]

			// Determine if player was MVP
			isMVP := false
			if match.MVPPlayerID != nil && *match.MVPPlayerID == player.ID {
				isMVP = true
			}

			// Generate the player stat
			playerStat := generatePlayerStat(player, match, seasonID, isMVP)
			playerStats = append(playerStats, playerStat)
		}
	}

	// Save all player stats to the database
	if err := db.Create(&playerStats).Error; err != nil {
		return fmt.Errorf("failed to create player stats: %w", err)
	}

	return nil
}
