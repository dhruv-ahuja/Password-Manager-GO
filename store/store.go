package store

type Store interface {
	SaveCreds([]byte) error

	RetrieveCreds(string, string, []byte) ([]map[string]string, error)

	ViewCreds(string, []byte) error

	EditCreds(string, []byte) error

	DeleteCreds(string, []byte) error
}
