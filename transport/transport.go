package transport

type Sender interface {
	Send([]byte) error
}
