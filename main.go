package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/gorilla/mux"
)

type RedditSystem struct {
	system *actor.ActorSystem
}

// Initialize the global engine actor system
var engineActor *actor.PID

func main() {
	// Initialize ProtoActor system and the RedditEngine actor
	system := actor.NewActorSystem()
	engineActor = system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return &RedditEngine{
			users:      make(map[string]*User),
			subreddits: make(map[string]*Subreddit),
			posts:      make(map[string]*Post),
			comments:   make(map[string]*Comment),
		}
	}))

	rs := RedditSystem{system: system}

	// Initialize HTTP server with routes
	router := mux.NewRouter()
	InitializeRoutes(router, &rs)

	// Start the server
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
