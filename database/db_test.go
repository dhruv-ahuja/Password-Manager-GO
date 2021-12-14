package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestTableExists(t *testing.T) {

	db, mock, err := sqlmock.New()
	defer db.Close()

	if err != nil {
		t.Fatalf("error when attempting to make a mock connection: %s", err)
	}

	repo := NewDBRepo(db)

	mock.ExpectExec("SELECT 'public.info'::regclass").WillReturnError(nil)

	if err := repo.TableExists(); err == nil {
		t.Log("Executed fine")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected %v, but got %s", nil, err)
	}

}
