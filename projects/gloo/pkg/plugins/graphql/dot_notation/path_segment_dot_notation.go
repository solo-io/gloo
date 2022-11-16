package dot_notation

import (
	"regexp"
	"strconv"

	"github.com/rotisserie/eris"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
)

var (
	// pattern to match basic accessor
	Accessor = regexp.MustCompile(`^[a-zA-Z_\-$:][a-zA-Z0-9_\-$:]*$`)
	// pattern to match array index
	Index = regexp.MustCompile(`^\[([0-9]+|\*)]$`)
	// valid segment openers
	Opener = regexp.MustCompile(`^(?:[0-9]|\*)$`)
	// valid properties in [] format
	Property = regexp.MustCompile(`^\[(?:'|")(.*)(?:'|")]$`)
	// an entire segment; kinda the culmination of all of the above patterns
	Segment = regexp.MustCompile(`^((?:[a-zA-Z_\-$:][a-zA-Z0-9_\-$:]*)|(?:\[(?:'.*?'|".*?")])|(?:\[\d+])|(?:\[\*]))`)
)

func DotNotationToPathSegments(input string) ([]*v2.PathSegment, error) {
	if input == "" {
		return nil, nil
	}

	var keys []*v2.PathSegment
	var position int
	var unparsed = input

	for len(unparsed) > 0 {
		m := Segment.FindStringSubmatch(unparsed)
		if len(m) < 1 {
			return nil, eris.Errorf("Unexpected char: %c, index: %d, key: %s", unparsed[0], position, input)
		}

		prop := m[1]
		var val *v2.PathSegment
		var matches []string
		if Accessor.MatchString(prop) {
			val = &v2.PathSegment{Segment: &v2.PathSegment_Key{Key: prop}}
		} else if matches = Index.FindStringSubmatch(prop); len(matches) > 0 {
			if matches[1] == "*" {
				val = &v2.PathSegment{Segment: &v2.PathSegment_All{All: true}}
			} else {
				i, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, eris.Wrapf(err, "not a int %d", i)
				}
				val = &v2.PathSegment{Segment: &v2.PathSegment_Index{Index: uint32(i)}}
			}
		} else {
			val = &v2.PathSegment{Segment: &v2.PathSegment_Key{Key: Property.FindString(prop)}}
		}

		keys = append(keys, val)
		var remainder string

		if len(unparsed) == len(prop) {
			remainder = ""
		} else {
			remainder = unparsed[len(prop):]
			var isDot = remainder[0] == '.'
			if len(remainder) <= 1 {
				var f = "bracket"
				if isDot {
					f = "dot"
				}
				return nil, eris.New("Unable to parse '" + input + "' due to trailing " + f + "!")
			}
			var nextChar = string(remainder[1])

			var pattern = Opener
			if isDot {
				pattern = Accessor
			}
			if !pattern.MatchString(nextChar) {
				return nil, eris.Errorf("Unexpected char: %s, index: %d, key: %s", nextChar, position+len(prop)+1, input)
			}
			if isDot {
				remainder = remainder[1:]
			}

		}
		position += len(unparsed) - len(remainder)
		unparsed = remainder

	}
	return keys, nil
}
