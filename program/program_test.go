package program

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// getMockDB returns a mock *sql.DB connection
// and a mock object generated using sqlmock library
func getMockDB() (*sql.DB, sqlmock.Sqlmock) {

	db, mock, err := sqlmock.New()

	if err != nil {
		log.Fatalf("An error encountered when generating new mock connection: %s", err)
	}

	return db, mock
}

func TestSaveCreds(t *testing.T) {

}

// the mockStore struct doesn't need anything except the Conn field
type mockStore struct {
	Conn *sql.DB
}

// The passed Conn argument is our mock DB connection
func newMockStore(Conn *sql.DB) *mockStore {
	return &mockStore{
		Conn: Conn,
	}
}

func (ms *mockStore) SaveCreds(encryptionKey []byte) error {

	return nil

}

func (ms *mockStore) RetrieveCreds(string, string, []byte) ([]map[string]string, error) {
	return nil, nil
}

func (ms *mockStore) ViewCreds(string, []byte) error {
	return nil
}

func (ms *mockStore) EditCreds(string, []byte) error {
	return nil
}

func (ms *mockStore) DeleteCreds(string, []byte) error {
	return nil
}
