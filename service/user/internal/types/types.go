// Code generated by goctl. DO NOT EDIT.
package types

type UserRequest struct {
	Username string `path:"username"`
	Token    string `json:"token"`
}

type UserReply struct {
	Status_code    int64  `json:"statusCode"`
	Status_msg     string `json:"statusMsg"`
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	Follow_count   int64  `json:"followCount"`
	Follower_count int64  `json:"followerCount"`
	Is_follow      bool   `json:"isFollow"`
}
