package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go_project.com/m/auth"
	"go_project.com/m/config"
	"go_project.com/m/follows"
	"go_project.com/m/session"
)

type User struct {
	ID string `json:"id"`
}

type UserResponse struct {
	Data []User `json:"data"`
}

type Follow struct {
	BroadcasterName string `json:"broadcaster_name"`
}

type FollowResponse struct {
	Total int      `json:"total"`
	Data  []Follow `json:"data"`
}

func GetFollows(follow follows.FollowStorage, auth auth.AuthStorage, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID := ctx.MustGet("sessionID").(session.SessionID)
		token := auth.GetToken(sessionID)

		// Construct the URL with query parameters
		url := "https://api.twitch.tv/helix/users"

		// Create a new HTTP request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set the required headers
		req.Header.Set("Client-ID", cfg.ClientID)
		req.Header.Set("Authorization", "Bearer "+token)

		// Create an HTTP client with timeout
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Getting User")
			fmt.Println(token)
			fmt.Println("Request failed with status code:", resp.StatusCode)
			return
		}

		body, err := io.ReadAll(resp.Body)
		var response UserResponse
		err = json.Unmarshal([]byte(body), &response)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}
		// Process the response body
		// Here, you can parse the JSON response and handle the data accordingly

		fmt.Println("Request successful")
		userID := response.Data[0].ID
		fmt.Printf("UserID: %s \n", userID)

		// Construct the URL with query parameters
		url = "https://api.twitch.tv/helix/channels/followed?user_id=" + userID + "&first=100"

		// Create a new HTTP request
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set the required headers
		req.Header.Set("Client-ID", cfg.ClientID)
		req.Header.Set("Authorization", "Bearer "+token)

		// Send the request
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Getting followed")
			fmt.Println("Request failed with status code:", resp.StatusCode)
			return
		}

		// Process the response body
		// Here, you can parse the JSON response and handle the data accordingly
		body, err = io.ReadAll(resp.Body)
		var res FollowResponse
		err = json.Unmarshal([]byte(body), &res)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		var channels []follows.Channel
		// Access the total and broadcaster_name
		fmt.Println("Total:", res.Total)
		for _, follow := range res.Data {
			fmt.Println("Broadcaster Name:", follow.BroadcasterName)
			channels = append(channels, follows.Channel(follow.BroadcasterName))
		}

		follow.AddFollows(sessionID, channels)

		fmt.Println("Request successful")

		ctx.Redirect(http.StatusSeeOther, "/")
		return

	}
}
