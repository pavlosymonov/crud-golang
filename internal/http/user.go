package http

import (
	"context"
	"encoding/json"
	"main/internal/models/user"
	"net/http"
)

func (s *Server) CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var user user.User
	json.NewDecoder(request.Body).Decode(&user)

	result, err := s.userService.Create(context.TODO(), user)
	if err != nil {
		http.Error(response, NewResponse(StatusError, err.Error()), http.StatusBadRequest)
		return
	}

	json.NewEncoder(response).Encode(result)
}

func (s *Server) GetUsersEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	users, err := s.userService.List(context.TODO())
	if err != nil {
		http.Error(response, NewResponse(StatusError, err.Error()), http.StatusBadRequest)
		return
	}

	json.NewEncoder(response).Encode(users)
}
