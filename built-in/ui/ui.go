package ui

type UI interface {
	PromptBool(question string, suggested bool) (v bool, ok bool)
	PromptString(question, suggested string) (v string, ok bool)
	HandleCrash()
}
