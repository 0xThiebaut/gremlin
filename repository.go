package gremlin

type Repository[Properties any] interface {
	Save(Element[Properties]) error
}
