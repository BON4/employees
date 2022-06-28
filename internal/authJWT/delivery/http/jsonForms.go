package http

type LoginForm struct {
	Username string
	Password string
}

type RegisterForm struct {
	Username string
	Password string
}

type Cookies struct {
	Access_Token  string `json:"access_token"`
	Refresh_Token string `json:"refresh_token"`
}
