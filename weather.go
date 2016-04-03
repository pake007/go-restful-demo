package main

import (
  "fmt"
  "io"
  "io/ioutil"
  "encoding/json"
  "net/http"
)

const (
  KEY = "f87dfd3af38ed44f157296b7150caacc"
)

type City struct {
  Name string `json:"name"`
}

// {"id":500,"main":"Rain","description":"light rain","icon":"10d"}
type Weather struct {
  Id int32 `json:"id"`
  Main string `json:"main"`
  Description string `json:"description"`
  Icon string `json:"icon"`
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
func getLocationHandler(w http.ResponseWriter, name string) {
  fmt.Println("Get weather of " + name)
  expired := weatherExpired(name)
  existingWeather := readWeather(name)
  // if no weather info in database or weather info expired (> 1 hour), request openweathermap api
  if expired || lend(existingWeather) == 0 {
    resp, _ := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + KEY + "&q=" + name)
    body, _ := ioutil.ReadAll(resp.Body)
    var f interface{}
    err := json.Unmarshal(body, &f)
    HandleError(err)
    m := f.(map[string]interface{})
    for k, v := range m {
      if k == "weather" {
        for _, vv := range v.([]interface{}) {
          fmt.Println(vv.(map[string]interface{}))
        }
        break
      }  
    }
    storeWeather(name, weather)
    // responseWeather(weather)   
  } else {
    fmt.Println("use existing weather")
    // responseWeather(existingWeather)
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
      getLocationHandler(w, name)
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