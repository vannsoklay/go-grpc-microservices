package domain

type User struct {
	ID       string
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

type UserRsp struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	Password string `json:"password_hash"`
}

