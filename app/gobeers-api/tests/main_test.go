package tests

import (
	"fmt"
	"testing"

	"github.com/phbpx/gobeers/business/data/dbtest"
	"github.com/phbpx/gobeers/foundation/docker"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}
