package controler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"socialMedia/database"
	"socialMedia/strucData"
	"time"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "page")

func NewUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var newUser strucData.User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newUser.Visible = true

		// Insert the new post document
		_, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return success response
		c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully!"})

	}
}
func NewPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the user ID from the URL parameter
		userIDParam := c.Param("userid")
		userID, err := primitive.ObjectIDFromHex(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		fmt.Printf("UserID: %s\n", userID.Hex())

		// Bind the incoming JSON to the Post struct
		var newPost strucData.Post
		if err := c.ShouldBindJSON(&newPost); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate a new ObjectID for the post
		newPost.PostID = primitive.NewObjectID()
		newPost.PostDate = time.Now() // Set the post date to the current time
		newPost.Visible = true

		fmt.Printf("New Post: %+v\n", newPost)

		// Define the update to add the new post to the user's posts array
		update := bson.M{
			"$push": bson.M{"post": newPost},
		}

		// Update the user's document by adding the new post
		result, err := userCollection.UpdateOne(ctx, bson.M{"_id": userID}, update, options.Update().SetUpsert(true))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("Update Result: %+v\n", result)

		// Return success response with the newly created post's ID
		c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully!", "post_id": newPost.PostID.Hex()})
	}
}

func NewCommentPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the user ID from the URL parameter
		sender := c.Param("sender")
		_, err := primitive.ObjectIDFromHex(sender)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender ID"})
			return
		}

		// Bind the incoming JSON to the Comment struct
		var newComment struct {
			OwnerPost string `json:"ownerPost" binding:"required"`
			PostID    string `json:"postid" binding:"required"`
			strucData.Comment
		}

		if err := c.ShouldBindJSON(&newComment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert post ID string to ObjectID
		postID, err := primitive.ObjectIDFromHex(newComment.PostID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		ownerPost, err := primitive.ObjectIDFromHex(newComment.OwnerPost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		// Generate a new ObjectID for the comment
		newComment.Comment.CommentID = primitive.NewObjectID().Hex()
		newComment.Comment.UserID = sender
		newComment.Comment.CommentDate = time.Now() // Set the comment date to the current time
		newComment.Comment.Visible = true

		// Define the update to add the new comment to the post's comments array
		update := bson.M{
			"$push": bson.M{"post.$[p].comments": newComment.Comment},
		}

		// Filter to update only the specific post in the user's document
		filter := bson.M{"_id": ownerPost}
		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{bson.M{"p._id": postID}},
		}

		// Update the user's document by adding the new comment to the post's comments array
		result, err := userCollection.UpdateOne(ctx, filter, update, &options.UpdateOptions{ArrayFilters: &arrayFilters})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return success response with the newly created comment's ID
		c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully!", "comment_id": newComment.Comment.CommentID,
			"Text": newComment.Comment.Text, "sender": sender, "result": result})
	}
}

func NewLikePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the user ID from the URL parameter
		userIDParam := c.Param("userid")
		_, err := primitive.ObjectIDFromHex(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Bind the incoming JSON to the Comment struct
		var newLike struct {
			UserID string `json:"userid" binding:"required"`
			PostID string `json:"postid" binding:"required"`
			strucData.Like
		}

		if err := c.ShouldBindJSON(&newLike); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := primitive.ObjectIDFromHex(newLike.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		// Convert post ID string to ObjectID
		postID, err := primitive.ObjectIDFromHex(newLike.PostID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		// Generate a new ObjectID for the comment
		newLike.Like.UserID = userIDParam

		// Define the update to add the new comment to the post's comments array
		update := bson.M{
			"$push": bson.M{"post.$[p].likes": newLike.Like},
		}

		// Filter to update only the specific post in the user's document
		filter := bson.M{"_id": userID}
		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{bson.M{"p._id": postID}},
		}

		// Update the user's document by adding the new comment to the post's comments array
		result, err := userCollection.UpdateOne(ctx, filter, update, &options.UpdateOptions{ArrayFilters: &arrayFilters})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return success response with the newly created comment's ID
		c.JSON(http.StatusCreated, gin.H{"message": "Like created successfully!", "Like": newLike.Like, "Userid": userIDParam, "result": result, "update": update, "filter": filter})
	}
}

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the user ID from the URL parameter
		userIDParam := c.Param("userid")
		userID, err := primitive.ObjectIDFromHex(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Bind the incoming JSON to get the post ID
		var postRequest struct {
			PostID string `json:"postid" binding:"required"`
		}

		if err := c.ShouldBindJSON(&postRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert post ID string to ObjectID
		postID, err := primitive.ObjectIDFromHex(postRequest.PostID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		// Define the update to set the post's Visible field to false
		update := bson.M{
			"$set": bson.M{"post.$[p].visible": false},
		}

		// Filter to update only the specific post in the user's document
		filter := bson.M{"_id": userID}
		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{bson.M{"p._id": postID}},
		}

		// Update the user's document by setting the post's Visible field to false
		result, err := userCollection.UpdateOne(ctx, filter, update, &options.UpdateOptions{ArrayFilters: &arrayFilters})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "Post visibility set to false", "result": result})
	}
}

func DeleteComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the user ID from the URL parameter
		userIDParam := c.Param("userid")
		userID, err := primitive.ObjectIDFromHex(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Bind the incoming JSON to get the post ID and comment ID
		var request struct {
			PostID    string `json:"postid" binding:"required"`
			CommentID string `json:"commentid" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert post ID string to ObjectID
		postID, err := primitive.ObjectIDFromHex(request.PostID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		// Define the update to set the comment's visible field to false
		update := bson.M{
			"$set": bson.M{"post.$[p].comments.$[c].visible": false},
		}

		// Filter to update only the specific post and comment in the user's document
		filter := bson.M{"_id": userID}
		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"p._id": postID},
				bson.M{"c.commentid": request.CommentID},
			},
		}

		// Debugging logs
		//c.JSON(http.StatusOK, gin.H{
		//	"message":      "Parameters received",
		//	"userID":       userID.Hex(),
		//	"postID":       postID.Hex(),
		//	"commentID":    request.CommentID,
		//	"filter":       filter,
		//	"update":       update,
		//	"arrayFilters": arrayFilters.Filters,
		//})

		// Update the user's document by setting the comment's visible field to false
		result, err := userCollection.UpdateOne(ctx, filter, update, &options.UpdateOptions{ArrayFilters: &arrayFilters})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "Comment visibility set to false", "result": result})
	}
}

func GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		userIDParam := c.Param("userid")
		userID, err := primitive.ObjectIDFromHex(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		filter := bson.M{"_id": userID}

		cursor, err := userCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var results []interface{}
		for cursor.Next(ctx) {
			var result bson.M
			if err := cursor.Decode(&result); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			}
			results = append(results, result)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

func FindUserByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get the username from the URL parameter
		username := c.Param("userName")
		//c.JSON(http.StatusInternalServerError, gin.H{"username": username})

		// Define the filter to search for users with the given username
		filter := bson.M{"username": username}

		// Find the users matching the filter
		cursor, err := userCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var userIDs []string
		for cursor.Next(ctx) {
			var user strucData.User
			if err := cursor.Decode(&user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			userIDs = append(userIDs, user.UserId.Hex())
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the list of user IDs
		c.JSON(http.StatusOK, gin.H{"user_ids": userIDs})
	}
}
