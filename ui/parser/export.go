package parser

// this is so we don't get into a filesystem-vs-caps battle
type WCLParser struct {
	*wcl
}

func WCLParserFromWcl(p *wcl) *WCLParser {
	return &WCLParser{
		wcl: p,
	}
}
