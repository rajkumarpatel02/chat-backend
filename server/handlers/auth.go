package handlers

import (
	"encoding/json"
	"net/http"

	"go-grpc-chat/config"
	"go-grpc-chat/models"
	"go-grpc-chat/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// ✅ Check if user already exists
	var existing models.User
	err := config.UserCollection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existing)
	if err != mongo.ErrNoDocuments && err != nil {
		http.Error(w, "Error checking existing user", http.StatusInternalServerError)
		return
	}
	if err == nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	// ✅ Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// ✅ Insert user
	_, err = config.UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var found models.User
	err := config.UserCollection.FindOne(context.TODO(), bson.M{"username": creds.Username}).Decode(&found)
	if err == mongo.ErrNoDocuments {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Error finding user", http.StatusInternalServerError)
		return
	}

	if !utils.CheckPasswordHash(creds.Password, found.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(found.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
