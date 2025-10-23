package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// --- GLOBAL VARIABLES ---
var (
	DB                *mongo.Database
	EmailFrom         string
	EmailPassword     string
	GoogleOauthConfig *oauth2.Config
	Cloud             *cloudinary.Cloudinary
)

// --- INITIALIZER ---
func InitConfig() {
	// ✅ Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("⚠️ .env file not found in current directory, trying parent directory...")
		_ = godotenv.Load("../.env")
	}

	// ✅ Load Email credentials
	EmailFrom = os.Getenv("EMAIL_FROM")
	EmailPassword = os.Getenv("EMAIL_PASSWORD")
	fmt.Println("DEBUG: EMAIL_FROM =", EmailFrom)
	fmt.Println("DEBUG: EMAIL_PASSWORD =", EmailPassword)
	if EmailFrom == "" || EmailPassword == "" {
		log.Println("⚠️ Warning: EMAIL_FROM or EMAIL_PASSWORD not set in .env")
	}

	// ✅ Connect MongoDB
	connectMongoDB()

	// ✅ Initialize Google OAuth
	initGoogleOAuth()

	// ✅ Initialize Cloudinary
	initCloudinary()
}

// --- MONGO DATABASE CONNECTION ---
func connectMongoDB() {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Failed to create MongoDB client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB ping failed:", err)
	}

	DB = client.Database(dbName)
	fmt.Println("✅ Connected to MongoDB successfully!")
}

// --- GOOGLE OAUTH CONFIGURATION ---
func initGoogleOAuth() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	fmt.Println("DEBUG: GOOGLE_CLIENT_ID =", clientID)
	fmt.Println("DEBUG: GOOGLE_CLIENT_SECRET =", clientSecret)
	fmt.Println("DEBUG: GOOGLE_REDIRECT_URI =", redirectURI)

	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  redirectURI,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		log.Println("⚠️ Warning: Google OAuth variables not set in .env")
	} else {
		fmt.Println("✅ Google OAuth configured successfully!")
	}
}

// --- CLOUDINARY CONFIGURATION ---
func initCloudinary() {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatal("❌ Failed to connect to Cloudinary:", err)
	}

	Cloud = cld
	fmt.Println("✅ Connected to Cloudinary successfully!")
}
