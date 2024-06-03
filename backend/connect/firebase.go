package connect

import (
	"backend/global"
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/appcheck"
)

// var (
// 	appCheckClient appcheck.Client
// )

// Firebase Admin configuration
func InitializeFirebase() (*firebase.App, error) {

	// Firebase configuration
	// https://firebase.google.com/docs/admin/setup
	conf := &firebase.Config{
		DatabaseURL:   global.Conf.FirebaseDatabaseURL,
		ProjectID:     global.Conf.FirebaseProjectID,
		StorageBucket: global.Conf.FirebaseStorageBucket,
	}

	// Initialize the Firebase Admin SDK
	// https://firebase.google.com/docs/admin/setup#initialize-sdk
	firebaseApp, err := firebase.NewApp(context.Background(), conf)
	if err != nil {
		return nil, err
	}

	return firebaseApp, nil
}

// Initialize Firebase App Check client
func InitializeFirebaseAppCheck() (appcheck.Client, error) {
	ctx := context.Background()

	firebaseApp, err := InitializeFirebase()
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	appCheckClient, err := firebaseApp.AppCheck(ctx)
	if err != nil {
		log.Fatalf("error initializing app check: %v\n", err)
	}

	log.Println("App Check initialized")

	return *appCheckClient, nil
}

// Initialize Firebase Authorization client
func InitializeFirebaseAuth() (*firebase.App, error) {
	ctx := context.Background()

	client, err := InitializeFirebase()
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	_, err = client.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing auth: %v\n", err)
	}

	log.Println("Auth initialized")

	return client, nil
}
