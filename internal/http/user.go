package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main/internal/models/user"
	"net/http"
)

func (s *Server) CreateUser(response http.ResponseWriter, request *http.Request) {
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

func (s *Server) GetUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	users, err := s.userService.List(context.TODO())
	if err != nil {
		http.Error(response, NewResponse(StatusError, err.Error()), http.StatusBadRequest)
		return
	}

	json.NewEncoder(response).Encode(users)
}

func (s *Server) DeleteUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)["id"]

	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		http.Error(response, NewResponse(StatusError, err.Error()), http.StatusInternalServerError)
	}

	err = s.userService.Delete(context.TODO(), _id)
	if err != nil {
		http.Error(response, NewResponse(StatusError, err.Error()), http.StatusBadRequest)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
