package api

import (
	"context"
	"fmt"
	"github.com/egylinux/hr/employees"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
)

type EmployeeManager interface {
	AddEmployee(emp employees.Employee) (bool, error)
	AddEmployees(emp []interface{}) (bool, error)
	Find(filter interface{}) ([]*employees.Employee, error)
	//FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	Delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (bool, error)
}

type EmployeeRouter struct {
	manager EmployeeManager
	ctx     context.Context
}

func NewRouter(manager EmployeeManager, ct context.Context) *echo.Echo {
	router := EmployeeRouter{manager: manager, ctx: ct}
	r := echo.New()
	r.Add(http.MethodPost, "/newemp", router.AddEmployee)
	r.Add(http.MethodGet, "/find", router.GetAll)
	r.Add(http.MethodGet, "/search/:id", router.Search)
	r.Add(http.MethodDelete, "/delete/:id", router.Delete)
	return r
}
func (e *EmployeeRouter) AddEmployee(c echo.Context) error {

	emp := new(employees.Employee)
	if err := c.Bind(emp); err != nil {
		fmt.Println("Error", err.Error())
		return err
	}

	_, err := e.manager.AddEmployee(*emp)
	//fmt.Println(exist)
	if err != nil {
		return c.String(http.StatusOK, "Error while saving , "+err.Error())
	}

	return c.String(http.StatusOK, "Saved Successfully")
}
func (e *EmployeeRouter) GetAll(c echo.Context) error {

	filter := bson.D{}
	employees, err := e.manager.Find(filter)
	fmt.Println("success len : ", len(employees))
	if err != nil {
		return c.String(http.StatusOK, "Error while reading data , "+err.Error())
	}
	fmt.Println("emp1")
	fmt.Println(employees[0].Employeename)
	name := ""
	for i := 0; i < len(employees); i++ {
		//fmt.Println(employees[i])
		name += employees[i].Employeename + "\n"
	}

	return c.JSON(http.StatusOK, employees)

}

func (e *EmployeeRouter) Search(c echo.Context) error {

	id := c.Param("id")
	intid, _ := strconv.Atoi(id)
	filter := bson.D{{"id", intid}}
	employees, err := e.manager.Find(filter)

	if err != nil {
		return c.String(http.StatusOK, "Error while reading data , "+err.Error())
	}
	if employees==nil {
		return c.String(http.StatusOK,"Employee Not Found")
	}

	name := ""
	for i := 0; i < len(employees); i++ {
		name += employees[i].Employeename + "\n"
	}

	return c.JSON(http.StatusOK, employees)

}

func (e *EmployeeRouter) Delete(c echo.Context) error {
	fmt.Println("start deleting")
	id := c.Param("id")
	intid, _ := strconv.Atoi(id)
	fmt.Printf("id: %v\n", intid)
	filter := bson.D{{"id", intid}}
	res, err := e.manager.Delete(e.ctx, filter)

	if err != nil {
		return c.String(http.StatusOK, "Error :"+err.Error())
	}
	fmt.Println(res)
	return c.String(http.StatusOK, "Employee Deleted Successfully")

}
