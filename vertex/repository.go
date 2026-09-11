package vertex

import (
	"fmt"
	"reflect"

	"github.com/0xThiebaut/gremlin"
	"github.com/0xThiebaut/gremlin/internal/element"
	"github.com/0xThiebaut/gremlin/internal/flatten"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

type Repository[Properties any] interface {
	New(*Properties) (gremlin.Element[Properties], error)
	Search(*Properties) (gremlin.Element[Properties], error)
	gremlin.Repository[Properties]
}

func New[Properties any](label string, remote *gremlingo.DriverRemoteConnection) Repository[Properties] {
	return &repository[Properties]{
		Label:      label,
		Remote:     remote,
		Repository: element.Repository[Properties](element.Vertex, label, remote),
	}
}

type repository[Properties any] struct {
	Label      string
	Remote     *gremlingo.DriverRemoteConnection
	Repository gremlin.Repository[Properties]
}

func (r *repository[Properties]) New(vertex *Properties) (gremlin.Element[Properties], error) {
	// create the vertex
	t := gremlingo.Traversal_().With(r.Remote).AddV(r.Label)

	// Set its properties
	props, err := flatten.Marshal(vertex)
	if err != nil {
		return nil, err
	}

	list := make(map[string]bool)
	for _, prop := range props {
		c := gremlingo.Cardinality.Single
		if list[prop.Label] {
			c = gremlingo.Cardinality.List
		} else {
			list[prop.Label] = true
		}
		t = t.Property(c, prop.Label, prop.Value)
	}

	// execute
	result, err := t.Next()
	if err != nil {
		return nil, err
	}

	// get identifier
	e, err := result.GetVertex()
	if err != nil {
		return nil, err
	}

	return element.New(e.Id, e.Label, vertex), nil
}

func (r *repository[Properties]) Search(vertex *Properties) (gremlin.Element[Properties], error) {
	// create the vertex
	t := gremlingo.Traversal_().With(r.Remote).V().HasLabel(r.Label)

	// Set its properties
	props, err := flatten.MarshalFilter(vertex, func(tag reflect.StructTag) bool {
		settings, ok := tag.Lookup("gremlin")
		return ok && settings == "filter"
	})
	if err != nil {
		return nil, err
	}

	for _, prop := range props {
		t = t.Has(prop.Label, prop.Value)
	}

	// execute
	result, err := t.Next()
	if err != nil {
		return nil, err
	}

	// get identifier
	v, err := result.GetVertex()
	if err != nil {
		return nil, err
	}

	props = props[:0]
	anys, ok := v.Properties.([]any)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T", v.Properties)
	}
	for _, a := range anys {
		switch prop := a.(type) {
		case gremlingo.VertexProperty:
			props = append(props, flatten.Property{Label: prop.Element.Label, Value: prop.Value})
		case *gremlingo.VertexProperty:
			props = append(props, flatten.Property{Label: prop.Element.Label, Value: prop.Value})
		default:
			return nil, fmt.Errorf("unsupported type %T", a)
		}
	}

	return element.New(v.Id, v.Label, vertex), flatten.Unmarshal(props, vertex)
}

func (r *repository[Properties]) Save(e gremlin.Element[Properties]) error {
	return r.Repository.Save(e)
}
