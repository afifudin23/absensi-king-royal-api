package main

import (
	"log"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/router"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("failed to initialize app context: %v", err)
	}
	defer func() {
		if err := config.CloseDB(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	env := config.GetEnv()
	if env == nil {
		log.Fatalf("failed to read app env")
	}

	r := router.New()
	log.Printf("starting %s on %s", env.AppName, env.Port)
	if err := r.Run(env.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
