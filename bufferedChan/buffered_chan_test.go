package bufferedChan

import (
	"testing"
	"time"
)

func TestHappy(t *testing.T) {
	c := NewBufferedChan[int]()
	c.Push(1)
	c.Push(2)
	c.Push(3)

	one, ok := c.Pull()
	assertOK(t, ok)
	two, ok := c.Pull()
	assertOK(t, ok)
	three, ok := c.Pull()
	assertOK(t, ok)

	c.Push(4)
	c.Push(5)

	four, ok := c.Pull()
	assertOK(t, ok)

	c.Push(6)

	five, ok := c.Pull()
	assertOK(t, ok)

	c.Push(7)
	assertOK(t, ok)
	c.Push(8)
	assertOK(t, ok)

	six, ok := c.Pull()
	assertOK(t, ok)
	c.NoMas()

	seven, ok := c.Pull()
	assertOK(t, ok)
	eight, ok := c.Pull()
	assertOK(t, ok)

	t.Logf("%d %d %d %d %d %d %d %d", one, two, three, four, five, six, seven, eight)
}

func TestBlock(t *testing.T) {
	c := NewBufferedChan[int]()
	go func() {
		time.Sleep(3 * time.Second)
		c.Push(1)
	}()
	t.Logf("Heading into a block")
	one, _ := c.Pull()

	t.Logf("%d", one)
}

func assertOK(t *testing.T, ok bool) {
	if !ok {
		t.Fatalf("Not OK")
	}
}
