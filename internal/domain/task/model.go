package task

const (
	StateCreate    = "create"
	StateCompleted = "completed"
	StateFail      = "fail"
)

type Entity struct {
	UserID    string
	Topic     string
	MessageID string
	Message   string
	State     string
}
