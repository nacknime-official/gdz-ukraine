package entity

type Subject struct {
	Name string
}
type Author struct {
	Name string
}
type Specification struct {
	Name string
}
type Year struct {
	Year int
}
type Topic struct {
	Name string
}
type Exercise struct {
	Name string
}
type Solution struct {
	// TODO
}

// TODO: better name
type Opts struct {
	Grade         int
	Subject       *Subject
	Author        *Author
	Specification *Specification
	Year          *Year
	Topics        []*Topic
	Exercise      *Exercise
}
