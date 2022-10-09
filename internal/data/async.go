package data

// Does task asynchronously and sends result or error.
// Waits for result to be received from channel.
//
// Reserved for future use.
func Async[R any](task func() (R, error)) (<-chan R, <-chan error) {
	out := make(chan R)
	e := make(chan error)

	go func() {
		defer func() {
			close(out)
			close(e)
		}()

		result, err := task()
		if err != nil {
			e <- err
		} else {
			out <- result
		}
	}()

	return out, e
}
