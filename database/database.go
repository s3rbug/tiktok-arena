package database

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"tiktok-arena/configuration"
	"tiktok-arena/models"
)

var DB *gorm.DB

func ConnectDB(config *configuration.EnvConfigModel) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost,
		config.DBUserName,
		config.DBUserPassword,
		config.DBName,
		config.DBPort,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database!\n", err.Error())
	}
	//	Extension for postgresql uuid support
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	//  Extension for search
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"fuzzystrmatch\"")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Tournament{},
		&models.Tiktok{},
	)
	if err != nil {
		log.Fatal("Migration Failed:\n", err.Error())
	}
	log.Println("Successfully connected to the database")
}

func GetUserByName(username string) (models.User, error) {
	var user models.User
	record := DB.Table("users").First(&user, "name = ?", username)
	return user, record.Error
}

func UserExists(username string) bool {
	var user models.User
	DB.Table("users").Select("id").First(&user, "name = ?", username)
	return user.ID != nil
}

func CreateNewUser(newUser *models.User) error {
	record := DB.Table("users").Create(&newUser)
	return record.Error
}

func GetTournamentById(tournamentId string) (models.Tournament, error) {
	var tournament models.Tournament
	record := DB.Table("tournaments").First(&tournament, "id = ?", tournamentId)
	return tournament, record.Error
}

func CheckIfTournamentExistsByName(name string) bool {
	var tournament models.Tournament
	DB.Table("tournaments").Select("id").First(&tournament, "name = ?", name)
	return tournament.ID != nil
}

func CheckIfNameIsTakenByOtherTournament(name string, id uuid.UUID) bool {
	var tournament models.Tournament
	DB.Table("tournaments").Select("id").First(&tournament, "name = ? AND id != ?", name, id)
	return tournament.ID != nil
}

func CheckIfTournamentExistsById(id uuid.UUID) bool {
	var tournament models.Tournament
	DB.Table("tournaments").Select("id").First(&tournament, "id = ?", id)
	return tournament.ID != nil
}

func CheckIfTournamentsExistsByIds(ids []string, userId uuid.UUID) (bool, error) {
	var tournaments []models.Tournament
	record := DB.Table("tournaments").Where("user_id = ? AND id IN ?", userId, ids).Find(&tournaments)
	if len(tournaments) != len(ids) {
		return false, record.Error
	}
	return true, record.Error
}

func CreateNewTournament(newTournament *models.Tournament) error {
	record := DB.Table("tournaments").Create(&newTournament)
	return record.Error
}

func EditTournament(t *models.Tournament) error {
	record := DB.Table("tournaments").Where("id = ?", &t.ID).
		Updates(map[string]interface{}{"id": &t.ID, "name": &t.Name, "size": &t.Size})
	return record.Error
}

func DeleteTournamentById(id uuid.UUID, userId uuid.UUID) error {
	record := DB.Table("tournaments").Where("id = ? AND user_id = ?", id, userId).Delete(&models.Tournament{})
	return record.Error
}

func DeleteTournamentsByIds(ids []string, userId uuid.UUID) error {
	record := DB.Table("tournaments").
		Where("user_id = ? AND id IN (?)", userId, ids).
		Delete(&models.Tournament{})
	return record.Error
}

func DeleteTiktoksByIds(ids []string) error {
	record := DB.Table("tiktoks").Where("tournament_id IN (?)", ids).Delete(&models.Tiktok{})
	return record.Error
}

func CreateNewTiktok(newTiktok *models.Tiktok) error {
	record := DB.Table("tiktoks").Create(&newTiktok)
	return record.Error
}

func EditTiktok(t *models.Tiktok) error {
	record := DB.Table("tiktoks").Where(&t.URL, &t.TournamentID).Updates(&t)
	return record.Error
}

func DeleteTiktoks(t []models.Tiktok) error {
	record := DB.Table("tiktoks").Delete(t)
	return record.Error
}

func CreateNewTiktoks(t []models.Tiktok) error {
	record := DB.Table("tiktoks").Create(t)
	return record.Error
}

func GetTournamentTiktoksById(tournamentId uuid.UUID) ([]models.Tiktok, error) {
	var tiktoks []models.Tiktok
	record := DB.Table("tiktoks").
		Select("*").
		Find(&tiktoks, "tournament_id = ?", tournamentId)
	return tiktoks, record.Error
}

func GetTournaments(queries models.PaginationQueries) (models.TournamentsResponse, error) {
	var tournaments []models.Tournament
	var totalTournaments int64
	DB.Table("tournaments").Count(&totalTournaments)
	record := DB.Table("tournaments").
		Scopes(Search(queries.SearchText)).
		Scopes(Sort(queries.SortName, queries.SortSize)).
		Scopes(Paginate(queries.Page, queries.Count)).
		Find(&tournaments)
	return models.TournamentsResponse{TournamentCount: totalTournaments, Tournaments: tournaments}, record.Error
}

func Search(searchText string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if searchText == "" {
			return db
		}
		return db.Select("*, levenshtein(name, ?) as distance", searchText).Order("distance")
	}
}

func Sort(sortName string, sortSize string) func(db *gorm.DB) *gorm.DB {
	var sortString string
	if sortName != "" && sortSize == "" {
		sortString = fmt.Sprintf("name %s", sortName)
	}
	if sortName == "" && sortSize != "" {
		sortString = fmt.Sprintf("size %s", sortSize)
	}
	if sortName != "" && sortSize != "" {
		sortString = fmt.Sprintf("name %s, size %s", sortName, sortSize)
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(sortString)
	}
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	offset := (page - 1) * pageSize
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(pageSize)
	}
}

func GetAllTournamentsForUserById(id uuid.UUID, queries models.PaginationQueries) (models.TournamentsResponse, error) {
	var tournaments []models.Tournament
	var totalTournaments int64
	DB.Table("tournaments").Where("user_id = ?", id).Count(&totalTournaments)
	record := DB.Table("tournaments").
		Scopes(Search(queries.SearchText)).
		Where("user_id = ?", id).
		Scopes(Sort(queries.SortName, queries.SortSize)).
		Scopes(Paginate(queries.Page, queries.Count)).
		Limit(100).Find(&tournaments)
	return models.TournamentsResponse{TournamentCount: totalTournaments, Tournaments: tournaments}, record.Error
}

func RegisterTiktokWinner(tournamentId uuid.UUID, tiktokURL string) error {
	record := DB.Table("tiktoks").Where("tournament_id = ? AND url = ?", tournamentId, tiktokURL).
		UpdateColumn("wins", gorm.Expr("wins + ?", 1))
	if record.Error != nil {
		return record.Error
	}
	record = DB.Table("tournaments").Where("id = ?", tournamentId).
		UpdateColumn("times_played", gorm.Expr("times_played + ?", 1))
	return record.Error
}
