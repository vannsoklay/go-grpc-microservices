package dto

type RegisterReq struct {
	Name     string `josn:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
