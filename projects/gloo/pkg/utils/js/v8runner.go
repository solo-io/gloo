package js

// #cgo LDFLAGS: -static-libstdc++
// #cgo linux LDFLAGS: -static-libgcc
import "C"

import (
	"sync"
	"time"

	"github.com/rotisserie/eris"
	v8 "rogchap.com/v8go"
)

var v8goMutex = &sync.Mutex{}

// Runner is any program that allows you to set the timeout and run the pogram. Use the Run() method to run the program.
type Runner interface {
	// SetTimeout will set the timeout for the program. If the program runs longer then the timeout, then the program will stop.
	SetTimeout(timeout time.Duration)
	// Run will run the program
	Run(input string) (string, error)
}

// V8RunnerInputOutput will attach a getInput method and a returnOutput method.
// getInput() method will return the input on the Run(input) method. The input value will be a string.
// returnOuput() return the output of the program. The return value will be a string.
type V8RunnerInputOutput struct {
	// sourceCode is the raw source code of the js file
	sourceCode string
	// filename the name of js file to run
	filename string
	// cache is the compiled v8 cache
	cache *v8.CompilerCachedData
	// timeout is the timeoiut to run the program, else it will fail with a timeout error. The default is 10 seconds.
	// Use SetTimeout to set.
	timeout time.Duration
}

func newIsolate() *v8.Isolate {
	v8goMutex.Lock()
	iso := v8.NewIsolate()
	v8goMutex.Unlock()
	return iso
}

func disposeIsolate(iso *v8.Isolate) {
	v8goMutex.Lock()
	iso.Dispose()
	v8goMutex.Unlock()
}

func newContext(iso *v8.Isolate, global *v8.ObjectTemplate) *v8.Context {
	// needs to have mutex, if isolate is not applied
	v8ctx := v8.NewContext(iso, global)
	return v8ctx
}

func NewV8RunnerInputOutput(filename, sourceCode string) (*V8RunnerInputOutput, error) {
	iso := newIsolate()
	defer disposeIsolate(iso)

	script, err := iso.CompileUnboundScript(sourceCode, filename, v8.CompileOptions{})
	if err != nil {
		return nil, eris.Wrapf(err, "unable to create V8RunnerInputOutput and compile %s", filename)
	}
	cacheData := script.CreateCodeCache()

	return &V8RunnerInputOutput{
		filename:   filename,
		sourceCode: sourceCode,
		cache:      cacheData,
		timeout:    time.Second * 10,
	}, nil
}

func (r *V8RunnerInputOutput) SetTimeout(timeout time.Duration) {
	r.timeout = timeout
}

// Run will run the JS program using v8go.
func (r *V8RunnerInputOutput) Run(input string) (string, error) {
	iso := newIsolate()
	defer disposeIsolate(iso)
	global := v8.NewObjectTemplate(iso) // a template that represents a JS Object
	if err := r.setInput(iso, global, input); err != nil {
		return "", err
	}
	outputChan, err := r.setOutput(iso, global)
	if err != nil {
		return "", err
	}
	v8ctx := newContext(iso, global)
	completed, errs := r.runScript(v8ctx, iso, global)
	defer close(completed)
	defer close(errs)

	select {
	case <-completed:
		v := <-outputChan
		return v, nil
	case err = <-errs:
		if err != nil {
			return "", eris.Wrapf(err, "error running %s", r.filename)
		}
	case <-time.After(r.timeout):
		vm := v8ctx.Isolate()   // get the Isolate from the context
		vm.TerminateExecution() // terminate the execution
		err = <-errs            // will get a termination error back from the running script
		return "", eris.Wrapf(err, "error %s timed out after %v", r.filename, r.timeout)
	}
	// we should never hit this case
	return "", eris.New(r.filename + " failed")
}

// setInput will set the input for the program.
func (r *V8RunnerInputOutput) setInput(iso *v8.Isolate, global *v8.ObjectTemplate, input string) error {
	inputVal, err := v8.NewValue(iso, input) // you can return a value back to the JS caller if required
	if err != nil {
		return eris.Wrapf(err, "unable to create value for v8go %s program: %s", r.filename, input)
	}
	// define the getInput method
	err = global.Set("getInput", v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		return inputVal
	}))
	if err != nil {
		return eris.Wrap(err, "unable to set getStitchingInput function")
	}
	return nil
}

// setOuput will return the output channel, receives the output of the program.
func (r *V8RunnerInputOutput) setOutput(iso *v8.Isolate, global *v8.ObjectTemplate) (chan string, error) {
	val := make(chan string, 1)
	// define the returnOutput method
	err := global.Set("returnOutput", v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		val <- info.Args()[0].String()
		return nil
	}))
	if err != nil {
		return nil, eris.Wrapf(err, "unable to set returnOutput function for %s", r.filename)
	}
	return val, nil
}

// runScript this will run the script. Returns a done channel and an error channel.
// ASYNC
func (r *V8RunnerInputOutput) runScript(v8ctx *v8.Context, iso *v8.Isolate, global *v8.ObjectTemplate) (chan *v8.Value, chan error) {
	completed := make(chan *v8.Value, 1)
	errs := make(chan error, 1)

	script, err := iso.CompileUnboundScript(r.sourceCode, r.filename, v8.CompileOptions{
		CachedData: r.cache,
	})
	// if there is no error then proceed running the script in a go routine
	if err == nil {
		go func() {
			// we don't need any output from the script, just the side effect of calling
			done, err := script.Run(v8ctx)
			if err != nil {
				errs <- err
				return
			}
			completed <- done
		}()
	} else {
		errs <- eris.Wrapf(err, "unable to compile %s", r.filename)
	}
	return completed, errs
}
