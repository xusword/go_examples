package bufferedChan

type BufferedChan[T interface{}] struct {
	In  chan T
	Out chan T
}

func NewBufferedChan[T interface{}]() *BufferedChan[T] {
	c := &BufferedChan[T]{
		In:  make(chan T),
		Out: make(chan T),
	}
	go c.start()

	return c
}

func (c *BufferedChan[T]) start() {
	draining := false
	buffer := []T{}
	headChan := make(chan T)
	for {
		select {
		case val, ok := <-c.In:
			if ok {
				select {
				// if value came in when headChan is occupied, push to buffer
				case head := <-headChan:
					buffer = append(buffer, head)
				default:
					// otherwise move on
				}
				go func() { headChan <- val }()
			} else {
				draining = true
			}
		case c.Out <- <-headChan:
			if len(buffer) > 0 {
				head := buffer[0]
				buffer = buffer[1:]
				go func() {
					headChan <- head
				}()
			} else {
				if draining {
					close(c.Out)
					return
				}
			}
		}
	}
}
