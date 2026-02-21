package main

import (
	"log"

	"github.com/afifudin23/absensi-king-royal-api/internal/config"
	"github.com/afifudin23/absensi-king-royal-api/internal/delivery/http/router"
)

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	r := router.New()
	log.Printf("starting %s on %s", env.AppName, env.Port)
	if err := r.Run(env.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
