package plugin

type option = func(*Plugin) error

func WithSettings(settings Settings) option {
	return func(p *Plugin) error {
		p.settings = settings
		return nil
	}
}

func WithEventCallback(callback func(Event)) option {
	return func(p *Plugin) error {
		p.OnEvent(callback)
		return nil
	}
}

func WithExitCallback(callback func()) option {
	return func(p *Plugin) error {
		p.OnExit(callback)
		return nil
	}
}
