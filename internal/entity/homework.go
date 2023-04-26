package entity

type Subject struct {
	ID   string
	Name string
}
type Author struct {
	ID   string
	Name string
}
type Specification struct {
	ID   string
	Name string
}
type Year struct {
	ID   string
	Year int
}
type Topic struct {
	ID   string
	Name string
}
type Exercise struct {
	ID   string
	Name string
}
type TopicOrExercise struct {
	Topic    *Topic
	Exercise *Exercise
}
type Solution struct {
	// TODO
}

// TODO: better name
type Opts struct {
	Class         int
	Subject       *Subject
	Author        *Author
	Specification *Specification
	Year          *Year
	Topics        []*Topic
	Exercise      *Exercise
}
