package userdto

type ReqLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=6"`
}

type ResLogin struct {
	Jwt string `json:"jwt"`
}
