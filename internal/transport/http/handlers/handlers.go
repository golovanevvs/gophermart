package handlers

type HandlersStr struct {
	service ServicesInt
}

type ServicesInt interface {
	RegisterUser() string
}

func New(sv ServicesInt) HandlersStr {
	return HandlersStr{
		service: sv,
	}
}
