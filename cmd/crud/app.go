package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/solsteace/misite/internal/controller"
	"github.com/solsteace/misite/internal/persistence"
	"github.com/solsteace/misite/internal/service"
)

type appState int

const (
	sTATE_READY appState = iota
	sTATE_NEED_ARG
	sTATE_OVER
)

const (
	fLAG_SOURCEFILE = "--source"
	fLAG_TARGET     = "--target"
	fLAG_ENTITY     = "--entity"
	fLAG_ACTION     = "--action"
	fLAG_HELP       = "--help"
)

func main() {
	args := os.Args
	state := sTATE_READY

	var sourceFile string
	var target string
	var entity string
	var action string
	var lastFlag string
	for _, arg := range args {
		switch state {
		case sTATE_READY:
			switch arg {
			case fLAG_TARGET, fLAG_SOURCEFILE, fLAG_ENTITY, fLAG_ACTION:
				state = sTATE_NEED_ARG
				lastFlag = arg
			case fLAG_HELP:
				state = sTATE_OVER
			}
		case sTATE_NEED_ARG:
			switch lastFlag {
			case fLAG_SOURCEFILE:
				sourceFile = arg
			case fLAG_TARGET:
				target = arg
			case fLAG_ENTITY:
				entity = arg
			case fLAG_ACTION:
				action = arg
			}
			state = sTATE_READY
		case sTATE_OVER:
			fmt.Println(tEXT_HELP)
			os.Exit(0)
		}
	}

	if len(args) == 1 {
		fmt.Println(tEXT_HELP)
		os.Exit(0)
	}
	switch state {
	case sTATE_NEED_ARG:
		switch lastFlag {
		case fLAG_ACTION:
			log.Fatalf("missing action argument")
		case fLAG_ENTITY:
			log.Fatalf("missing entity argument")
		case fLAG_SOURCEFILE:
			log.Fatalf("missing data source file argument")
		case fLAG_TARGET:
			log.Fatalf("missing target argument")
		}
	}
	switch "" {
	case entity:
		log.Fatalf("missing entity argument")
	case sourceFile:
		log.Fatalf("missing data source file argument")
	case target:
		log.Fatalf("missing target argument")
	case action:
		log.Fatalf("missing action argument")
	}

	dbCfg, err := pgx.ParseConfig(target)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	dbConn := sqlx.NewDb(stdlib.OpenDB(*dbCfg), "pgx")
	defer dbConn.Close()

	db := persistence.NewPg(dbConn)
	service := service.NewService(&db)
	controller := controller.NewController(service, "", "", "")

	var handler func(*os.File) error
	switch entity {
	case "a", "articles":
		switch action {
		case "a", "add":
			handler = controller.InsertArticles
		case "u", "update":
			handler = controller.UpsertArticles
		case "d", "delete":
			handler = controller.DeleteArticles
		}
	case "at", "article_tags":
		switch action {
		case "a", "add":
			handler = controller.InsertArticleTags
		case "u", "update":
			handler = controller.UpsertArticleTags
		case "d", "delete":
			handler = controller.DeleteArticleTags
		}
	case "p", "projects":
		switch action {
		case "a", "add":
			handler = controller.InsertProjects
		case "u", "update":
			handler = controller.UpsertProjects
		case "d", "delete":
			handler = controller.DeleteProjects
		}
	case "pt", "project_tags":
		switch action {
		case "a", "add":
			handler = controller.InsertProjectTags
		case "u", "update":
			handler = controller.UpsertProjectTags
		case "d", "delete":
			handler = controller.DeleteProjectTags
		}
	case "pl", "project_links":
		switch action {
		case "a", "add":
			handler = controller.InsertProjectLinks
		case "u", "update":
			handler = controller.UpsertProjectLinks
		case "d", "delete":
			handler = controller.DeleteProjectLinks
		}
	case "t", "tags":
		switch action {
		case "a", "add":
			handler = controller.InsertTags
		case "u", "update":
			handler = controller.UpsertTags
		case "d", "delete":
			handler = controller.DeleteTags
		}
	case "s", "series":
		switch action {
		case "a", "add":
			handler = controller.InsertSeries
		case "u", "update":
			handler = controller.UpsertSeries
		case "d", "delete":
			handler = controller.DeleteSeries
		}
	}
	if handler == nil {
		log.Fatalf("unknown entity or handler type")
	}

	f, err := os.Open(sourceFile)
	if err != nil {
		log.Fatalf("opening data file: %s", err.Error())
	}
	if err := handler(f); err != nil {
		log.Fatalf("handling action: %s", err.Error())
	}
}

const tEXT_HELP = `NOTE ============
* -> mandatory flags

FLAGS ===========
help - do you need help?

*action - what do you want to do?
- (a)dd
- (u)pdate
- (d)pdate

*source - where the app should look the data from to do the action?

*target - where the action should be applied to?

*entity - to what object the action should be applied to?
- (a)rticles
- (a)rticle_(t)ags
- (p)rojects
- (p)roject_(t)ags
- (p)roject_(l)inks
- (t)ags
- (s)eries `
