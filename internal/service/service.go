package service

type serviceStr struct {
	db DatabaseInt
}

type DatabaseInt interface {
	SaveToDB() string
}

func New(db DatabaseInt) serviceStr {
	return serviceStr{
		db: db,
	}
}
