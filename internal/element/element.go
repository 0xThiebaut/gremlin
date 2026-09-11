package element

import "github.com/0xThiebaut/gremlin"

func New[Properties any](id any, label string, properties *Properties) gremlin.Element[Properties] {
	return &element[Properties]{
		id:         id,
		label:      label,
		properties: properties,
	}
}

type element[Properties any] struct {
	id         any
	label      string
	properties *Properties
}

func (e *element[Properties]) ID() any {
	return e.id
}

func (e *element[Properties]) Label() string {
	return e.label
}

func (e *element[Properties]) Properties() *Properties {
	return e.properties
}
