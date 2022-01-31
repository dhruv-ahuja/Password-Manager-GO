package store

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupForTests() (*sql.DB, sqlmock.Sqlmock, *mockStore) {
	// generates a mock *sql.DB connection to use for testing
	db, mock, err := sqlmock.New()

	if err != nil {
		log.Fatalf("An error encountered when generating new mock connection: %s", err)
	}

	ms := newMockStore(db)

	return db, mock, ms
}

func TestSaveCreds(t *testing.T) {

	db, mock, ms := setupForTests()
	defer db.Close()

	mock.ExpectExec("INSERT INTO info").
		WithArgs("reddit", "test123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := ms.SaveCreds([]byte("ok")); err != nil {
		t.Errorf("error when running func: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}

}

func TestRetrieveCreds(t *testing.T) {

	query := "SELECT * FROM info WHERE key ILIKE ? ORDER BY id ASC"

	db, mock, ms := setupForTests()
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "key", "encrypted_pw"}).AddRow("1", "reddit", "test123").AddRow("2", "Reddit", "testing")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM info WHERE key ILIKE ? ORDER BY id ASC`)).WithArgs("%red%").WillReturnRows(mockRows)

	var want []map[string]string

	want = append(want, map[string]string{"id": "1", "key": "reddit", "password": "test123"}, map[string]string{"id": "2", "key": "Reddit", "password": "testing"})

	credList, err := ms.RetrieveCreds(query, "red", []byte("test"))

	if err != nil {
		t.Log(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations unmet: %s", err)
	}

	// check if credList matches want
	assert.Equal(t, want, credList, "Function returns incorrect values!")

}

// TestViewCreds will fail. Since ViewCreds method uses RetrieveCreds as it's
// base, we can test a failing condition here instead of copy-pasting
func TestViewCreds(t *testing.T) {

	db, mock, ms := setupForTests()
	defer db.Close()

	// mockRows := sqlmock.NewRows([]string{"id", "key", "encrypted_pw"}).AddRow("1", "reddit", "test123").AddRow("2", "Reddit", "testing")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM info WHERE key ILIKE ? ORDER BY id ASC`)).WithArgs("%spotify%").WillReturnError(fmt.Errorf("no accounts found"))

	_ = ms.ViewCreds("spotify", []byte("test"))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectatations: %s", err)
	}

	// the query returns the error that we had specified, means our test passed

}

func TestEditCreds(t *testing.T) {

	db, mock, ms := setupForTests()
	defer db.Close()

	mockRows := mock.NewRows([]string{"id", "key", "encrypted_pw"}).
		AddRow("10", "spotify", "randomPassword")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM info`)).WithArgs("%spot%").WillReturnRows(mockRows)

	usrInfo := map[string]string{
		"id": "10", "key": "spotify", "encrypted_pw": "randomPassword",
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE info`)).
		WithArgs(usrInfo["key"], usrInfo["encrypted_pw"], usrInfo["id"]).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := ms.EditCreds("spot", []byte("test")); err != nil {
		t.Log(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations unmet: %s", err)
	}

}

func TestDeleteCreds(t *testing.T) {

	db, mock, ms := setupForTests()
	defer db.Close()

	mockRows := mock.NewRows([]string{"id", "key", "encrypted_pw"}).
		AddRow("10", "spotify", "randomPassword")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM info`)).WithArgs("%spot%").WillReturnRows(mockRows)

	// sample data for testing
	usrInfo := map[string]string{
		"id": "10", "key": "spotify", "encrypted_pw": "randomPassword",
	}

	mock.ExpectExec(`DELETE FROM info`).WithArgs(usrInfo["id"]).WillReturnResult(sqlmock.NewResult(1, 1))

	if err := ms.DeleteCreds("spot", []byte("test")); err != nil {
		t.Log(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
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
		return nil, err
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

func (ms *mockStore) ViewCreds(key string, encryptionKey []byte) error {
	// this is basically the same as RetrieveCreds since it just uses
	// that to get its work done. RetrieveCreds serves as the common
	// helper function for all DB queries.
	query := "SELECT * FROM info WHERE key ILIKE ? ORDER BY id ASC"

	credList, err := ms.RetrieveCreds(query, "spotify", encryptionKey)

	if err != nil || len(credList) == 0 {
		return err
	}

	return nil

}

func (ms *mockStore) EditCreds(key string, encryptionKey []byte) error {

	query := "SELECT * FROM info WHERE key ILIKE ? ORDER BY id ASC"

	_, err := ms.Conn.Query(query, "%"+key+"%")

	if err != nil {
		return fmt.Errorf("error executing query: %s", err)
	}

	usrInfo := map[string]string{
		"id": "10", "key": "spotify", "encrypted_pw": "randomPassword",
	}

	editQuery := "UPDATE info SET key=?, encrypted_pw=? WHERE id=?"

	_, err = ms.Conn.Exec(editQuery, usrInfo["key"], usrInfo["encrypted_pw"], usrInfo["id"])

	if err != nil {
		return fmt.Errorf("error executing update query: %s", editQuery)
	}

	return nil
}

func (ms *mockStore) DeleteCreds(key string, encryptionKey []byte) error {

	query := "SELECT * FROM info WHERE key ILIKE ? ORDER BY id ASC"

	_, err := ms.Conn.Query(query, "%"+key+"%")

	if err != nil {
		return fmt.Errorf("error executing query: %s", err)
	}

	usrInfo := map[string]string{
		"id": "10", "key": "spotify", "encrypted_pw": "randomPassword",
	}

	deleteQuery := "DELETE FROM info WHERE id=?"

	_, err = ms.Conn.Exec(deleteQuery, usrInfo["id"])

	if err != nil {
		return fmt.Errorf("error executing deletion query: %s", err)
	}

	return nil
}
