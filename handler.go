package main

import (
  "fmt"
  "io"
  "io/ioutil"
  "encoding/json"
  "net/http"
  "strings"
  "github.com/pake007/go-restful-demo/weatherapi"
  "github.com/pake007/go-restful-demo/redis"
)

// index
func indexHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    locations := redis.ListLocations()
    err := json.NewEncoder(w).Encode(locations)
    HandleError(err)
  } else {
    fmt.Println("Not a valid action!")
    return
  }
}

// create
func createHandler(w http.ResponseWriter, r *http.Request) {
  var city City
  if r.Method == "POST" {
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    HandleError(err)
    defer r.Body.Close()

    // create the city
    if err := json.Unmarshal(body, &city); err != nil || len(city.Name) == 0 {
      responseUnprocessable(w)
    } else {
      exists := redis.LocationExists(strings.ToLower(city.Name))
      if exists {
        responseConflict(w)
      } else {
        redis.CreateLocation(strings.ToLower(city.Name))
        responseCreated(w)
      }
    }
  } else {
    fmt.Println("Not a valid action!")
    return
  }
}

// show location weather info
func getWeatherHandler(w http.ResponseWriter, r *http.Request) {
  name := getLoctionName(r)
  if len(name) == 0 {
    return
  }
  if !redis.LocationExists(name) {
    responseNotfound(w)
    return
  }
  expired := redis.WeatherExpired(name)
  existingWeather := redis.ReadWeather(name)
  // if no weather info in database or weather info expired (> 1 hour), request openweathermap api
  if expired || len(existingWeather) == 0 {
    weatherResp := weatherapi.RequestWeather(name)
    redis.StoreWeather(name, string(weatherResp))
    responseWeather(w, weatherResp)   
  } else {
    fmt.Println("use existing weather of " + name)
    responseWeather(w, []byte(existingWeather))
  }
}

// delete
func deleteHandler(w http.ResponseWriter, r *http.Request) {
  name := getLoctionName(r)
  if len(name) == 0 {
    return
  }
  exists := redis.LocationExists(name)
  if exists {
    redis.DeleteLocation(name)
    responseOK(w)
  } else {
    responseNotfound(w)
  }
}

// GET or DELETE dispatcher
func locationHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
      getWeatherHandler(w, r)
    case "DELETE":
      deleteHandler(w, r)
    default:
      fmt.Println("Not a valid action!")
      return
  }
}

// ------------ helper method for get or delete, parse the location name from url ------------
func getLoctionName(r *http.Request) string {
  name := r.URL.Path[len("/location/"):]
  if len(name) == 0 {
    fmt.Println("Not a valid name!")
    return ""
  }
  return strings.ToLower(name)
}