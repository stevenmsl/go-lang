/*
A pipeline consists of stages. Each stage will take the data in,
perform a transformation on it, and then send the data back out.

Stage
- A stage consumes and returns the same type.
- A stage must be reified by the language so that it may be passed around.
  (Need to do more research on this. (nomad, high order functions))

Batch processing
- Each stage will process chunks of data at once

Stream processing
- Each stage will receive and emit one element at a time



*/

package main

import (
	"fmt"
	"math/rand"
	"stage"

	"time"
)

func main() {
	runPipelineString()
	//runPipelineTakeFn()
	//runPipelineTake()
	//runPipelineC()
	//runPipelineS()
	//runPipelineBP()
}

func runPipelineString() {
	done := make(chan interface{})
	defer func() {
		close(done)
		fmt.Println("runPipelineTakeFn closed")
		time.Sleep(3 * time.Second)
	}()

	var message string
	for token := range stage.ToString(done, stage.Take(done, stage.Repeat(done, "I", "am."), 5)) {
		message += token
	}
	fmt.Printf("message: %s...\n", message)

}

func runPipelineTakeFn() {
	done := make(chan interface{})
	defer func() {
		close(done)
		fmt.Println("runPipelineTakeFn closed")
		time.Sleep(3 * time.Second) //Wait a bit so you can see the message printed when a stage is closed.
	}()
	rand := func() interface{} { return rand.Int() }
	for num := range stage.Take(done, stage.RepeatFn(done, rand), 10) {
		fmt.Printf("%v\n", num)
	}

}

func runPipelineTake() {
	done := make(chan interface{})
	defer func() {
		close(done)
	}()

	for num := range stage.Take(done, stage.Repeat(done, 1, 2), 10) {
		fmt.Printf("%v ", num)
	}
	fmt.Println()
}

/*
When you look at the output, it looks quite random:
Generated the number 1
1 + 2 = 2 done by stage Added by 2
3 x 3 = 9 done by stage Multiplied by 3
1 x 1 = 1 done by stage Multiplied by 1

In the above example, it appears to be in iteration 0 stage “Added by 2” finished before stage “Multiplied by 1”,
which it’s not possible as the later stage “Added by 2” is depending on the outcome of the previous stage “Multiplied by 1”.
I think the reason for this inconsistence is due to the time needed between you pushed the result to the channel
and you printed the result.
Don’t forget all stages are executing concurrently, and the next stage probably already finished
by the time you finished printing the result in the current stage.

*/
func runPipelineC() {
	done := make(chan interface{})
	defer close(done)
	intStream := stage.Generator(done, 1, 2, 3, 4)
	//each stage of the pipeline is executing concurrently
	pipeline := stage.MultiplyC(done, stage.AddC(done, stage.MultiplyC(done, intStream, 1), 2), 3)
	index := 0
	//If the range expression is a channel, at most one iteration variable is permitted.
	for v := range pipeline {
		fmt.Printf("Iteration %v Completed. Produced the number %v\n", index, v)
		index++
	}
}

func runPipelineS() {
	ints := []int{1, 2, 3, 4}
	for _, v := range ints {
		//We instantiate the pipeline for each iteration: 4 times in this case
		fmt.Println(stage.MultiplyS(stage.AddS(stage.MultiplyS(v, 2), 1), 2))
	}
}

func runPipelineBP() {
	ints := []int{1, 2, 3, 4}
	for _, v := range stage.AddBP(stage.MultiplyBP(ints, 2), 1) {
		fmt.Println(v)
	}
}
