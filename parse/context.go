package parse

// ContextType categorizes Context.
type ContextType int

// ContextType values.
const (
	CommandContext ContextType = iota
	ArgContext
	RedirFilenameContext
)

// Context contains information from the AST useful for tab completion.
type Context struct {
	Typ         ContextType
	CommandTerm *TermNode
	PrevTerms   *TermListNode
	PrevFactors *TermNode
	ThisFactor  *FactorNode
}

type PlainContext struct {
	Typ         ContextType
	CommandTerm string
	PrevTerms   []string
	PrevFactors string
	ThisFactor  *FactorNode
}

type notPlainFactor struct{}

func evalPlainFactor(fn *FactorNode) string {
	if fn == nil {
		return ""
	}
	if fn.Typ != StringFactor {
		panic(notPlainFactor{})
	}
	return fn.Node.(*StringNode).Text
}

func evalPlainTerm(tn *TermNode) (word string) {
	if tn == nil {
		return
	}
	for _, n := range tn.Nodes {
		word += evalPlainFactor(n)
	}
	return
}

func evalPlainTermList(tn *TermListNode) (words []string) {
	if tn == nil {
		return
	}
	for _, n := range tn.Nodes {
		words = append(words, evalPlainTerm(n))
	}
	return
}

func (c *Context) EvalPlain() (pctx *PlainContext) {
	defer func() {
		r := recover()
		if _, ok := r.(notPlainFactor); ok {
			pctx = nil
		} else if r != nil {
			panic(r)
		}
	}()
	return &PlainContext{
		Typ:         c.Typ,
		CommandTerm: evalPlainTerm(c.CommandTerm),
		PrevTerms:   evalPlainTermList(c.PrevTerms),
		PrevFactors: evalPlainTerm(c.PrevFactors),
		ThisFactor:  c.ThisFactor,
	}
}
