package internal

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

func Run() {
	LoadEnv()

	// About `search_path` param...
	// - https://github.com/jackc/pgx/issues/1013
	// The repo author suggests to pass query `search_path` param. The OP kinda
	// indicates the solution works on pgx + pgxpool setup.
	// In pgx + sqlx db connection setup, however, it seems that the solution
	// doesn't work. Not sure about the actual technical details, tho

	// - https://github.com/riverqueue/river/issues/15
	// Other folks indicates that the user of this repo could just explicitly
	// pass its value to `RuntimeParams`
	dbCfg, err := pgx.ParseConfig(dB_URL)
	dbCfg.RuntimeParams["search_path"] = dB_SCHEMA
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	dbConn := sqlx.NewDb(stdlib.OpenDB(*dbCfg), "pgx")
	defer dbConn.Close()

	app := chi.NewRouter()
	store := persistence.NewPg(dbConn)
	service := service.NewService(&store)
	controller := controller.NewController(service, aLPINE_URL, hTMX_URL)

	app.Use(middleware.RequestID)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	route.NewRouter(controller).UseOn(app)

	port := 10000
	fmt.Printf(fmt.Sprintf("Server listening on port %d...\n", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), app); err != nil {
		log.Fatalf("server listening: %v", err)
	}
}
