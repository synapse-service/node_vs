package plugin

type Settings struct {
	ID string
}

func (s Settings) IsEmpty() bool { return s.ID == "" }
