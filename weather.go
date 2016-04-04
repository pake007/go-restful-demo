package main

import (
  "fmt"
  "io"
  "io/ioutil"
  "encoding/json"
  "net/http"
)

const (
  ApiAddress = "http://api.openweathermap.org/data/2.5/weather?APPID="
  ApiKey = "f87dfd3af38ed44f157296b7150caacc"
)

type City struct {
  Name string `json:"name"`
}


// index
func indexHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    locations := listLocations()
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
      exists := locationExists(city.Name)
      if exists {
        responseConflict(w)
      } else {
        createLocation(city.Name)
        responseCreated(w)
      }
    }
  } else {
    fmt.Println("Not a valid action!")
    return
  }
}

// show location weather info
func getLocationWeatherHandler(w http.ResponseWriter, name string) {
  if !locationExists(name) {
    responseNotfound(w)
    return
  }
  expired := weatherExpired(name)
  existingWeather := readWeather(name)
  // if no weather info in database or weather info expired (> 1 hour), request openweathermap api
  if expired || len(existingWeather) == 0 {
    weatherResp := requestWeatherAPI(name)
    storeWeather(name, string(weatherResp))
    responseWeather(w, weatherResp)   
  } else {
    fmt.Println("use existing weather of " + name)
    responseWeather(w, []byte(existingWeather))
  }
}

// delete
func deleteLocationHandler(w http.ResponseWriter, name string) {
  fmt.Println("Delete city " + name)
  exists := locationExists(name)
  if exists {
    deleteLocation(name)
    responseOK(w)
  } else {
    responseNotfound(w)
  }
}

// GET or DELETE
func locationHandler(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Path[len("/location/"):]
  if len(name) == 0 {
    fmt.Println("Not a valid name!")
    return
  }
  switch r.Method {
    case "GET":
      getLocationWeatherHandler(w, name)
    case "DELETE":
      deleteLocationHandler(w, name)
    default:
      fmt.Println("Not a valid action!")
      return
  }
}

func HandleError(err error) {
  if err != nil {
    panic(err)
  }
}

func main() {  
  if dbErr != nil {
    fmt.Println("Can't connect to redis:", dbErr)
    return
  }
  defer client.Close()

  http.HandleFunc("/locations", indexHandler)
  http.HandleFunc("/location", createHandler)
  http.HandleFunc("/location/", locationHandler)
  http.ListenAndServe(":8080", nil)
}