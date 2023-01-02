package main

import (
  "net/http"
  "fmt"
  "strings"
  "log"
  "encoding/json"
  "github.com/gorilla/mux"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

func main() {
  // server := http.Server{
  //   Addr: ":8080",
  // }
  r := mux.NewRouter()
  r.HandleFunc("/", http.HandlerFunc(hoge))
  r.HandleFunc("/users", getUsers).Methods("GET")
  r.HandleFunc("/users", createUser).Methods("POST")
  r.HandleFunc("/users/{id}", getUser).Methods("GET")
  r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
  r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
  http.ListenAndServe(":8080", r)
}

type User struct {
  Id int `json:id`
  FirstName string `json:first_name`
  LastName string `json:last_name`
  Age int `json:age`
}

type UserRequest struct {
  FirstName string `json:"first_name"`
  LastName string `json:"last_name"`
  Age int `json:"age"`
}

func hoge(w http.ResponseWriter, _ *http.Request) {
//  db,err := sql.Open("mysql", "test_user:test_password@tcp(127.0.0.1:13306)/test_db")
//  if err != nil {
//    log.Fatal(err)
//  }
//  rows,err := db.Query("SELECT id, first_name, last_name, age FROM users")
//  if err != nil {
//    log.Fatal(err)
//  }
//  defer rows.Close()
//  for rows.Next() {
//    var user User
//    err := rows.Scan(&user.id, &user.first_name, &user.last_name, &user.age)
//    if err != nil {
//      log.Fatal(err)
//    }
//    log.Println(user.id, user.first_name, user.last_name, user.age)
//  }
//  defer db.Close()
//  fmt.Print("Hello world");
}

func getUsers(w http.ResponseWriter, r *http.Request) {
  fmt.Print("/users");
}

func getUser(w http.ResponseWriter, r *http.Request) {
  fmt.Print("/user");
  params := mux.Vars(r)
  fmt.Print(params);
  userId := strings.TrimPrefix(r.URL.Path, "/users/")
  fmt.Print(userId);
}

func createUser(w http.ResponseWriter, r *http.Request) {
  db,err := sql.Open("mysql", "test_user:test_password@tcp(127.0.0.1:13306)/test_db")
  if err != nil {
    log.Fatal(err)
  }
  stmt,err := db.Prepare("INSERT INTO users(first_name,last_name,age,created_at,updated_at) VALUES(?,?,?,NOW(),NOW())")
  if err != nil {
    log.Fatal(err)
  }

  body := make([]byte, r.ContentLength)
  r.Body.Read(body)
  var userRequest UserRequest
  json.Unmarshal(body, &userRequest)

  user := User{FirstName: userRequest.FirstName, LastName: userRequest.LastName, Age: userRequest.Age}

  res,err := stmt.Exec(user.FirstName, user.LastName, user.Age)
  if err != nil {
    log.Fatal(err)
  }
  lastId, err := res.LastInsertId()
  if err != nil {
    log.Fatal(err)
  }
  rowCnt, err := res.RowsAffected()
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
  defer stmt.Close()
}

func updateUser(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  userId := params["id"]

  db,err := sql.Open("mysql", "test_user:test_password@tcp(127.0.0.1:13306)/test_db")
  if err != nil {
    log.Fatal(err)
  }
  body := make([]byte, r.ContentLength)
  r.Body.Read(body)
  var userRequest UserRequest
  json.Unmarshal(body, &userRequest)

  user := User{FirstName: userRequest.FirstName, LastName: userRequest.LastName, Age: userRequest.Age}

  stmt,err := db.Prepare("UPDATE users SET first_name = ?, last_name = ?, age = ?, updated_at = NOW() WHERE users.id = ?")
  if err != nil {
    log.Fatal(err)
  }

  res,err := stmt.Exec(user.FirstName, user.LastName, user.Age, userId)
  if err != nil {
    log.Fatal(err)
  }
  lastId, err := res.LastInsertId()
  if err != nil {
    log.Fatal(err)
  }
  rowCnt, err := res.RowsAffected()
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
  defer stmt.Close()
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  userId := params["id"]

  db,err := sql.Open("mysql", "test_user:test_password@tcp(127.0.0.1:13306)/test_db")
  if err != nil {
    log.Fatal(err)
  }

  stmt,err := db.Prepare("DELETE FROM users WHERE id = ?")
  if err != nil {
    log.Fatal(err)
  }

  res,err := stmt.Exec(userId)
  if err != nil {
    log.Fatal(err)
  }
  lastId, err := res.LastInsertId()
  if err != nil {
    log.Fatal(err)
  }
  rowCnt, err := res.RowsAffected()
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
  defer stmt.Close()
}
