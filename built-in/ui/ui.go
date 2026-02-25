package ui

type UI interface {
	PromptBool(question string, suggested bool) (bool, ok bool)
	PromptString(question, suggested string) (string, ok bool)
	HandleCrash()
}
