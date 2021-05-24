package employees

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
}
// DBConnector contains dataAccess functionalities
type DBConnector interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	/*FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*UpdateResult, error)
	UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*UpdateResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*DeleteResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*DeleteResult, error)

	 */
}
// NewManager instantiate employee manager
func NewManager(connector DBConnector) *Manager {
	return &Manager{
		connector: connector,
	}
}
func (m *Manager) AddEmployee(emp Employee)  (bool, error) {

	_,err:= m.connector.InsertOne(context.TODO(),emp)
	if err != nil{
		return  true,nil
	}else{
		return  false,err
	}
}
func (m *Manager) AddEmployees(emp []interface{})  (bool, error) {

	_,err:= m.connector.InsertMany(context.TODO(),emp)
	if err != nil{
		return  true,nil
	}else{
		return  false,err
	}
}
func (m *Manager) Find(filter interface{})  ([]*Employee, error) {

	 collec,err := m.connector.Find(context.TODO(), filter)

		var results []*Employee


		for collec.Next(context.TODO()) {
			var elem Employee
			err := collec.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}

			results = append(results, &elem)
		}


	if err != nil{

		return  results,nil
	}else{
		return  results,err
	}
}