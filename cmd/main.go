package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/SmilingAli3n/crud_refactored/pkg/auth"
	"github.com/SmilingAli3n/crud_refactored/pkg/cache"
	"github.com/SmilingAli3n/crud_refactored/pkg/repos"
	"github.com/SmilingAli3n/crud_refactored/pkg/response"

	"github.com/gorilla/mux"
)

var c = cache.New(time.Minute)
var once sync.Once

func init() {
	go once.Do(func() {
		c.Init()
	})
}

// @title Ticket API
// @version 1.0
// @description This is a sample Ticket CRUD

// @Router /tickets/{id}

// @host tickets.swagger.io
// @BasePath /tickets
// @Param id path int true "int, id of ticket"
// @Router /tickets/{id} [post]
// @Router /tickets/{id} [put]
// @Router /tickets/{id} [get]
// @Router /tickets/{id} [delete]
func main() {
    r := mux.NewRouter()
	r.HandleFunc("/tickets/{id}", ticketsHandler)
	http.Handle("/", r)
	fmt.Println("Server started")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ticketsHandler(w http.ResponseWriter, req *http.Request) {
	resp := response.New()
	defer resp.Send(w)
	if !auth.Authorized(req) {
		resp.StatusCode = http.StatusUnauthorized
		return
	}

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Print(err)
	}
	if req.Method == http.MethodPost {
		repos.CreateTicket(req, resp)
	} else if req.Method == http.MethodGet {
        if id != 0 {
    		repos.GetAllTickets(req, resp, c)
        } else {
		    repos.GetTicketById(req, resp, c, int64(id))
        }
	} else if req.Method == http.MethodPut {
		repos.UpdateTicket(req, resp, int64(id))
	} else if req.Method == http.MethodDelete {
		repos.DeleteTicketById(req, resp, int64(id))
	} else {
		resp.Status = fmt.Sprintf("Method %v is not supported", req.Method)
		resp.StatusCode = http.StatusMethodNotAllowed
		return
	}
}
