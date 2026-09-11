package vertex

import (
	"image/color"
	"os"
	"reflect"
	"testing"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

type pixel struct {
	X     int `gremlin:"filter"`
	Y     int `gremlin:"filter"`
	Color color.RGBA
}

func TestRepository_Search(t *testing.T) {
	remote, err := gremlingo.NewDriverRemoteConnection(os.Getenv("GREMLIN_URL"))
	if err != nil {
		t.Fatal(err)
	}

	repo := New[pixel]("Pixels", remote)

	in := &pixel{
		X: 10,
		Y: 14,
		Color: color.RGBA{
			R: 0,
			G: 64,
			B: 130,
			A: 10,
		},
	}

	vin, err := repo.New(in)

	if err != nil {
		t.Fatal(err)
	} else if vin.ID() == nil {
		t.Fail()
	}

	out := &pixel{
		X: in.X,
		Y: in.Y,
	}

	vout, err := repo.Search(out)
	if err != nil {
		t.Fatal(err)
	} else if vout.ID() == nil {
		t.Fail()
	} else if !reflect.DeepEqual(vin, vout) {
		t.Fail()
	}
}
