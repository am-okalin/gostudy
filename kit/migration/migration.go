package migration

type FileMarshaler interface {
	Marshal(v any) (fname string, err error)
}

type FileUnmarshaler interface {
	Unmarshal(fname string, v any) error
}

func Unmarshal(fname string, v any) error {
	return nil
}

func Marshal(v any) (fname string, err error) {
	return "", nil
}
