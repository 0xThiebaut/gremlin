package edge

import (
	"github.com/0xThiebaut/gremlin"
	"github.com/0xThiebaut/gremlin/internal/element"
	"github.com/0xThiebaut/gremlin/internal/flatten"
	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

type Repository[Properties any] interface {
	New(from gremlin.Identifiable, edge *Properties, to gremlin.Identifiable) (gremlin.Element[Properties], error)
	gremlin.Repository[Properties]
}

func New[Properties any](label string, remote *gremlingo.DriverRemoteConnection) Repository[Properties] {
	return &repository[Properties]{
		Label:      label,
		Remote:     remote,
		Repository: element.Repository[Properties](element.Edge, label, remote),
	}
}

type repository[Properties any] struct {
	Label      string
	Remote     *gremlingo.DriverRemoteConnection
	Repository gremlin.Repository[Properties]
}

func (r *repository[Properties]) New(from gremlin.Identifiable, edge *Properties, to gremlin.Identifiable) (gremlin.Element[Properties], error) {
	// create the edge
	t := gremlingo.Traversal_().With(r.Remote).AddE(r.Label).From(gremlingo.T__.V(from.ID())).To(gremlingo.T__.V(to.ID()))

	// Set its properties
	props, err := flatten.Marshal(edge)
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
	e, err := result.GetEdge()
	if err != nil {
		return nil, err
	}

	return element.New(e.Id, e.Label, edge), nil
}

func (r *repository[Properties]) Save(e gremlin.Element[Properties]) error {
	return r.Repository.Save(e)
}
