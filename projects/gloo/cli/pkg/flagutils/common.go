package flagutils

import "github.com/spf13/pflag"

func AddOutputFlag(set *pflag.FlagSet, strptr *string) {
	set.StringVarP(strptr, "output", "o", "", "output format: (yaml, json, table)")
}

func AddFileFlag(set *pflag.FlagSet, strptr *string) {
	set.StringVarP(strptr, "file", "f", "", "file to be read or written to")
}
