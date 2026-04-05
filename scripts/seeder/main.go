package main

import (
	"log"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/database/seeder"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("failed init config: %v", err)
	}

	log.Println("sedding start...")

	seeder.SeedUsers()
	seeder.SeedPayrollSettings()

	log.Println("sedding finish...")
}
