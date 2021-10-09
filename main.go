package main

import (
	"context"
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"encoding/json"
)

//
//
//
//
//
//
//
//
// Users
type Users struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

type UserController struct {
	collection *mongo.Collection
}

func createUserController(userCollection *mongo.Collection) *UserController {
	return &UserController{userCollection}
}

// Add user information to the user database
func (uc UserController) addUser(response http.ResponseWriter, request *http.Request, p httprouter.Params) {
	response.Header().Add("content-type", "application/json")
	var user Users
	json.NewDecoder(request.Body).Decode(&user)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := uc.collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}

func (uc UserController) getUsers(response http.ResponseWriter, request *http.Request, p httprouter.Params) {
	response.Header().Add("content-type", "application/json")
	var UserArray []Users
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := uc.collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message: "` + err.Error() + `"}"`))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user Users
		cursor.Decode(&user)
		UserArray = append(UserArray, user)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message: "` + err.Error() + `"}"`))
		return
	}
	json.NewEncoder(response).Encode(UserArray)
}

// Fetch user information based on the given user-id
func (uc UserController) getSingleUser(response http.ResponseWriter, request *http.Request, parameters httprouter.Params) {
	response.Header().Add("content-type", "application/json")
	id_string := parameters.ByName("id")

	id, _ := primitive.ObjectIDFromHex(id_string)
	var user Users
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := uc.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message: "` + err.Error() + `"}"`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

// Deletes all the users
// used while production and debugging
func (uc UserController) deleteUsers(response http.ResponseWriter, request *http.Request, p httprouter.Params) {
	response.Header().Add("content-type", "application/json")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	uc.collection.Drop(ctx)
}

//
//
//
//
//
//
//
// Posts
type Posts struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AccountId        primitive.ObjectID `json:"accountid,omitempty" bson:"accountid,omitempty"`
	Caption          string             `json:"caption,omitempty" bson:"caption,omitempty"`
	Image_URL        string             `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	Posted_Timestamp time.Time          `json:"posted_Timestamp,omitempty" bson:"Posted_Timestamp,omitempty"`
}

type PostsController struct {
	postscollection *mongo.Collection
}

func createPostController(collection *mongo.Collection) *PostsController {
	return &PostsController{collection}
}

// Add post information to the post database
func (pc PostsController) addPost(response http.ResponseWriter, request *http.Request, parameters httprouter.Params) {
	response.Header().Add("content-type", "application/json")
	var post Posts
	json.NewDecoder(request.Body).Decode(&post)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := pc.postscollection.InsertOne(ctx, post)
	json.NewEncoder(response).Encode(result)
}

func (pc PostsController) getPosts(response http.ResponseWriter, request *http.Request, parameters httprouter.Params) {
	response.Header().Add("content-type", "application/json")
	var postArray []Posts
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := pc.postscollection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message: "` + err.Error() + `"}"`))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var post Posts
		cursor.Decode(&post)
		postArray = append(postArray, post)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message: "` + err.Error() + `"}"`))
		return
	}
	json.NewEncoder(response).Encode(postArray)
}

// Fetch post based on the given post-id
func (pc PostsController) getSinglePost(response http.ResponseWriter, request *http.Request, parameters httprouter.Params) {
	response.Header().Add("content-type", "application/json")
	id_string := parameters.ByName("id")

	id, _ := primitive.ObjectIDFromHex(id_string)
	var post Posts
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := pc.postscollection.FindOne(ctx, bson.M{"_id": id}).Decode(&post)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message: "` + err.Error() + `"}"`))
		return
	}
	json.NewEncoder(response).Encode(post)
}

// lists all the posts posted by a particular user
func (pc PostsController) getPostByUser(response http.ResponseWriter, request *http.Request, parameters httprouter.Params) {
	response.Header().Add("content-type", "application/json")
	id_string := parameters.ByName("id")

	id, _ := primitive.ObjectIDFromHex(id_string)
	var postArray []Posts
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := pc.postscollection.Find(ctx, bson.M{"accountid": id})

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message: "` + err.Error() + `"}"`))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var post Posts
		cursor.Decode(&post)
		postArray = append(postArray, post)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message: "` + err.Error() + `"}"`))
		return
	}
	json.NewEncoder(response).Encode(postArray)
}

// Deletes all the posts
// used while production and debugging
func (pc PostsController) deleteposts(response http.ResponseWriter, request *http.Request, parameters httprouter.Params) {
	response.Header().Add("content-type", "application/json")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	pc.postscollection.Drop(ctx)
}

//
//
//
//
//
//
//
//
// main
// var client mongo.Client

func main() {
	fmt.Println("Start.")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	usercollection := client.Database("instadupe").Collection("users")
	userController := createUserController(usercollection)

	postcollection := client.Database("instadupe").Collection("posts")
	postController := createPostController(postcollection)

	r := httprouter.New()

	r.POST("/users", userController.addUser)
	r.GET("/users", userController.getUsers)
	r.GET("/users/:id", userController.getSingleUser)
	r.DELETE("/users", userController.deleteUsers)

	r.POST("/posts", postController.addPost)
	r.GET("/posts", postController.getPosts)
	r.GET("/posts/:id", postController.getSinglePost)
	r.DELETE("/posts", postController.deleteposts)

	r.GET("/posts/:id/users", postController.getPostByUser)

	log.Fatal(http.ListenAndServe(":8080", r))

}
