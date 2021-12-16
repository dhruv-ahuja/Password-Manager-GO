package Store

import "database/sql"

type Store interface {
	SaveCreds([]byte) error
	ViewCreds(string, []byte) error
	EditCreds(string, []byte) error
	DeleteCreds(string, []byte) error
	RetrieveCreds(string, []byte, *sql.DB) ([]map[string]string, error)
}
