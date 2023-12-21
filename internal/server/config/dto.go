package config

import "time"

type signIn struct {
	Username *string `json:"username" validate:"required"`
	Password *string `json:"password" validate:"required"`
}

type adminUpdate struct {
	Password    *string `json:"password" validate:"required"`
	Username    *string `json:"username"`
	NewPassword *string `json:"newPassword"`
}

type appConfig struct {
	AutoStart *bool   `json:"autoStart"`
	UploadDir *string `json:"uploadDir" validate:"dir"`
	TempDir   *string `json:"tempDir" validate:"dir"`
}

type transferPath struct {
	Path *string `json:"path" validate:"required"`
}

type authConfig struct {
	TokensExpiresAfter  *time.Duration `json:"tokensExpiresAfter"`
	TokenExpireTolerant *time.Duration `json:"tokenExpireTolerant"`
}

type groupUpdate struct {
	Name  *string   `json:"name"`
	Users *[]string `json:"users"`
}

type groupCreate struct {
	Name  *string   `json:"name" validate:"required"`
	Users *[]string `json:"users" validate:"required"`
}
