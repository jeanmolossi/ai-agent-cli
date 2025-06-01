package foundation

type Json interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
	MarshalString(any) (string, error)
	UnmarshalString(string, any) error
}
