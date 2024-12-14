package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

// Define message types
type RegisterUser struct {
	Username string
}

type CreateSubreddit struct {
	Name        string
	Description string
}

type JoinSubreddit struct {
	Username  string
	Subreddit string
}

type LeaveSubreddit struct {
	Username  string
	Subreddit string
}

type CreatePost struct {
	Title     string
	Content   string
	Author    string
	Subreddit string
}

type CreateComment struct {
	Content  string
	Author   string
	PostID   string // Post ID should be handled correctly in a real scenario
	ParentID string // Optional: ID of the parent comment
}

type Upvote struct {
	UserID    string // Username in a real scenario
	MediaType string //type of media
	TargetID  string // Can be a Post or Comment ID
}

type Downvote struct {
	UserID    string // Username in a real scenario
	MediaType string //type of media
	TargetID  string // Can be a Post or Comment ID
}

type SendDirectMessage struct {
	From    string
	To      string
	Content string
}

type GetUserFeed struct {
	Username string
	writer   http.ResponseWriter
}

// User represents a Reddit user.
type User struct {
	ID       string
	Username string
	Karma    int
	Inbox    []string // List of direct messages received.
}

// Post represents a Reddit post.
type Post struct {
	ID        string
	Title     string
	Content   string
	Author    *User      // Reference to the user who created the post.
	Subreddit *Subreddit // Reference to the subreddit where the post was made.
	Upvotes   int
	Downvotes int
	CreatedAt time.Time
}

// Comment represents a comment on a post.
type Comment struct {
	ID        string
	Content   string
	Author    *User   // Reference to the user who created the comment.
	Post      *Post   // Reference to the post where the comment was made.
	ParentID  *string // Optional: ID of the parent comment for hierarchical comments.
	Upvotes   int
	Downvotes int
	CreatedAt time.Time
}

// Subreddit represents a subreddit.
type Subreddit struct {
	ID          string           // Unique identifier for the subreddit.
	Name        string           // Name of the subreddit.
	Description string           // Description of the subreddit.
	Members     map[string]*User // Map of usernames to User objects who are members of the subreddit.
}

// RedditEngine is the main actor for the Reddit clone engine.
type RedditEngine struct {
	users      map[string]*User      // Map of username to User details.
	subreddits map[string]*Subreddit // Map of subreddit name to Subreddit details.
	// usernames  map[string]string     // Map of username to user ID for quick lookup.
	posts    map[string]*Post    // Map of post ID to Post details.
	comments map[string]*Comment // Map of comment ID to Comment details.
}

// to get user feed json object
type PostInfo struct {
	SubredditName string `json:"subreddit_name"`
	Title         string `json:"title"`
	AuthorName    string `json:"author_name"`
}

// Receive handles incoming messages for the RedditEngine actor.
func (re *RedditEngine) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *RegisterUser:
		re.registerUser(msg.Username, context)
	case *CreateSubreddit:
		re.createSubreddit(msg.Name, msg.Description, context)
	case *JoinSubreddit:
		re.joinSubreddit(msg.Username, msg.Subreddit, context)
	case *LeaveSubreddit:
		re.leaveSubreddit(msg.Username, msg.Subreddit, context)
	case *CreatePost:
		re.createPost(msg.Title, msg.Content, msg.Author, msg.Subreddit, context)
	case *CreateComment:
		re.createComment(msg.Content, msg.Author, msg.PostID, msg.ParentID, context)
	case *Upvote:
		re.upvote(msg.UserID, msg.MediaType, msg.TargetID, context)
	case *Downvote:
		re.downvote(msg.UserID, msg.MediaType, msg.TargetID, context)
	case *SendDirectMessage:
		re.sendDirectMessage(msg.From, msg.To, msg.Content, context)
	case *GetUserFeed:
		re.getUserFeed(msg.Username, msg.writer, context)
	default:
		fmt.Println("Engine Initiallised")
	}
}

func (re *RedditEngine) registerUser(username string, context actor.Context) {
	if _, exists := re.users[username]; exists {
		fmt.Printf("Username %s already taken\n", username)
		context.Respond(false)
		return
	}
	user := &User{ID: username, Username: username}
	re.users[username] = user
	fmt.Printf("Registered new user: %s\n", username)
	context.Respond(true)
}

func (re *RedditEngine) createSubreddit(name, description string, context actor.Context) {
	if _, exists := re.subreddits[name]; exists {
		fmt.Printf("Subreddit %s already exists\n", name)
		context.Respond(false)
		return
	}
	subreddit := &Subreddit{
		ID:          name,
		Name:        name,
		Description: description,
		Members:     make(map[string]*User),
	}
	re.subreddits[name] = subreddit
	fmt.Printf("Created new subreddit: %s\n", name)
	context.Respond(true)
}

func (re *RedditEngine) joinSubreddit(username, subredditName string, context actor.Context) {
	user, userExists := re.users[username]
	if !userExists {
		fmt.Printf("No such user with username %s\n", username)
		context.Respond(301)
		return
	}

	subreddit, subExists := re.subreddits[subredditName]
	if !subExists {
		fmt.Printf("No such subreddit with name %s\n", subredditName)
		context.Respond(302)
		return
	}

	subreddit.Members[username] = user
	fmt.Printf("User %s joined subreddit %s\n", username, subreddit.Name)
	context.Respond(200)
}

func (re *RedditEngine) leaveSubreddit(username, subredditName string, context actor.Context) {
	subreddit, exists := re.subreddits[subredditName]
	if !exists {
		fmt.Printf("No such subreddit with name %s\n", subredditName)
		context.Respond(301)
		return
	}
	delete(subreddit.Members, username)
	fmt.Printf("User %s left subreddit %s\n", username, subreddit.Name)
	context.Respond(200)
}

