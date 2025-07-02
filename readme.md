# RedditAPI

A simple Reddit-like API backend implemented in Go, powered by ProtoActor for actor-based concurrency and Gorilla Mux for HTTP routing.

## Overview

This project implements a Reddit-style backend API with the following features:

- User registration
- Subreddit creation and joining
- Post creation, upvoting, and downvoting
- Comment creation, upvoting, and downvoting
- Sending direct messages between users
- Fetching personalized user feeds

The backend uses **ProtoActor** (an actor model framework for Go) to manage internal state and concurrency, and **Gorilla Mux** for routing HTTP REST API endpoints.

## Project Structure

- `main.go` — Entry point, initializes ProtoActor system and HTTP server.
- `routers.go` — Defines HTTP API routes and handlers.
- `responses.go` — Utility functions for consistent JSON API responses.
- `go.mod` — Module dependencies.


## Tech Stack

- Go 1.23+
- [ProtoActor-Go](https://github.com/asynkron/protoactor-go) for actor concurrency model
- [Gorilla Mux](https://github.com/gorilla/mux) for HTTP routing
- JSON-based REST API

## Installation

```bash
git clone https://github.com/yourusername/RedditAPI.git
cd RedditAPI
```
## Run app

```bash
go run engine.go main.go
```

## API endpoints supported

| Method | Endpoint            | Description                | Request Body (JSON)                                                                              | Response                 |
| ------ | ------------------- | -------------------------- | ------------------------------------------------------------------------------------------------ | ------------------------ |
| POST   | `/register`         | Register a new user        | `{ "username": "user123" }`                                                                      | Success or error message |
| POST   | `/subreddit/create` | Create a new subreddit     | `{ "name": "golang", "description": "Go subreddit" }`                                            | Success or error message |
| POST   | `/subreddit/join`   | Join a subreddit           | `{ "username": "user123", "subreddit": "golang" }`                                               | Success or error message |
| POST   | `/post/create`      | Create a new post          | `{ "title": "Hello", "content": "World", "author": "user123", "subreddit": "golang" }`           | Success or error message |
| POST   | `/comment/create`   | Create a new comment       | `{ "content": "Nice post!", "author": "user123", "post_id": "postid", "parent_id": "optional" }` | Success or error message |
| POST   | `/post/upvote`      | Upvote a post or comment   | `{ "user_id": "user123", "media_type": "post/comment", "target_id": "postid/commentid" }`        | Success or error message |
| POST   | `/post/downvote`    | Downvote a post or comment | `{ "user_id": "user123", "media_type": "post/comment", "target_id": "postid/commentid" }`        | Success or error message |
| POST   | `/message/send`     | Send a direct message      | `{ "from": "user123", "to": "user456", "content": "Hello!" }`                                    | Success or error message |
| GET    | `/feed/{username}`  | Get personalized user feed | None                                                                                             | JSON feed data           |
