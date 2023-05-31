package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hakonslie/twitchfriendmaker/auth"
	"github.com/hakonslie/twitchfriendmaker/config"
	"github.com/hakonslie/twitchfriendmaker/logger"
	"github.com/hakonslie/twitchfriendmaker/session"
)

type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
}

func exchangeCodeForAccessToken(cfg *config.Config, code string) (*TokenResponse, error) {
	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", map[string][]string{
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {cfg.RedirectURL},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, err
	}
	return &tokenResponse, nil
}

func Redirect(sessions *session.Session, authStorage *auth.AuthStorage, cfg *config.Config, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.MustGet("sessionID").(session.SessionID)
		code := c.Query("code")
		if code != "" {
			state := c.Query("state")

			if !sessions.CompareState(sessionID, state) {
				log.Error("Something wrong with scope or state")
				fmt.Println("state or scope doesnt match")
				return
			}

			accessToken, err := exchangeCodeForAccessToken(cfg, code)
			if err != nil {
				log.Error("Couldn't get access token from code")
				fmt.Println(err)
				return
			}
			authStorage.AddToken(sessionID, auth.TokenData{AccessToken: accessToken.AccessToken, RefreshToken: accessToken.RefreshToken})

			c.Redirect(http.StatusSeeOther, "/")
			return
		}

		errorCode := c.Query("error")
		errorDescription := c.Query("error_description")
		if errorCode != "" {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"title":            "Authorization Denied",
				"errorCode":        errorCode,
				"errorDescription": errorDescription,
			})
			return
		} else {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"title": "Authorization Rejected",
			})
			return
		}
	}
}
