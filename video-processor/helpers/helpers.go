package helpers

import (
	"os/exec"
	"sync"
)

func RunCmd(cmd *exec.Cmd) (string, error) {
	b, err := cmd.CombinedOutput()
	if err != nil {
		return string(b), err
	}

	return string(b), nil
}

type ThreadElement struct {
	Element interface{}
}

func Threadify(numOfThreads int, elements []ThreadElement, f func(args ...interface{})) {
	length := len(elements)
	each := length / numOfThreads
	acc := length - (numOfThreads * each)

	var wg sync.WaitGroup

	wg.Add(numOfThreads)

	start := 0

	for i := 0; i < numOfThreads; i++ {
		running := each

		if acc > 0 {
			running++
			acc--
		}

		go func(start, i int) {
			for j := 0; j < running; j++ {
				e := elements[start]
				f(e.Element)
				start++
			}

			wg.Done()
		}(start, i)

		start += running
	}

	wg.Wait()
}
