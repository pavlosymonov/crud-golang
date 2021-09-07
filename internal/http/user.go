package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	userpb "main/proto"
	"net/http"
)

func (s *Server) CreateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var user userpb.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(response, NewResponse(StatusError, err.Error()), http.StatusInternalServerError)
		return
	}

	userReq := &userpb.CreateUserReq{User: &user}

	result, err := s.userClient.CreateUser(context.TODO(), userReq)
	if err != nil {
		http.Error(response, NewResponse(StatusError, err.Error()), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(response).Encode(result)
	if err != nil {
		http.Error(response, NewResponse(StatusError, err.Error()), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	req := &userpb.ListUsersReq{}
	stream, err := s.userClient.ListUsers(context.TODO(), req)
	if err != nil {
		http.Error(response, NewResponse(StatusError, err.Error()), http.StatusBadRequest)
		return
	}
	var users []*userpb.User
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(response, NewResponse(StatusError, err.Error()), http.StatusInternalServerError)
			return
		}
		users = append(users, res.GetUser())
	}

	err = json.NewEncoder(response).Encode(users)
	if err != nil {
		http.Error(response, NewResponse(StatusError, err.Error()), http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)["id"]

	req := &userpb.DeleteUserReq{
		Id: params,
	}

	_, err := s.userClient.DeleteUser(context.TODO(), req)
	if err != nil {
		http.Error(response, NewResponse(StatusError, err.Error()), http.StatusBadRequest)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
