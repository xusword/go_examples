package bufferedChan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHappy(t *testing.T) {
	c := NewBufferedChan[int]()
	c.Push(1)
	c.Push(2)
	c.Push(3)

	var val int
	var ok bool

	val, ok = c.Pull()
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	val, ok = c.Pull()
	assert.True(t, ok)
	assert.Equal(t, 2, val)

	val, ok = c.Pull()
	assert.True(t, ok)
	assert.Equal(t, 3, val)

	c.Push(4)
	c.Push(5)

	val, ok = c.Pull()
	assert.True(t, ok)
	assert.Equal(t, 4, val)

	c.Push(6)

	val, ok = c.Pull()
	assert.True(t, ok)
	assert.Equal(t, 5, val)

	c.Push(7)
	c.Push(8)

	val, ok = c.Pull()
	assert.True(t, ok)
	assert.Equal(t, 6, val)
	c.NoMas()

	val, ok = c.Pull()
	assert.True(t, ok)
	assert.Equal(t, 7, val)

	val, ok = c.Pull()
	assert.True(t, ok)
	assert.Equal(t, 8, val)

	_, ok = c.Pull()
	assert.False(t, ok)
}

func TestBlock(t *testing.T) {
	c := NewBufferedChan[int]()
	go func() {
		time.Sleep(10 * time.Microsecond)
		c.Push(1)
	}()
	t.Logf("Heading into a block")
	val, ok := c.Pull()
	assert.True(t, ok)
	assert.Equal(t, 1, val)
}
