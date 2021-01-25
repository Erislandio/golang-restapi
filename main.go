package main

import (
	"github.com/erislandio/web/restapi/database"
	"github.com/erislandio/web/restapi/routes"
)

func main() {
	database.Init()
	routes := routes.Init()

	routes.Logger.Fatal(routes.Start(":8080"))

}
