package base

// base.AppError defines app-level errors.
type AppError struct {
	Code string `json:"code,omitempty"`
	Err  error  `json:"error,omitempty"`
}

func (ae *AppError) Error() string {
	return ae.Err.Error()
}

func NewAppError(code string, err error) *AppError {
	return &AppError{code, err}
}

type Status struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func NewStatus(code string, msg string) *Status {
	return &Status{Code: code, Msg: msg}
}

var (
	SuccessStatus = *NewStatus("0", "")
)
