package api
import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelreservation/db"
	"context"
	"testing"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) tearDown(t *testing.T) {
	if err := tdb.client.Database(db.TestDB).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setupUserHandlerTest(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.TestURI))
	if err != nil {
		t.Fatal(err)
	}
	return &testdb{
		client: client,
		Store: &db.Store{
			User: db.NewMongoUserStore(client),
		},
	}
}