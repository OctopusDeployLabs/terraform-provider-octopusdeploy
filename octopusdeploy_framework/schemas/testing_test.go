package schemas

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var createSharedContainer = flag.Bool("createSharedContainer", false, "Set to true to run integration tests in containers")

func TestSomething(t *testing.T) {
	log.Print("Hi")
	assert.Equal(t, 1, 1)
}
