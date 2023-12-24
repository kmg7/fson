package transfer

type signIn struct {
	Username *string `json:"username" validate:"required"`
	Password *string `json:"password" validate:"required"`
}

type userUpdate struct {
	Password    *string `json:"password" validate:"required"`
	Username    *string `json:"username"`
	NewPassword *string `json:"newPassword"`
}
