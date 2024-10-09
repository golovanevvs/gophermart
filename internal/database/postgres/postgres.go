package postgres

import "fmt"

type postgresStr struct {
}

func New() postgresStr {
	return postgresStr{}
}

func (db postgresStr) SaveToDB() string {
	return fmt.Sprintf("SaveToDB")
}
