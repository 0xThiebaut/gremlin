package edge

import (
	"os"
	"testing"

	"github.com/0xThiebaut/gremlin/vertex"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

type person struct {
	Name string
	Age  int
}

type knows struct {
	Reason string
}

func TestRepository_New(t *testing.T) {
	remote, err := gremlingo.NewDriverRemoteConnection(os.Getenv("GREMLIN_URL"))
	if err != nil {
		t.Fatal(err)
	}

	persons := vertex.New[person]("Person", remote)

	john, err := persons.New(&person{
		Name: "John Do",
		Age:  18,
	})
	if err != nil {
		t.Fatal(err)
	}

	katie, err := persons.New(&person{
		Name: "Jane Do",
		Age:  20,
	})
	if err != nil {
		t.Fatal(err)
	}

	repo := New[knows]("Knows", remote)

	rel, err := repo.New(john, &knows{Reason: "Colleague"}, katie)
	if err != nil {
		t.Fatal(err)
	} else if rel.ID() == nil {
		t.Fail()
	}
}
