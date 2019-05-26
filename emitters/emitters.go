package emitters

type Emitter interface {
	Emit(commands ...Command) error
}

const (
	Panasonic Encoding = 241
)

type Encoding int

type Command struct {
	Encoding Encoding
	Payload  []byte
}