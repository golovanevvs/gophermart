package database

type database interface {
	SaveToDB()
	GetFromDB()
}
