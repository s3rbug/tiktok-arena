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

func CheckIfUserExists(username string) bool {
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

func CreateNewTournament(newTournament *models.Tournament) error {
	record := DB.Table("tournaments").Create(&newTournament)
	return record.Error
}

func EditTournament(t *models.Tournament) error {
	record := DB.Table("tournaments").Where("id = ?", &t.ID).
		Updates(map[string]interface{}{"id": &t.ID, "name": &t.Name, "size": &t.Size})
	return record.Error
}

func DeleteTournamentById(id uuid.UUID) error {
	record := DB.Table("tournaments").Where("id = ?", id).Delete(&models.Tournament{})
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

func GetTournamentTiktoksById(tournamentId string) ([]models.Tiktok, error) {
	var tiktoks []models.Tiktok
	record := DB.Table("tiktoks").
		Select([]string{"TournamentID", "URL", "Wins", "AvgPoints"}).
		Find(&tiktoks, "tournament_id = ?", tournamentId)
	return tiktoks, record.Error
}

func GetAllTournaments() ([]models.Tournament, error) {
	var tournaments []models.Tournament
	record := DB.Table("tournaments").Select("*").Limit(100).Find(&tournaments)
	return tournaments, record.Error
}
