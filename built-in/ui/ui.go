package ui

type UI interface {
	PromptBool(question string) bool
	PromptString(question string) bool
	HandleCrash()
}
