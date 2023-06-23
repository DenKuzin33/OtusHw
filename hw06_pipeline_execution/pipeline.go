package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	prev := in
	result := make(chan interface{})

	for _, stage := range stages {
		prev = stage(prev)
	}

	go func() {
		for {
			select {
			case <-done:
				close(result)
				return
			case v, opened := <-prev:
				if opened {
					result <- v
				} else {
					close(result)
					return
				}
			}
		}
	}()
	return result
}
