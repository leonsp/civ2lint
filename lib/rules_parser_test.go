package lib

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestRulesParser_Parse(t *testing.T) {
	g := NewWithT(t)
	content := "@CIVILIZE\nExtra Advance 6,3,0,X7,no,0,0; X6\nExtra Advance 7,3,0,X7,no,0,0; X7"
	parser := &RulesParser{}
	err := parser.Parse(content)
	civilize := parser.Sections["@CIVILIZE"]

	g.Expect(err).To(BeNil())
	g.Expect(civilize).To(HaveLen(2))
	g.Expect(civilize[0].Text).To(ContainSubstring("Extra Advance 6"))
}
