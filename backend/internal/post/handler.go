package post

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/studyplatform/backend/pkg/database"
	"github.com/studyplatform/backend/pkg/models"
)

// removeFromSlice removes an item from a slice and returns the new slice
func removeFromSlice(slice []string, item string) []string {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// ListPostsHandler returns all posts
func ListPostsHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		posts := mongoClient.GetCollection(database.CollectionNames.Posts)
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// For social feed, show all posts. For user posts, filter by author.
		// Check if this is for user's own posts (query parameter)
		filter := bson.M{}
		if c.Query("mine") == "true" {
			filter = bson.M{"author_id": userIDStr}
		}

		cursor, err := posts.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer cursor.Close(ctx)

		var postList []models.PostForResponse
		for cursor.Next(ctx) {
			var post models.Post
			if err := cursor.Decode(&post); err == nil {
				fmt.Printf("DEBUG: ListPosts - Processing post ID: %s, AuthorID: %s\n", post.ID.Hex(), post.AuthorID)
				fmt.Printf("DEBUG: ListPosts - Post raw data: %+v\n", post)

				postResponse := post.ToResponse(userIDStr)

				// Fetch author username
				var user models.User
				// Convert authorID string to ObjectID for lookup
				authorObjID, err := primitive.ObjectIDFromHex(post.AuthorID)
				if err != nil {
					fmt.Printf("Error converting author ID %s to ObjectID: %v\n", post.AuthorID, err)
					postResponse.AuthorName = "Unknown User"
				} else {
					fmt.Printf("DEBUG: ListPosts - Looking up user with ObjectID: %v\n", authorObjID)
					err := users.FindOne(ctx, bson.M{"_id": authorObjID}).Decode(&user)
					if err == nil {
						postResponse.AuthorName = user.Username
						postResponse.AuthorUniqueID = user.UniqueID
						fmt.Printf("DEBUG: Found username '%s' for author ID %s\n", user.Username, post.AuthorID)
					} else {
						fmt.Printf("DEBUG: User not found for author ID %s: %v\n", post.AuthorID, err)
						postResponse.AuthorName = "Unknown User"
					}
				}

				// Populate comment author names and unique IDs
				for i, comment := range postResponse.Comments {
					var commentUser models.User
					commentAuthorObjID, err := primitive.ObjectIDFromHex(comment.AuthorID)
					if err == nil {
						err := users.FindOne(ctx, bson.M{"_id": commentAuthorObjID}).Decode(&commentUser)
						if err == nil {
							postResponse.Comments[i].AuthorName = commentUser.Username
							postResponse.Comments[i].AuthorUniqueID = commentUser.UniqueID
						} else {
							postResponse.Comments[i].AuthorName = "Unknown User"
						}
					} else {
						postResponse.Comments[i].AuthorName = "Unknown User"
					}
				}

				postList = append(postList, postResponse)
			}
		}

		c.JSON(http.StatusOK, gin.H{"posts": postList})
	}
}

// CreatePostHandler creates a new post
func CreatePostHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		fmt.Printf("DEBUG: CreatePost - User ID from context: %s\n", userIDStr)

		var req models.CreatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Printf("DEBUG: CreatePost - JSON binding failed: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		fmt.Printf("DEBUG: CreatePost - Request content: %s\n", req.Content)

		posts := mongoClient.GetCollection(database.CollectionNames.Posts)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		postMap := bson.M{
			"content":    req.Content,
			"author_id":  userIDStr,
			"likes":      []string{},
			"comments":   []models.Comment{},
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}

		fmt.Printf("DEBUG: CreatePost - About to insert post: %+v\n", postMap)

		res, err := posts.InsertOne(ctx, postMap)
		if err != nil {
			fmt.Printf("DEBUG: CreatePost - Insert failed: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
			return
		}

		fmt.Printf("DEBUG: CreatePost - Insert successful, ID: %v\n", res.InsertedID)

		var post models.Post
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			post.ID = oid
		}
		err = posts.FindOne(ctx, bson.M{"_id": post.ID}).Decode(&post)
		if err != nil {
			fmt.Printf("DEBUG: CreatePost - Failed to fetch created post: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created post"})
			return
		}

		fmt.Printf("DEBUG: CreatePost - Fetched post: %+v\n", post)

		// Fetch author username
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		var user models.User
		err = users.FindOne(ctx, bson.M{"unique_id": userIDStr}).Decode(&user)
		if err != nil {
			fmt.Printf("DEBUG: CreatePost - Failed to fetch user info: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
			return
		}

		fmt.Printf("DEBUG: CreatePost - Found user: %s (ID: %s)\n", user.Username, user.ID.Hex())

		postResponse := post.ToResponse(userIDStr)
		postResponse.AuthorName = user.Username
		postResponse.AuthorUniqueID = user.UniqueID

		fmt.Printf("DEBUG: Created post with author username '%s' for user ID %s\n", user.Username, userIDStr)

		c.JSON(http.StatusCreated, gin.H{"post": postResponse})
	}
}

