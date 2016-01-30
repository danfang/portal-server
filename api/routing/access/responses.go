package access

// A LoginResponse contains credentials to make authenticated requests.
//
// swagger:response loginResponse
type LoginResponse struct {
	// in: body
	Body loginResponse
}

type loginResponse struct {
	UserToken string `json:"user_token"`
	UserUUID  string `json:"user_id"`
}
