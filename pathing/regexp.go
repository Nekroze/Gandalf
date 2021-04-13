package pathing

import (
	"fmt"
	"regexp"

	"github.com/Nekroze/Gandalf/check"
)

type RegexpExtractor struct {
	Expression *regexp.Regexp
}

func (rc RegexpExtractor) Extractor(source, captureGroup string) (found []string, err error) {
	// TOOD: memoize this as it will be common to run it across the same source in rapid succession
	matches := rc.Expression.FindStringSubmatch(source)
	found = []string{}
	if len(matches) < 1 {
		return found, fmt.Errorf("failed to match regexp %#v", rc.Expression.String())
	}
	captureGroupIndex := rc.Expression.SubexpIndex(captureGroup)
	if captureGroupIndex == -1 {
		return found, fmt.Errorf("failed to match regexp %#v and extract capture group %#v", rc.Expression.String(), captureGroup)
	}
	if match := matches[captureGroupIndex]; match != "" {
		found = append(found, match)
	}
	return found, nil
}

func NewRegexpExtractor(expression string) RegexpExtractor {
	return RegexpExtractor{
		Expression: regexp.MustCompile(expression),
	}
}

func RegexpChecks(expression string, pcs PathChecks) check.Func {
	rc := NewRegexpExtractor(expression)
	return Checks(rc.Extractor, pcs)
}
