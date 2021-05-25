package main

import (
	"context"
	"fmt"
	"github.com/egylinux/hr/api"
	"github.com/egylinux/hr/employees"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)



func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx :=  context.Background()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")
	// get collection as ref
	collection := client.Database("hrdb").Collection("employees")

	manager := employees.NewManager(collection,ctx)

	router := api.NewRouter(manager,ctx)
	router.Add(http.MethodGet, "/", home)
	router.Logger.Fatal(router.Start(":1323"))

	 if err=client.Disconnect(ctx);err != nil{
	 	fmt.Println(err.Error())
	 }

}
// Handler
func home(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome, HR Manager")
}