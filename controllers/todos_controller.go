package controllers

import (
    "encoding/json"
    "net/http"
    "io"
    "io/ioutil"
    "to-go/models"
    "to-go/storage"
    "github.com/gorilla/mux"
    "strings"
)

type TodoRequest struct {
    Todo models.Todo
    AuthToken string
}

type TodoFailResponse struct {
    Message    string `json:"message"`
}

func newTodoFailResponse() TodoFailResponse {
    response := TodoFailResponse{}
    response.Message = "Could not find a todo item with the given ID."
    return response
}

func TodosIndex(w http.ResponseWriter, r *http.Request) {
    db := storage.GetActiveDB()
    todos := models.Todos{}
    authToken := strings.Join(r.Header["X-Auth-Token"], "")
    db.Joins("join users on user_id = users.id").Where("users.auth_token = ?", authToken).Find(&todos)

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(todos); err != nil {
        panic(err)
    }
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
    db := storage.GetActiveDB()
    vars := mux.Vars(r)
    todoId := vars["todoId"]
    todo := models.Todo{}
    db.Find(&todo, todoId)

    if todo.ID == 0 {
        response := newTodoFailResponse()
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(404)
        if err := json.NewEncoder(w).Encode(&response); err != nil {
            panic(err)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(todo); err != nil {
        panic(err)
    }
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {
    db := storage.GetActiveDB()
    todo :=  models.Todo{}
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &todo); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    db.Create(&todo)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(&todo); err != nil {
        panic(err)
    }
}

func TodoUpdate(w http.ResponseWriter, r *http.Request) {
    db := storage.GetActiveDB()
    vars := mux.Vars(r)
    todoId := vars["todoId"]
    todo := models.Todo{}
    db.Find(&todo, todoId)

    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &todo); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    db.Save(&todo)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(&todo); err != nil {
        panic(err)
    }
}

func TodoDelete(w http.ResponseWriter, r *http.Request) {
    db := storage.GetActiveDB()
    vars := mux.Vars(r)
    todoId := vars["todoId"]
    todo := models.Todo{}
    db.Find(&todo, todoId)
    db.Delete(&todo, todoId)

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(&todo); err != nil {
        panic(err)
    }

}

