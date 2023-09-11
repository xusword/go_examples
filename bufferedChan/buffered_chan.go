package bufferedChan

// credits: https://medium.com/capital-one-tech/building-an-unbounded-channel-in-go-789e175cd2cd

type BufferedChan[T any] struct {
	in     chan T
	out    chan T
	buffer []T
}

func NewBufferedChan[T any]() *BufferedChan[T] {
	c := &BufferedChan[T]{
		in:     make(chan T),
		out:    make(chan T),
		buffer: []T{},
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

func (c *BufferedChan[T]) receive(val T, ok bool) {
	if ok {
		c.buffer = append(c.buffer, val)
	} else {
		c.in = nil
	}
}

func (c *BufferedChan[T]) start() {
	for len(c.buffer) > 0 || c.in != nil {
		if len(c.buffer) == 0 {
			// if buffer empty you can only receive
			val, ok := <-c.in
			c.receive(val, ok)
		} else {
			// once buffer is not, we can receive or send
			select {
			case val, ok := <-c.in:
				c.receive(val, ok)
			// the referenced article used an outCh() function to make the lhs nil-able
			// This is a very nice trick and would be nice to have if we have multiple
			// case statements. However since we only 2 code paths I find having the
			// if statement before select being more concise
			case c.out <- c.buffer[0]:
				c.buffer = c.buffer[1:]
			}
		}
	}
	close(c.out)
}
