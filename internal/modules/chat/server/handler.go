package chat_server

type ServerHandler struct {
	ServerService *ServerService
}

func NewServerHandler(service *ServerService) *ServerHandler {
	return &ServerHandler{ServerService: service}
}
