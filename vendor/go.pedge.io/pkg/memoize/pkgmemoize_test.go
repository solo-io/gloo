package pkgmemoize

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFibonacciEMemoFunc(t *testing.T) {
	calls := 0
	fibonacci := NewEMemoFunc(
		func(i int, f EMemoFunc) (interface{}, error) {
			calls++
			if i == 0 {
				return uint64(0), nil
			}
			if i == 1 {
				return uint64(1), nil
			}
			n1, err := f(i - 1)
			if err != nil {
				return 0, err
			}
			n2, err := f(i - 2)
			if err != nil {
				return 0, err
			}
			return n1.(uint64) + n2.(uint64), nil
		},
	)
	result, err := fibonacci(93)
	require.NoError(t, err)
	require.Equal(t, uint64(12200160415121876738), result.(uint64))
	require.Equal(t, 94, calls)
}

func TestFibonacciEMemoFuncError(t *testing.T) {
	calls := 0
	fibonacci := NewEMemoFunc(
		func(i int, f EMemoFunc) (interface{}, error) {
			calls++
			if i == 0 {
				return uint64(0), nil
			}
			if i == 1 {
				return uint64(1), nil
			}
			if i == 5 {
				return 0, fmt.Errorf("5")
			}
			n1, err := f(i - 1)
			if err != nil {
				return 0, err
			}
			n2, err := f(i - 2)
			if err != nil {
				return 0, err
			}
			return n1.(uint64) + n2.(uint64), nil
		},
	)
	_, err := fibonacci(93)
	require.Error(t, err, "5")
	require.Equal(t, 89, calls)
}

func TestFibonacciMemoFunc(t *testing.T) {
	calls := 0
	fibonacci := NewMemoFunc(
		func(i int, f MemoFunc) interface{} {
			calls++
			if i == 0 {
				return uint64(0)
			}
			if i == 1 {
				return uint64(1)
			}
			return f(i-1).(uint64) + f(i-2).(uint64)
		},
	)
	require.Equal(t, uint64(12200160415121876738), fibonacci(93))
	require.Equal(t, 94, calls)
}
