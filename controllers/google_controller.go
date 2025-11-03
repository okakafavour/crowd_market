package controllers

import (
	"CROWD_MARKET/config"
	"CROWD_MARKET/services"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

// Redirect-based Google login (for backend OAuth flow)
func GoogleLogin(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// OAuth callback for backend login
func GoogleCallBack(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code not found"})
		return
	}

	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to exchange token"})
		return
	}

	client := config.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode user info"})
		return
	}

	email, ok1 := userInfo["email"].(string)
	name, ok2 := userInfo["name"].(string)
	if !ok1 || !ok2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user data from Google"})
		return
	}

	user, _ := services.FindUserByEmail(email)
	if user == nil {
		services.CreateGoogleUser(name, email)
	}

	tokenString, _ := services.GenerateJWT(email)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"name":  name,
			"email": email,
		},
	})
}

// ✅ Web login using Flutter/React with ID token
func GoogleLoginWeb(c *gin.Context) {
	var body struct {
		IdToken string `json:"idToken"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// ✅ Validate the Google ID token
	payload, err := idtoken.Validate(context.Background(), body.IdToken, "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid ID token"})
		return
	}

	// ✅ Allow both backend and frontend client IDs
	validAudiences := []string{
		os.Getenv("GOOGLE_CLIENT_ID"), // Backend Client ID (from .env)
		"208804397507-0bu3dp1neogk2a5nl75s85cvtthvcc81.apps.googleusercontent.com", // Frontend Client ID
	}

	aud, _ := payload.Claims["aud"].(string)
	valid := false
	for _, v := range validAudiences {
		if aud == v {
			valid = true
			break
		}
	}

	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized audience"})
		return
	}

	// ✅ Extract user details from token
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	// ✅ Create or find user
	user, _ := services.FindUserByEmail(email)
	if user == nil {
		services.CreateGoogleUser(name, email)
	}

	// ✅ Generate your app's JWT token
	token, err := services.GenerateJWT(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate JWT"})
		return
	}

	// ✅ Respond with token and user info
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"email": email,
		"name":  name,
	})
}
