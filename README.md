# Gremlin
A proof-of-concept abstraction library for Gremlin ingestion.
You should not run this in production or rely on this library.

## Example Usage
```go
package main

import (
	"image/color"
	"log"
	"os"
	"time"

	"github.com/0xThiebaut/gremlin/edge"
	"github.com/0xThiebaut/gremlin/vertex"
	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

type Person struct {
	Name string `gremlin:"filter"`
	Age uint8
}

type Action struct {
	Time time.Time
}

type Pixel struct {
	X     int `gremlin:"filter"`
	Y     int `gremlin:"filter"`
	Color color.RGBA
}

func main() {
	// create a connection
	remote, err := gremlingo.NewDriverRemoteConnection(os.Getenv("GREMLIN_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	// define the repositories
	pixels := vertex.New[Pixel]("Pixel", remote)
	artists := vertex.New[Person]("Artist", remote)
	authored := edge.New[Action]("Drew", remote)

	// create a new artist
	_, err = artists.New(&Person{
		Name: "John Do",
		Age: 18,
	})

	if err != nil {
		log.Fatalln(err)
	}

	// create a new pixel
	pixel, err := pixels.New(&Pixel{
		X: 10,
		Y: 14,
		Color: color.RGBA{
			R: 0,
			G: 64,
			B: 130,
			A: 10,
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	// search for an existing artist
	john := &Person{Name: "John Do"}
	artist, err := artists.Search(john)
	if err != nil {
		log.Fatalln(err)
	}

	// set the pixel's author
	_, err = authored.New(artist, &Action{Time: time.Now()}, pixel)
	if err != nil {
		log.Fatalln(err)
	}
}
```