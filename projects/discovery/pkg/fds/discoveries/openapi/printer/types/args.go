package types

import . "github.com/graphql-go/graphql"

type PrintDescriptionParams struct {
	Description  string
	Indentation  string
	FirstInBlock bool
}

func NewPrintDescriptionParams() *PrintDescriptionParams {
	return &PrintDescriptionParams{
		FirstInBlock: true,
	}
}

type PrintArgsParams struct {
	Args        []*Argument
	Indentation string
}
