package service

type option = func(*Service) error

func WithConfig(cfg Config) option {
	return func(s *Service) error {
		s.cfg = cfg
		return nil
	}
}

func WithSettings(settings Settings) option {
	return func(s *Service) error {
		s.settings = settings
		return nil
	}
}
