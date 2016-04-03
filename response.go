package main

import (
  "encoding/json"
  "net/http"
)

type ResponseError struct {
  Error string
}

func responseUnprocessable(w http.ResponseWriter) {
  response := ResponseError{"Unprocessable Entity"}
  js, err := json.Marshal(response)
  HandleError(err)
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(422)
  w.Write(js)
  w.Write([]byte("\n"))  // just for view convinience
}

func responseConflict(w http.ResponseWriter) {
  response := ResponseError{"Name already exists"}
  js, err := json.Marshal(response)
  HandleError(err)
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusConflict)
  w.Write(js)
  w.Write([]byte("\n"))
}

func responseCreated(w http.ResponseWriter) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)
}

func responseOK(w http.ResponseWriter) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
}

func responseNotfound(w http.ResponseWriter) {
  response := ResponseError{"Can't find the location"}
  js, err := json.Marshal(response)
  HandleError(err)
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusNotFound)
  w.Write(js)
  w.Write([]byte("\n"))
}