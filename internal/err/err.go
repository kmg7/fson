package err

type AppError struct {
	Internal bool
	Code     string
	Messages []string
	Err      error
}
