package access

type loginResponse struct {
	UserToken string `json:"user_token"`
	UserUUID  string `json:"user_id"`
}
