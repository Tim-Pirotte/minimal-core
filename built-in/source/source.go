package source

// TODO limit content to 4 GiB
type Sources struct {
	Sources []Source
}

type Source struct {
	Name    string
	Content string
}
