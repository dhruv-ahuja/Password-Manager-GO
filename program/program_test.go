package program

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
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

	db, mock := getMockDB()
	defer db.Close()

	ms := newMockStore(db)

	p := New()
	p.store = ms

	mock.ExpectExec("INSERT INTO info").
		WithArgs("reddit", "test123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := p.store.SaveCreds([]byte("ok")); err != nil {
		t.Errorf("error when running func: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}

}

func TestRetrieveCreds(t *testing.T) {

	query := "SELECT * FROM info WHERE key ILIKE ? ORDER BY id ASC"

	db, mock := getMockDB()
	defer db.Close()

	ms := newMockStore(db)

	p := New()
	p.store = ms

	mockRows := sqlmock.NewRows([]string{"id", "key", "encrypted_pw"}).AddRow("1", "reddit", "test123").AddRow("2", "Reddit", "testing")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM info WHERE key ILIKE ? ORDER BY id ASC`)).WithArgs("%red%").WillReturnRows(mockRows)

	_, err := p.store.RetrieveCreds(query, "red", []byte("test"))

	if err != nil {
		t.Log(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations unmet: %s", err)
	}

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

	usrInfo := make(map[string]string, 2)
	usrInfo["key"], usrInfo["password"] = "reddit", "test123"

	query := "INSERT INTO info (key, encrypted_pw) VALUES (?, ?)"

	if _, err := ms.Conn.Exec(query,
		usrInfo["key"], usrInfo["password"]); err != nil {
		return err
	}

	return nil

}

func (ms *mockStore) RetrieveCreds(query string, key string, encryptionKey []byte) ([]map[string]string, error) {

	rows, err := ms.Conn.Query(query, "%"+key+"%")

	if err != nil {
		return nil, fmt.Errorf("error executing query: %s", err)
	}

	var credList []map[string]string

	for rows.Next() {

		usrInfo := make(map[string]string, 3)
		var id, key, password string

		err := rows.Scan(&id, &key, &password)

		if err != nil {
			return nil, fmt.Errorf("error reading query data: %s", err)
		}

		usrInfo["id"], usrInfo["key"], usrInfo["password"] = id, key, password

		credList = append(credList, usrInfo)
	}

	return credList, nil

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
