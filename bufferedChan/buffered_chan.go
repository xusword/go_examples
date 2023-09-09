package bufferedChan

type BufferedChan[T interface{}] struct {
	in  chan T
	out chan T
}

func NewBufferedChan[T interface{}]() *BufferedChan[T] {
	c := &BufferedChan[T]{
		in:  make(chan T),
		out: make(chan T),
	}
	go c.start()

	return c
}

func (c *BufferedChan[T]) Push(val T) {
	c.in <- val
}

func (c *BufferedChan[T]) NoMas() {
	close(c.in)
}

func (c *BufferedChan[T]) Pull() (T, bool) {
	v, ok := <-c.out
	return v, ok
}

func (c *BufferedChan[T]) start() {
	buffer := []T{}
	// out channel should not be hooked to our select statement if buffer is empty
	out := func() chan T {
		if len(buffer) == 0 {
			return nil
		} else {
			return c.out
		}
	}
	// head value doesn't matter when buffer is empty
	head := func() T {
		if len(buffer) == 0 {
			var null T
			return null
		} else {
			return buffer[0]
		}
	}
	for len(buffer) > 0 || c.in != nil {
		select {
		case val, ok := <-c.in:
			if ok {
				buffer = append(buffer, val)
			} else {
				c.in = nil
			}
		case out() <- head():
			buffer = buffer[1:]
		}
	}
	close(c.out)
}
