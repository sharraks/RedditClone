package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Initialize routes
func InitializeRoutes(router *mux.Router, rs *RedditSystem) {
	router.HandleFunc("/register", RegisterUserHandler(rs)).Methods("POST")
	router.HandleFunc("/subreddit/create", CreateSubredditHandler(rs)).Methods("POST")
	router.HandleFunc("/subreddit/join", JoinSubredditHandler(rs)).Methods("POST")
	router.HandleFunc("/post/create", CreatePostHandler(rs)).Methods("POST")
	router.HandleFunc("/comment/create", CreateCommentHandler(rs)).Methods("POST")
	router.HandleFunc("/post/upvote", UpvoteHandler(rs)).Methods("POST")
	router.HandleFunc("/post/downvote", DownvoteHandler(rs)).Methods("POST")
	router.HandleFunc("/message/send", SendDirectMessageHandler(rs)).Methods("POST")
	router.HandleFunc("/feed/{username}", GetUserFeedHandler(rs)).Methods("GET")
}

// Handle user registration
func RegisterUserHandler(rs *RedditSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Username string `json:"username"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			JSONError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Create the RegisterUser message and send it to the engine actor
		result := rs.system.Root.RequestFuture(engineActor, &RegisterUser{Username: request.Username}, 1*time.Second)

		if resp, err := result.Result(); resp == true && err == nil {
			// Respond with success message
			JSONSuccess(w, "User registered successfully")
		} else {
			JSONError(w, 200, "Username already taken")
		}
	}
}

// Handle subreddit creation
func CreateSubredditHandler(rs *RedditSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			JSONError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		result := rs.system.Root.RequestFuture(engineActor, &CreateSubreddit{Name: request.Name, Description: request.Description}, 1*time.Second)

		if resp, err := result.Result(); resp == true && err == nil {
			// Respond with success message
			JSONSuccess(w, "Subreddit created successfully")
		} else {
			JSONError(w, 200, "Subreddit already exists")
		}
	}
}

// Handle joining a subreddit
func JoinSubredditHandler(rs *RedditSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Username  string `json:"username"`
			Subreddit string `json:"subreddit"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			JSONError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Send the JoinSubreddit message to the engine actor
		result := rs.system.Root.RequestFuture(engineActor, &JoinSubreddit{Username: request.Username, Subreddit: request.Subreddit}, 1*time.Second)

		if resp, err := result.Result(); resp == 200 && err == nil {
			// Respond with success message
			JSONSuccess(w, "Subreddit joined successfully")
		} else {
			if resp == 301 {
				JSONError(w, 403, "No such username")
			} else if resp == 302 {
				JSONError(w, 403, "No such subreddit")
			}
		}
	}
}

// Handle post creation
func CreatePostHandler(rs *RedditSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Title     string `json:"title"`
			Content   string `json:"content"`
			Author    string `json:"author"`
			Subreddit string `json:"subreddit"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			JSONError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Send the createPost message to the engine actor
		result := rs.system.Root.RequestFuture(engineActor, &CreatePost{
			Title:     request.Title,
			Content:   request.Content,
			Author:    request.Author,
			Subreddit: request.Subreddit,
		}, 1*time.Second)

		if resp, err := result.Result(); resp == 200 && err == nil {
			// Respond with success message
			JSONSuccess(w, "Post created successfully")
		} else {
			if resp == 301 {
				JSONError(w, 403, "No such username")
			} else if resp == 302 {
				JSONError(w, 403, "No such subreddit")
			}
		}
	}
}

// Handle comment creation
func CreateCommentHandler(rs *RedditSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Content  string `json:"content"`
			Author   string `json:"author"`
			PostID   string `json:"post_id"`
			ParentID string `json:"parent_id,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			JSONError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Send the CreateComment message to the engine actor
		result := rs.system.Root.RequestFuture(engineActor, &CreateComment{
			Content:  request.Content,
			Author:   request.Author,
			PostID:   request.PostID,
			ParentID: request.ParentID,
		}, 1*time.Second)

		if resp, err := result.Result(); resp == 200 && err == nil {
			// Respond with success message
			JSONSuccess(w, "Comment created successfully")
		} else {
			if resp == 301 {
				JSONError(w, 403, "No such username")
			} else if resp == 302 {
				JSONError(w, 403, "No such post")
			}
		}
	}
}

// Handle upvoting a post or comment
func UpvoteHandler(rs *RedditSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			UserID    string `json:"user_id"`
			MediaType string `json:"media_type"`
			TargetID  string `json:"target_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			JSONError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		/// Send the CreateComment message to the engine actor
		result := rs.system.Root.RequestFuture(engineActor, &Upvote{UserID: request.UserID, MediaType: request.MediaType, TargetID: request.TargetID}, 1*time.Second)

		if resp, err := result.Result(); resp == 201 && err == nil {
			// Respond with success message
			JSONSuccess(w, "Post Upvoted successfully")
		} else if resp == 202 {
			// Respond with success message
			JSONSuccess(w, "Comment Upvoted successfully")
		} else {
			if resp == 301 {
				JSONError(w, 403, "No such post")
			} else if resp == 302 {
				JSONError(w, 403, "No such comment")
			}
		}
	}
}

// Handle downvoting a post or comment
func DownvoteHandler(rs *RedditSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			UserID    string `json:"user_id"`
			MediaType string `json:"media_type"`
			TargetID  string `json:"target_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			JSONError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Send the Downvote message to the engine actor
		result := rs.system.Root.RequestFuture(engineActor, &Downvote{UserID: request.UserID, MediaType: request.MediaType, TargetID: request.TargetID}, 1*time.Second)

		if resp, err := result.Result(); resp == 201 && err == nil {
			// Respond with success message
			JSONSuccess(w, "Post Downvoted successfully")
		} else if resp == 202 {
			// Respond with success message
			JSONSuccess(w, "Comment Downvoted successfully")
		} else {
			if resp == 301 {
				JSONError(w, 403, "No such post")
			} else if resp == 302 {
				JSONError(w, 403, "No such comment")
			}
		}
	}
}

// Handle sending direct messages
func SendDirectMessageHandler(rs *RedditSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			From    string `json:"from"`
			To      string `json:"to"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			JSONError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Send the SendDirectMessage message to the engine actor
		result := rs.system.Root.RequestFuture(engineActor, &SendDirectMessage{
			From:    request.From,
			To:      request.To,
			Content: request.Content,
		}, 1*time.Second)

		if resp, err := result.Result(); resp == 200 && err == nil {
			// Respond with success message
			JSONSuccess(w, "DM sent successfully")
		} else {
			if resp == 301 {
				JSONError(w, 403, "Sender doesn't exist")
			} else if resp == 302 {
				JSONError(w, 403, "Receiver doesn't exist")
			}
		}
	}
}

// Handle getting a user's feed
func GetUserFeedHandler(rs *RedditSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		// Send the GetUserFeed message to the engine actor
		result := rs.system.Root.RequestFuture(engineActor, &GetUserFeed{Username: username, writer: w}, 1*time.Second)

		if resp, err := result.Result(); resp == 200 && err == nil {
			return
		}
	}
}
