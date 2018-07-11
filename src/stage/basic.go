package stage

//ToString stage:type assertion, Change the data type of the channel from interface {} to string
func ToString(
	done <-chan interface{},
	valueStream <-chan interface{},
) <-chan string {
	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}

	}()
	return stringStream
}