// LikePostHandler likes or unlikes a post
func LikePostHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		var req models.LikePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		postID, err := primitive.ObjectIDFromHex(req.PostID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		posts := mongoClient.GetCollection(database.CollectionNames.Posts)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if user already liked the post
		var post models.Post
		err = posts.FindOne(ctx, bson.M{"_id": postID}).Decode(&post)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		// Check if user already liked
		alreadyLiked := false
		for _, likeID := range post.Likes {
			if likeID == userIDStr {
				alreadyLiked = true
				break
			}
		}

		var update bson.M
		if alreadyLiked {
			// Unlike
			update = bson.M{"$pull": bson.M{"likes": userIDStr}}
		} else {
			// Like
			update = bson.M{"$push": bson.M{"likes": userIDStr}}
		}

		_, err = posts.UpdateOne(ctx, bson.M{"_id": postID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
			return
		}

		// Create notification for post author when post is liked (not when unliked)
		if !alreadyLiked && post.AuthorID != userIDStr {
			notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)

			// Get the liker's username for the notification
			users := mongoClient.GetCollection(database.CollectionNames.Users)
			var liker models.User
			if err := users.FindOne(ctx, bson.M{"unique_id": userIDStr}).Decode(&liker); err == nil {
				postLikeNotification := models.CreatePostLikeNotification(post.AuthorID, liker.ID.Hex(), liker.Username, postID.Hex(), post.Content)
				_, err = notifications.InsertOne(ctx, postLikeNotification)
				if err != nil {
					fmt.Printf("DEBUG: LikePost - Warning: Failed to create notification: %v\n", err)
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
	}
}

// CreateCommentHandler creates a new comment on a post
func CreateCommentHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		postID := c.Param("postId")
		if postID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID required"})
			return
		}

		var req models.CreateCommentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		postObjID, err := primitive.ObjectIDFromHex(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		posts := mongoClient.GetCollection(database.CollectionNames.Posts)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		comment := models.Comment{
			ID:        primitive.NewObjectID(),
			PostID:    postID,
			Content:   req.Content,
			AuthorID:  userIDStr,
			Likes:     []string{},
			CreatedAt: time.Now(),
		}

		update := bson.M{"$push": bson.M{"comments": comment}}
		_, err = posts.UpdateOne(ctx, bson.M{"_id": postObjID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
			return
		}

		// Fetch the user to get the username for the response
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		var user models.User
		err = users.FindOne(ctx, bson.M{"unique_id": userIDStr}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
			return
		}

		// Create notification for post author when comment is added (if commenter is not the post author)
		// Get the post to find the author ID
		var postData models.Post
		if err := posts.FindOne(ctx, bson.M{"_id": postObjID}).Decode(&postData); err == nil {
			if postData.AuthorID != userIDStr {
				notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
				postCommentNotification := models.CreatePostCommentNotification(postData.AuthorID, userIDStr, user.Username, postID)
				_, err = notifications.InsertOne(ctx, postCommentNotification)
				if err != nil {
					fmt.Printf("DEBUG: CreateComment - Warning: Failed to create notification: %v\n", err)
				}
			}
		}

		// Create comment response with author info
		commentResponse := models.CommentForResponse{
			ID:             comment.ID.Hex(),
			PostID:         postID,
			Content:        comment.Content,
			AuthorID:       userIDStr,
			AuthorUniqueID: user.UniqueID,
			AuthorName:     user.Username,
			LikesCount:     0,
			IsLiked:        false,
			CreatedAt:      comment.CreatedAt,
		}

		c.JSON(http.StatusCreated, gin.H{"comment": commentResponse})
	}
}

// LikeCommentHandler handles liking/unliking a comment
func LikeCommentHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		commentID := c.Param("commentId")
		if commentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Comment ID required"})
			return
		}

		fmt.Printf("DEBUG: LikeComment - Comment ID: %s\n", commentID)

		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		fmt.Printf("DEBUG: LikeComment - User ID: %s\n", userIDStr)

		posts := mongoClient.GetCollection(database.CollectionNames.Posts)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Convert commentID string to ObjectID for the query
		commentObjID, err := primitive.ObjectIDFromHex(commentID)
		if err != nil {
			fmt.Printf("DEBUG: LikeComment - Invalid comment ID format: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID format"})
			return
		}

		fmt.Printf("DEBUG: LikeComment - Converted comment ObjectID: %s\n", commentObjID.Hex())

		// Find the post containing the comment
		var post models.Post
		err = posts.FindOne(ctx, bson.M{"comments._id": commentObjID}).Decode(&post)
		if err == mongo.ErrNoDocuments {
			fmt.Printf("DEBUG: LikeComment - Comment not found in database\n")
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		} else if err != nil {
			fmt.Printf("DEBUG: LikeComment - Database error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		fmt.Printf("DEBUG: LikeComment - Found post with ID: %s\n", post.ID.Hex())
		fmt.Printf("DEBUG: LikeComment - Post has %d comments\n", len(post.Comments))

		// Find the comment and toggle like
		for i, comment := range post.Comments {
			fmt.Printf("DEBUG: LikeComment - Checking comment %d: ID=%s, Hex=%s, comparing with %s\n", i, comment.ID, comment.ID.Hex(), commentID)
			if comment.ID.Hex() == commentID {
				fmt.Printf("DEBUG: LikeComment - Found matching comment at index %d\n", i)
				// Toggle like status
				liked := false
				for _, like := range comment.Likes {
					if like == userIDStr {
						liked = true
						break
					}
				}

				if liked {
					// Unlike
					post.Comments[i].Likes = removeFromSlice(post.Comments[i].Likes, userIDStr)
				} else {
					// Like
					post.Comments[i].Likes = append(post.Comments[i].Likes, userIDStr)

					// Create notification for comment author when comment is liked (if liker is not the comment author)
					if comment.AuthorID != userIDStr {
						notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)

						// Get the liker's username for the notification
						users := mongoClient.GetCollection(database.CollectionNames.Users)
						var liker models.User
						if err := users.FindOne(ctx, bson.M{"unique_id": userIDStr}).Decode(&liker); err == nil {
							commentLikeNotification := models.CreateCommentLikeNotification(comment.AuthorID, liker.ID.Hex(), liker.Username, commentID, comment.Content)
							_, err = notifications.InsertOne(ctx, commentLikeNotification)
							if err != nil {
								fmt.Printf("DEBUG: LikeComment - Warning: Failed to create notification: %v\n", err)
							}
						}
					}
				}

				// Update the post
				_, err := posts.UpdateOne(ctx, bson.M{"_id": post.ID}, bson.M{"$set": bson.M{"comments": post.Comments}})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"message": "Comment like updated"})
				return
			}
		}

		fmt.Printf("DEBUG: LikeComment - Comment not found after searching through all comments\n")
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
	}
}

