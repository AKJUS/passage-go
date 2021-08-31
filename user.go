package passage

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/resty.v1"
)

type User struct {
	ID            string    `json:"id"`
	Active        bool      `json:"active"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	LastLogin     time.Time `json:"last_login_at"`
}

func (a *App) GetUser(userID string) (*User, error) {
	type respUser struct {
		User User `json:"user"`
	}
	var userBody respUser

	response, err := resty.New().R().
		SetAuthToken(a.Config.APIKey).
		SetResult(&userBody).
		Get(fmt.Sprintf("https://api.passage.id/v1/apps/%v/users/%v", a.ID, userID))
	if err != nil {
		return nil, errors.New("network error: could not get Passage User")
	}
	if response.StatusCode() == http.StatusNotFound {
		return nil, fmt.Errorf("passage User with ID \"%v\" does not exist", userID)
	}
	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get Passage User")
	}
	user := userBody.User

	return &user, nil
}
