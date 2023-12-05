package server

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
