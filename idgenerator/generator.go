package idgenerator

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// New returns an ID generator which generate a unique id.
func New(componentName, componentID string) IDGenerator {
	return &generator{
		componentID:   componentID,
		componentName: componentName,
	}
}

// IDGenerator interface
type IDGenerator interface {
	NextID() string
}

type generator struct {
	componentID   string
	componentName string
}

// Generate a unique id with following format: <componentName>-<componentID>-<time>-<random number>
func (g *generator) NextID() string {
	var id = rand.Uint64()
	var time = strconv.FormatInt(time.Now().Unix(), 10)
	return fmt.Sprintf("%s-%s-%s-%020d", g.componentName, g.componentID, time, id)
}
