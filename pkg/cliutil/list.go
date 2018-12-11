package cliutil

// TODO(mitchdraft) move this to a util file
func Contains(a []string, s string) bool {
	for _, n := range a {
		if s == n {
			return true
		}
	}
	return false
}
