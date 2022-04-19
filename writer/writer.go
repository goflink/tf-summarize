package writer

import (
	"github.com/dineshba/terraform-plan-summary/terraform_state"
	"io"
)

type Writer interface {
	Write(writer io.Writer) error
}

func CreateWriter(tree, separateTree, drawable bool, terraformState terraform_state.TerraformState) Writer {
	if tree {
		return NewTreeWriter(terraformState.ResourceChanges, drawable)
	}
	if separateTree {
		return NewSeparateTree(terraformState.AllChanges(), drawable)
	}
	return NewTableWriter(terraformState.AllChanges())
}
