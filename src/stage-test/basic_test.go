/*
- Benchmark functions start with Benchmark.
- Benchmark functions are run several times by the testing package.
  The value of b.N will increase each time until the benchmark runner is satisfied with the stability of the benchmark.
- Each benchmark must execute the code under test b.N times.
*/

package stage_test

import (
	"stage"
	"testing"
)

//BenchmarkGeneric ...
func BenchmarkGeneric(b *testing.B) {
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()
	for range stage.ToString(done, stage.Take(done, stage.Repeat(done, "a"), b.N)) {

	}
}

func BenchmarkTyped(b *testing.B) {
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()
	for range stage.TakeS(done, stage.RepeatS(done, "a"), b.N) {

	}

}
