package redis

import (
  "github.com/fzzy/radix/redis"
  "time"
)

const (
  ADDRESS = "localhost:6379"
)

var (
  client, dbErr = redis.DialTimeout("tcp", ADDRESS, time.Duration(10)*time.Second)
)

func StartRedisClient() error {
  return dbErr
} 

func CloseRedisClient() {
  client.Close()
}

// check if the location already stored in database
func LocationExists(name string) bool {
  r, err := client.Cmd("EXISTS", "city:" + name).Bool()
  HandleError(err)
  return r
}

// create a location entity in database
func CreateLocation(name string) {
  addErr := client.Cmd("SADD", "cities", name).Err
  HandleError(addErr)
  setErr := client.Cmd("HMSET", "city:" + name, "name", name).Err
  HandleError(setErr)
}

// delete a location from database
func DeleteLocation(name string) {
  delErr := client.Cmd("DEL", "city:" + name).Err
  HandleError(delErr)
  remErr := client.Cmd("SREM", "cities", name).Err
  HandleError(remErr)
}

// list all locations in database
func ListLocations() []string {
  ls, err := client.Cmd("SMEMBERS", "cities").List()
  HandleError(err)
  return ls
}

// check if the location weather info expired (1 hour)
func WeatherExpired(name string) bool {
  now := time.Now()
  results, err := client.Cmd("HMGET", "city:" + name, "updated_at").List()  // weird the HMGET return a list result
  if err != nil {
    return true
  }
  updatedAtStr := results[0]
  updatedAt, _ := time.Parse(time.RFC3339, updatedAtStr)
  expireAt := updatedAt.Add(1 * time.Hour)
  return now.After(expireAt)
}

// store weather info in database
func StoreWeather(name string, weather string) {
  err := client.Cmd("HMSET", "city:" + name, "weather", weather, "updated_at", time.Now().Format(time.RFC3339)).Err
  HandleError(err)
}

// get weather info from database
func ReadWeather(name string) string {
  results, err := client.Cmd("HMGET", "city:" + name, "weather").List()
  if err != nil {
    return ""
  }
  return results[0]
}

func HandleError(err error) {
  if err != nil {
    panic(err)
  }
}