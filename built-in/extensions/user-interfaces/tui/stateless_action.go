package tui

type SimpleAction struct {
	Name        string
	Description string
	ShortCut    rune
	FnToRun     func()
}

func (s *SimpleAction) GetName() string {
	return s.Name
}

func (s *SimpleAction) GetDescription() string {
	return s.Description
}

func (s *SimpleAction) GetShortCut() rune {
	return s.ShortCut
}

func (s *SimpleAction) Run() {
	if s.FnToRun != nil {
		s.FnToRun()
	}
}
