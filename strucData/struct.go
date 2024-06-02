package strucData

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	User     string             `json:"user" `
	Username string             `json:"username" `
	UserId   primitive.ObjectID `json:"userId" bson:"_id""`
	UserIcon UserIcon           `json:"usericon" `
	Post     []Post             `json:"post" `
	Visible  bool               `json:"visible" `
}

type UserIcon struct {
	IconURL string `json:"iconurl"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PrivacySettings struct {
	Visibility   string   `json:"visibility"`
	AllowedUsers []string `json:"allowed_users"`
}

type Picture struct {
	URL string `json:"url"`
}

type Comment struct {
	CommentID   string    `json:"comment_id"`
	UserID      string    `json:"UserID"`
	Text        string    `json:"text"`
	UserName    string    `json:"UserName"`
	CommentDate time.Time `json:"comment_date"`
	Visible     bool      `json:"visible"`
}

type Like struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type Post struct {
	PostID          primitive.ObjectID `json:"post_id" bson:"_id"`
	PostDate        time.Time          `json:"post_date"`
	Description     string             `json:"description"`
	Tags            []string           `json:"tags"`
	Location        Location           `json:"location"`
	PrivacySettings PrivacySettings    `json:"privacy_settings"`
	Pictures        []Picture          `json:"pictures"`
	Comments        []Comment          `json:"comments"`
	Likes           []Like             `json:"likes"`
	Visible         bool               `json:"visible"`
}
