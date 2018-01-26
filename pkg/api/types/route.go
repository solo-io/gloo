package types

type Route struct {
	From Matcher
	To   isDestination
	// plugin content is unknown to the core api
	// plugins are passed
	Plugins []Spec
}

type Matcher struct {
	Path    isPathMatcher
	Headers map[string]string
	Verbs   []string
}

type isDestination interface {
	isDestination()
}

type FunctionDestination struct {
	FunctionName string
}

func (d FunctionDestination) isDestination() {}

type UpstreamDestination struct {
	UpstreamName string
	// If rewriteprefix is an empty string,
	// the incoming prefix will be preserved
	// Otherwise, it wil be replaced with RewritePrefix
	RewritePrefix string
}

func (d UpstreamDestination) isDestination() {}

type isPathMatcher interface {
	isPathMatcher()
}

type ExactPathMatcher struct {
	Path string
}

func (p ExactPathMatcher) isPathMatcher() {}

type PrefixPathMatcher struct {
	Prefix string
}

func (p PrefixPathMatcher) isPathMatcher() {}

type RegexPathMatcher struct {
	Regex string
}

func (p RegexPathMatcher) isPathMatcher() {}
