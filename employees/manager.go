package employees

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Employee Struct to set Employee Info
type Employee struct {
	Id           int
	Employeename string
	Age          int
	Active       bool
}

// Manager to manage Employee Data
type Manager struct {
	connector DBConnector
	ctx       context.Context
}

// DBConnector contains dataAccess functionalities
type DBConnector interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	/*UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*UpdateResult, error)
	UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*UpdateResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*DeleteResult, error)


	*/
}

// NewManager instantiate employee manager
func NewManager(connector DBConnector, ct context.Context) *Manager {
	return &Manager{
		connector: connector,
		ctx:       ct,
	}
}
func (m *Manager) AddEmployee(emp Employee) (bool, error) {

	_, err := m.connector.InsertOne(m.ctx, emp)
	if err != nil {
		return false,err
	}
	return true, nil
}
func (m *Manager) AddEmployees(emp []interface{}) (bool, error) {

	_, err := m.connector.InsertMany(m.ctx, emp)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (m *Manager) Delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (bool, error) {

	res, err := m.connector.DeleteOne(m.ctx, filter)
	if err != nil {
		return false, err
	}
	fmt.Println(res)
	return true, nil

}
func (m *Manager) Find(filter interface{}) ([]*Employee, error) {

	collec, err := m.connector.Find(m.ctx, filter)
	if err!=nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var results []*Employee

	for collec.Next(m.ctx) {
		var elem Employee
		err := collec.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	if err != nil {

		return results, nil
	} else {
		return results, err
	}
}
