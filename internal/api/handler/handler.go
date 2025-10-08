package handler

type Handler struct {
	service ServiceI
}

func New(s ServiceI) *Handler {
	return &Handler{service: s}
}
