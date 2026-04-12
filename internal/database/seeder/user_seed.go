package seeder

import (
	"log"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/model"
	"github.com/afifudin23/absensi-king-royal-api/pkg/utils"
	"github.com/google/uuid"
)

func SeedUsers() {
	db := config.GetDB()

	adminPassword, err := utils.HashPassword("Admin123")
	if err != nil {
		log.Printf("Failed hash password for admin@kingroyal.com: %v\n", err)
		return
	}

	userPassword, err := utils.HashPassword("User123")
	if err != nil {
		log.Printf("Failed hash password for user@kingroyal.com: %v\n", err)
		return
	}

	users := []model.User{
		{FullName: "Admin", Email: "admin@kingroyal.com", Password: adminPassword, Role: "admin"},
		{FullName: "User", Email: "user@kingroyal.com", Password: userPassword, Role: "user"},
	}

	for _, user := range users {
		var existing model.User
		result := db.Where("email = ?", user.Email).Find(&existing)
		if result.Error != nil {
			log.Printf("Failed find user %s: %v\n", user.Email, result.Error)
			continue
		}

		if result.RowsAffected == 0 {
			user.ID = uuid.NewString()

			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed seed user %s: %v\n", user.Email, err)
				continue
			}

			existing = user
		}

		// Ensure user profile exists for the created user
		var profile model.UserProfile
		result = db.Where("user_id = ?", existing.ID).Find(&profile)
		if result.Error != nil {
			log.Printf("Failed find user_profile for %s: %v\n", user.Email, result.Error)
			continue
		}

		if result.RowsAffected == 0 {
			profile = model.UserProfile{
				ID:     uuid.NewString(),
				UserID: existing.ID,
			}

			if err := db.Create(&profile).Error; err != nil {
				log.Printf("Failed seed user_profile for %s: %v\n", user.Email, err)
				continue
			}
		}
	}
}
