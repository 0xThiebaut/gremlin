package element

import (
	"fmt"
	"reflect"

	"github.com/0xThiebaut/gremlin"
	"github.com/0xThiebaut/gremlin/internal/flatten"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

func Repository[Properties any](k kind, label string, remote *gremlingo.DriverRemoteConnection) gremlin.Repository[Properties] {
	return &repository[Properties]{
		Kind:   k,
		Label:  label,
		Remote: remote,
	}
}

type kind string

const (
	Edge   kind = "edge"
	Vertex kind = "vertex"
)

type repository[Properties any] struct {
	Kind   kind
	Label  string
	Remote *gremlingo.DriverRemoteConnection
}

func (r *repository[Properties]) Save(e gremlin.Element[Properties]) error {
	if e.Label() != r.Label {
		return fmt.Errorf("invalid label")
	}
	if id := reflect.ValueOf(e.ID()); !id.IsValid() || id.IsNil() || id.IsZero() {
		return fmt.Errorf("invalid identifier")
	}
	// select the element using V(e.ID()) or E(e.ID())
	t := gremlingo.Traversal_().With(r.Remote).GetGraphTraversal()
	switch r.Kind {
	case Edge:
		t = t.E(e.ID())
	case Vertex:
		t = t.V(e.ID())
	default:
		return fmt.Errorf("unknown element kind %q", r.Kind)
	}
	// enforce a label check
	t = t.HasLabel(r.Label)
	// drop existing properties
	t = t.SideEffect(gremlingo.T__.Properties().Drop())

	// save new properties
	props, err := flatten.Marshal(e.Properties())
	if err != nil {
		return err
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
	return <-t.Iterate()
}