func (re *RedditEngine) createPost(title, content, authorName, subredditName string, context actor.Context) {
	user, userExists := re.users[authorName]
	if !userExists {
		fmt.Printf("No such user with username %s\n", authorName)
		context.Respond(301)
		return
	}

	subreddit, subExists := re.subreddits[subredditName]
	if !subExists {
		fmt.Printf("No such subreddit with name %s\n", subredditName)
		context.Respond(302)
		return
	}

	postId := fmt.Sprintf("%s_post_%d", authorName, len(re.posts)+1)
	post := &Post{
		ID:        postId,
		Title:     title,
		Content:   content,
		Author:    user,
		Subreddit: subreddit,
		CreatedAt: time.Now(),
	}
	re.posts[postId] = post
	fmt.Printf("Created new post in subreddit %s by user %s with id %s\n", subreddit.Name, authorName, postId)
	context.Respond(200)
}

func (re *RedditEngine) createComment(content, authorName, postId, parentId string, context actor.Context) {
	user, userExists := re.users[authorName]
	if !userExists {
		fmt.Printf("No such user with username %s\n", authorName)
		context.Respond(301)
		return
	}

	post, postExists := re.posts[postId]
	if !postExists {
		fmt.Printf("No such post with ID %s\n", postId)
		context.Respond(302)
		return
	}

	commentId := fmt.Sprintf("%s_comment_%d", authorName, len(re.comments)+1)
	comment := &Comment{
		ID:        commentId,
		Content:   content,
		Author:    user,
		Post:      post,
		CreatedAt: time.Now(),
	}
	if parentId != "" { // If it's a reply to another comment
		comment.ParentID = &parentId
	}
	re.comments[commentId] = comment
	fmt.Printf("Created new comment on post %s by user %s with id %s\n", postId, authorName, commentId)
	context.Respond(200)
}

func (re *RedditEngine) upvote(userId, mediaType string, targetId string, context actor.Context) {
	if mediaType == "Post" {
		if post, exists := re.posts[targetId]; exists { // Upvoting a post
			post.Upvotes++
			fmt.Printf("User %s upvoted post %s\n", userId, targetId)
			context.Respond(201)
			return
		} else {
			fmt.Printf("No such Post with ID %s for upvote\n", targetId)
			context.Respond(301)
		}
	} else {
		if mediaType == "Comment" {
			if comment, exists := re.comments[targetId]; exists { // Upvoting a comment
				comment.Upvotes++
				fmt.Printf("User %s upvoted comment %s\n", userId, targetId)
				context.Respond(202)
				return
			}
		} else {
			fmt.Printf("No such comment with ID %s for upvote\n", targetId)
			context.Respond(302)
		}
	}
}

func (re *RedditEngine) downvote(userId, mediaType string, targetId string, context actor.Context) {
	if mediaType == "Post" {
		if post, exists := re.posts[targetId]; exists { // Upvoting a post
			post.Downvotes++
			fmt.Printf("User %s downvoted post %s\n", userId, targetId)
			context.Respond(201)
			return
		} else {
			fmt.Printf("No such Post with ID %s for downvote\n", targetId)
			context.Respond(301)
		}
	} else {
		if mediaType == "Comment" {
			if comment, exists := re.comments[targetId]; exists { // Upvoting a comment
				comment.Downvotes++
				fmt.Printf("User %s downvoted comment %s\n", userId, targetId)
				context.Respond(202)
				return
			}
		} else {
			fmt.Printf("No such comment with ID %s for upvote\n", targetId)
			context.Respond(302)
		}
	}
}

func (re *RedditEngine) sendDirectMessage(fromUsername, toUsername, content string, context actor.Context) {
	toUser, exists := re.users[toUsername]
	if !exists {
		fmt.Printf("No such recipient with username %s\n", toUsername)
		context.Respond(301)
		return
	}
	fromUser, exists := re.users[fromUsername]

	if !exists {
		fmt.Printf("No such sender with username %s\n", fromUsername)
		context.Respond(302)
		return
	}

	message := fmt.Sprintf("From:% s -% s", fromUsername, content)
	toUser.Inbox = append(toUser.Inbox, message)
	fmt.Printf("Direct message sent from%s to%s \n", fromUsername, toUsername)
	context.Respond(200)
}

func (re *RedditEngine) getUserFeed(username string, w http.ResponseWriter, context actor.Context) {
	user, exists := re.users[username]
	if !exists {
		fmt.Printf("No such user with username%s \n", username)
		jsonData := map[string]interface{}{
			"error": "User doesn't exist",
		}
		JSONSuccess(w, jsonData)
		context.Respond(200)
		return
	}

	fmt.Printf("Feed fetched for%s :\n", username)

	// var posts []PostInfo

	// for _, post := range re.posts {
	// 	if _, member := post.Subreddit.Members[user.ID]; member {
	// 		posts = append(posts, PostInfo{
	// 			SubredditName: post.Subreddit.Name,
	// 			Title:         post.Title,
	// 			AuthorName:    post.Author.Username,
	// 		})
	// 	}
	// }
	result := make(map[string]interface{})
	posts := []map[string]string{}

	for _, post := range re.posts {
		if _, member := post.Subreddit.Members[user.ID]; member {
			postInfo := map[string]string{
				"subreddit": post.Subreddit.Name,
				"title":     post.Title,
				"author":    post.Author.Username,
			}
			posts = append(posts, postInfo)
		}
	}

	result["posts"] = posts

	// jsonData, err := json.Marshal(posts)

	// if err != nil {
	// 	fmt.Println("Error marshalling to JSON:", err)
	// 	context.Respond(302)
	// 	return
	// }
	JSONFeed(w, result)
	context.Respond(200)
}
