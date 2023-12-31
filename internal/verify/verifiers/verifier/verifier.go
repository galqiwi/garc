package verifier

type EMessageLevel int64

const (
	Warn EMessageLevel = iota
	Error
)

var MessageLevelMap = map[EMessageLevel]string{
	Error: "Error",
	Warn:  "Warn",
}

type ErrorMessage struct {
	Path   string
	Reason string
	Level  EMessageLevel
}

type Verifier interface {
	Verify() []ErrorMessage
	Name() string
}

func NewErrorMessageFromError(err error) ErrorMessage {
	return ErrorMessage{
		Path:   "",
		Reason: err.Error(),
		Level:  Error,
	}
}
