package css

import (
	"bytes"
	"path/filepath"
	"strings"
)

type CSSBuild struct {
	*Basecss3Listener
	ClassName  map[string]struct{}
	SourceFile string
}

var KnownClasses = make(map[string]struct{})

var _ css3Listener = &CSSBuild{}

func NewCSSBuild(source string) *CSSBuild {
	return &CSSBuild{
		Basecss3Listener: &Basecss3Listener{},
		ClassName:        make(map[string]struct{}),
		SourceFile:       source,
	}
}

func (b *CSSBuild) RelativePath(path string) string {
	dir := filepath.Dir(b.SourceFile)
	result := filepath.Join(dir, path)
	return filepath.Clean(result)
}

func (b *CSSBuild) EnterKnownRuleset(ctx *KnownRulesetContext) {

}

func (b *CSSBuild) ExitKnownRuleset(ctx *KnownRulesetContext) {
	//x := ctx.SelectorGroup().GetText()
}

func (b *CSSBuild) EnterSelector(ctx *SelectorContext) {

}

func (b *CSSBuild) ExitSelector(ctx *SelectorContext) {
	raw := ctx.AllSimpleSelectorSequence()
	buf := &bytes.Buffer{}
	for i := 0; i < len(raw); i++ {
		buf.WriteString(" " + raw[i].GetText())
	}
	s := strings.TrimSpace(buf.String())
	b.ClassName[s] = struct{}{}
	KnownClasses[s] = struct{}{}
}
