package routing

// Rule describes how dialed numbers are routed (extension, trunk, queue).
type Rule struct {
	Pattern string
	Target  string
	Type    string // extension | trunk | queue | ivr
}

// Engine evaluates dialplan rules.
type Engine struct {
	rules []Rule
}

func NewEngine(rules ...Rule) *Engine {
	return &Engine{rules: rules}
}

func (e *Engine) Match(destination string) (Rule, bool) {
	for _, rule := range e.rules {
		if rule.Pattern == destination {
			return rule, true
		}
	}
	return Rule{}, false
}
