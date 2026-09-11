package gremlin

type Identifiable interface {
	ID() any
	Label() string
}

type Element[Properties any] interface {
	Identifiable
	Properties() *Properties
}
