package api

import (
	"fmt"
	"github.com/egylinux/hr/employees"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)
type EmployeeManager interface {
	AddEmployee(emp employees.Employee) (bool, error)
	AddEmployees(emp []interface{}) (bool, error)
	Find(filter interface{}) ([]*employees.Employee, error)
}

type EmployeeRouter struct {
	manager EmployeeManager
}

func NewRouter(manager EmployeeManager) *echo.Echo {
	router := EmployeeRouter{manager: manager}
	r := echo.New()
	r.Add(http.MethodPost, "/newemp", router.AddEmployee)
	r.Add(http.MethodGet, "/find", router.GetAll)
	return r
}
func (e *EmployeeRouter) AddEmployee(c echo.Context) error {

	emp:= new(employees.Employee)
	if err := c.Bind(emp); err != nil {
		fmt.Println("Error",err.Error())
		return err
	}


	_, err := e.manager.AddEmployee(*emp)
	//fmt.Println(exist)
	if err != nil {
		return c.String(http.StatusOK, "Error while saving , "+err.Error())
	}

	return c.String(http.StatusOK, "Saved Successfully")
}
func (e *EmployeeRouter) GetAll(c echo.Context)   error {

	filter := bson.D{}
	employees, err := e.manager.Find(filter)
fmt.Println("success len : ",len(employees))
	if err != nil {
		return c.String(http.StatusOK, "Error while reading data , "+err.Error())
	}
fmt.Println("emp1")
	fmt.Println(employees[0].Employeename)
	name:=""
	for i := 0; i < len(employees); i++ {
		//fmt.Println(employees[i])
		name+= employees[i].Employeename +"\n"
	}

	return c.JSON(http.StatusOK, employees)

}