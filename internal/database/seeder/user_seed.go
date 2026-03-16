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
		{ID: uuid.NewString(), FullName: "Admin", Email: "admin@kingroyal.com", Password: adminPassword, Role: "admin"},
		{ID: uuid.NewString(), FullName: "User", Email: "user@kingroyal.com", Password: userPassword, Role: "user"},
	}

	for _, user := range users {
		var existing model.User
		err := db.Where("email = ?", user.Email).Attrs(user).FirstOrCreate(&existing).Error
		if err != nil {
			log.Printf("Failed seed user %s: %v\n", user.Email, err)
			continue
		}

		var profile model.UserProfile
		err = db.Where("user_id = ?", existing.ID).Attrs(model.UserProfile{UserID: existing.ID}).FirstOrCreate(&profile).Error
		if err != nil {
			log.Printf("Failed seed user_profile for %s: %v\n", user.Email, err)
		}
	}
}
