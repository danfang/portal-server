package util

import (
	"encoding/json"
	"io/ioutil"
	"portal-server/api/errs"
)

// A GoogleUser is a user as represented by Google OAuth.
type GoogleUser struct {
	Sub           string `json:"sub"`
	Aud           string `json:"aud"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}

var (
	googleAUDs = []string{
		"1045304436932-9vtokstg18sq2hu26hipueithq7sb0bq.apps.googleusercontent.com", // Android
		"1045304436932-564pg9gi9lee05mg45frg7kigd7h5775.apps.googleusercontent.com", // Chrome
	}
)

// GetGoogleUser takes an id token and retrieves a user profile
// from Google and returns a GoogleUser.
func GetGoogleUser(gc *WebClient, idToken string) (*GoogleUser, error) {
	res, err := gc.HTTPClient.Get(gc.BaseURL + "?id_token=" + idToken)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 400 {
		return nil, errs.ErrInvalidGoogleIDToken
	}
	if res.StatusCode != 200 {
		return nil, errs.ErrGoogleOAuthUnavailable
	}
	var user GoogleUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}
	validAUD := checkAUD(&user)
	if !validAUD {
		return nil, errs.ErrInvalidGoogleIDToken
	}
	return &user, nil
}

func checkAUD(user *GoogleUser) bool {
	validAUD := false
	for _, aud := range googleAUDs {
		if aud == user.Aud {
			validAUD = true
		}
	}
	return validAUD
}
