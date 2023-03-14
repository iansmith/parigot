package tree

type ImportSectionNode struct {
	TextItem_ []TextItem
}

func NewImportSectionNode() *ImportSectionNode {
	return &ImportSectionNode{}
}
