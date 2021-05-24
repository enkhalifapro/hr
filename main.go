package main

import (
	"context"
	"fmt"
	"github.com/egylinux/hr/api"
	"github.com/egylinux/hr/employees"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")
	// get collection as ref
	collection := client.Database("hrdb").Collection("employees")

	manager := employees.NewManager(collection)
	/*

		_, err = manager.AddEmployee(employees.Employee{Id: 11, Employeename: "Ahmed Mahmoud"})
		if err != nil {
			fmt.Println("Error Inserting Employee ", err.Error())
		}else {
			fmt.Println("inserted successfully")
		}



	*/


	router := api.NewRouter(manager)
	router.Add(http.MethodGet, "/", home)
	router.Logger.Fatal(router.Start(":1323"))

	defer client.Disconnect(ctx)
}
// Handler
func home(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome, HR Manager")
}