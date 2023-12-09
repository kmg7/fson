package server

import "time"

type AppConfigDTO struct {
	AutoStart bool   `json:"autoStart"`
	UploadDir string `json:"uploadDir" validate:"required,dir"`
	TempDir   string `json:"tempDir" validate:"required,dir"`
}
type TransferPathDTO struct {
	Path string `json:"path" validate:"required"`
}

type TransferConfigDTO struct {
	Transfer []TransferPathDTO `json:"transfer" validate:"required,dive"`
}

type AuthConfigDTO struct {
	TokensExpiresAfter  time.Duration `json:"tokensExpiresAfter"`
	TokenExpireTolerant time.Duration `json:"tokenExpireTolerant"`
}

type AdminSignInDTO struct {
	Username *string `json:"username" validate:"required"`
	Password *string `json:"password" validate:"required"`
}

type AdminUpdateDTO struct {
	Password    *string `json:"password" validate:"required"`
	Username    *string `json:"username"`
	NewPassword *string `json:"newPassword"`
}

type GroupUpdateDTO struct {
	Name  *string   `json:"name"`
	Users *[]string `json:"users"`
}

type GroupCreateDTO struct {
	Name  *string   `json:"name" validate:"required"`
	Users *[]string `json:"users" validate:"required"`
}
