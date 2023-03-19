package tree

type ImportSectionNode struct {
	TextItem_                []TextItem
	LineNumber, ColumnNumber int
}

func NewImportSectionNode() *ImportSectionNode {
	return &ImportSectionNode{}
}
