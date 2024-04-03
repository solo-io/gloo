package kubectl

import "io"

type Kubectl struct {
	Receiver io.Writer
}

func MustKubectl(receiver io.Writer) *Kubectl {
	if receiver == nil {
		panic("receiver must not be nil")
	}
	return &Kubectl{
		Receiver: receiver,
	}
}
