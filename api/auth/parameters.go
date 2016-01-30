package auth

// AuthenticationHeaders are the required headers to make authenticated
// requests.
//
// swagger:parameters addDevice messageHistory getDevices
type AuthenticationHeaders struct {
	// in: header
	// required: true
	UserToken string `json:"X-USER-TOKEN"`

	// in: header
	// required: true
	UserID string `json:"X-USER-ID"`
}