// DeletePostHandler handles deleting a post
func DeletePostHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("id")
		if postID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID required"})
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		objID, err := primitive.ObjectIDFromHex(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		posts := mongoClient.GetCollection(database.CollectionNames.Posts)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if post exists and user is the author
		var post models.Post
		err = posts.FindOne(ctx, bson.M{"_id": objID}).Decode(&post)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Debug logging
		fmt.Printf("DEBUG: Delete post - Post AuthorID: %s, User ID from context: %s\n", post.AuthorID, userIDStr)
		fmt.Printf("DEBUG: Delete post - Post AuthorID type: %T, User ID type: %T\n", post.AuthorID, userIDStr)
		fmt.Printf("DEBUG: Delete post - Post AuthorID length: %d, User ID length: %d\n", len(post.AuthorID), len(userIDStr))
		fmt.Printf("DEBUG: Delete post - Post AuthorID bytes: %v, User ID bytes: %v\n", []byte(post.AuthorID), []byte(userIDStr))

		// Check if user is the author - normalize both IDs to strings
		postAuthorID := post.AuthorID
		if postAuthorID != userIDStr {
			fmt.Printf("DEBUG: Delete post - Permission denied. Post author: %s, Current user: %s\n", postAuthorID, userIDStr)
			fmt.Printf("DEBUG: Delete post - String comparison failed. Checking if they're equal after trimming...\n")
			if strings.TrimSpace(postAuthorID) == strings.TrimSpace(userIDStr) {
				fmt.Printf("DEBUG: Delete post - IDs are equal after trimming whitespace\n")
				postAuthorID = strings.TrimSpace(postAuthorID)
			} else {
				c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
				return
			}
		}

		// Delete the post
		_, err = posts.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
	}
}
