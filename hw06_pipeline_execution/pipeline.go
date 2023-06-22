package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In, done In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	prev := in

	for _, stage := range stages {
		prev = stage(prev, done)
	}

	return prev
}
