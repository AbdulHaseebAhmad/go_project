package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AbdulHaseebAhmad/go_project/internal/config"
	student "github.com/AbdulHaseebAhmad/go_project/internal/httpHandler"
	"github.com/AbdulHaseebAhmad/go_project/internal/storage/postgress"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup

	storage, storageerr := postgress.New(cfg) // the wrapper package we created, return the instance and error
	if storageerr != nil {
		log.Fatal(storageerr)
		return
	}
	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("Path", cfg.Storage_path))
	// setup router
	router := http.NewServeMux() // server mux is server multiplexer that routes http request to its specific handler functions

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome To the App"))
	})

	// getting handler func from the student package, same like in js we have that call back function this route creates a student
	router.HandleFunc("POST /api/students/create", student.New(storage))

	// getting handler func from the student package, same like in js we have that call back function. this route is to get the student by id
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	// getting handler func from the student package, same like in js we have that call back function. this route is to get the student list
	router.HandleFunc("GET /api/students/", student.GetStudentList(storage))

	// getting handler func from the student package, same like in js we have that call back function. this route is to get the student list
	router.HandleFunc("GET /api/students/delete/{id}", student.DeleteStudent(storage))
	// setup server

	server := http.Server{
		Addr:    cfg.HttpServer.Address,
		Handler: router,
	}
	fmt.Println("Server starting at:", cfg.HttpServer.Address)

	// Before adding gracefull shutdown
	// err := server.ListenAndServe()

	// if err != nil {
	// 	log.Fatal("Failed To Start Server")
	// } else {
	// }

	// After adding gracefull shutdown

	done := make(chan os.Signal, 1) // make a chanel name done, has a buffer size of 1 and has type of os.signal

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // signal.notify rallies signal to c, the first argument is
	// where we want (which channel)to direct the incoming signal, theen the
	// remaining arguments are which type of signals are we interested in

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("Failed To Start Server")
		} else {
		}
	}()

	<-done                                   // this will block, the code below until a signal is recieved
	slog.Info("Shutting down gracefully...") //this will only run once a signal is recieved into done channel.

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// below we can kill other process and db connectins etc
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Info("failed to shhutdown", slog.String("erro", err.Error()))
	}

	slog.Info("Server Shut Down Successfully")
}
