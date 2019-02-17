package fibonacci_tdd_kata

type Fibonacci struct{}

func (f Fibonacci) Fib(steps int) (result int) {
	if steps == 1 || steps == 2 {
		result = 1
	} else if steps > 2 {
		result = f.Fib(steps-1) + f.Fib(steps-2)
	}

	return
}
