package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/solsteace/misite/internal/controller"
	"github.com/solsteace/misite/internal/persistence"
	"github.com/solsteace/misite/internal/route"
	"github.com/solsteace/misite/internal/service"
)

func main() {
	loadEnv()

	dbCfg, err := pgx.ParseConfig(dB_URL)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	dbConn := sqlx.NewDb(stdlib.OpenDB(*dbCfg), "pgx")
	defer dbConn.Close()

	app := chi.NewRouter()
	store := persistence.NewPg(dbConn)
	service := service.NewService(&store)
	controller := controller.NewController(
		service,
		iNDEX_URL,
		aLPINE_URL,
		hTMX_URL)

	app.Use(middleware.RequestID)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	route.NewRouter(controller).UseOn(app)

	port := 10000
	fmt.Printf("Server listening on port %d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), app); err != nil {
		log.Fatalf("server listening: %v", err)
	}
}
