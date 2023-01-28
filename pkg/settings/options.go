package settings

type option = func(*Settings) error

func WithFileName(filename string) option {
	return func(s *Settings) error {
		s.filename = filename
		return nil
	}
}
