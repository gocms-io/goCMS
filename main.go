package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/gin-gonic/gin"
	"github.com/menklab/goCMS/context"
	"github.com/menklab/goCMS/controllers"
	"github.com/menklab/goCMS/database"
	"github.com/menklab/goCMS/repositories"
	"github.com/menklab/goCMS/services"
	"log"
)

var app *Engine

type Engine struct {
	Gin               *gin.Engine
	ControllersGroup  *controllers.ControllersGroup
	ServicesGroup     *services.ServicesGroup
	RepositoriesGroup *repositories.RepositoriesGroup
	Database          *database.Database
}

//go:generate apidoc -c ./ -i ./models -i ./controllers/ -o ./docs/ -f ".*\\.go$" -f ".*\\.js$"
//go:generate go run ./utility/gen/gen.docs.go -o ./docs/docs.js ./models ./controllers/

func Default() *Engine {

	// init config environment vars
	context.Init()

	// setup database
	db := database.Default()

	// migrate cms db
	db.MigrateCMSSql()

	// start gin with defaults
	r := gin.Default()
	// setup repositories
	rg := repositories.DefaultRepositoriesGroup(db)

	// setup services
	sg := services.DefaultServicesGroup(rg)

	// setup controllers
	cg := controllers.DefaultControllerGroup(r, sg)

	// create engine
	engine := Engine{
		Gin:               r,
		ControllersGroup:  cg,
		ServicesGroup:     sg,
		RepositoriesGroup: rg,
		Database:          db,
	}
	return &engine
}

func (engine *Engine) Listen(uri string) {

	err := engine.Gin.Run(uri)
	log.Println(err.Error())

}

func main() {

	// startup defaults
	app = Default()

	// start server and listen
	port := context.Config.Port

	if port == "" {
		port = "8080"
	}

	app.Listen(":" + port)
}


