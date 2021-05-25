package employees

import (
	"context"
	"errors"
	"github.com/egylinux/hr/employees/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddEmploy(t *testing.T) {
	//Arrange
	db := &mocks.DBConnector{}
	ctx := context.Background()
	emp := Employee{Id: 5, Employeename: "ahmed", Age: 35, Active: true}
	db.On("InsertOne", ctx, emp).Return(nil, nil)
	emp0 := Employee{Id: 6, Employeename: "ahmed", Age: 35, Active: true}
	db.On("InsertOne", ctx, emp0).Return(nil, errors.New("Inserting Error"))
	//Act
	mgr := NewManager(db, ctx)

	emp2 := Employee{Id: 5, Employeename: "ahmed", Age: 35, Active: true}
	ok, err := mgr.AddEmployee(emp2)
	//Assert
	assert.Equal(t,true ,ok )
	assert.Equal(t, nil, err)

	//Act


	emp3 := Employee{Id: 6, Employeename: "ahmed", Age: 35, Active: true}
	ok, err = mgr.AddEmployee(emp3)
	//Assert
	assert.Equal(t,false ,ok )
	assert.Error(t, err)

}
