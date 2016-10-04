package konturtransferbot

// Logger is a simple logging wrapper interface
type Logger interface {
	Log(...interface{}) error
}
