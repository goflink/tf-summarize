package parser

import (
	"strings"

	"github.com/dineshba/terraform-plan-summary/reader"
	"github.com/dineshba/terraform-plan-summary/terraform_state"
)

type Parser interface {
	Parse() (terraform_state.TerraformState, error)
}

func CreateParser(data []byte, fileName string) (Parser, error) {
	if fileName != reader.StdinFileName && !strings.HasSuffix(fileName, ".json") {
		return NewBinaryParser(fileName), nil
	}
	return NewJsonParser(data), nil
}
