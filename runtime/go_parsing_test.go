package runtime

import (
	"testing"
	"text/template/parse"
)

func TestExtractVariablesFromTemplate(t *testing.T) {
	templateStr := `SELECT * FROM {{.env.partner_table_name}} WITH SAMPLING {{.env.partner_table_name}}% .... {{.user.domain}}`

	expected := []string{"env.partner_table_name", "env.partner_table_name", "user.domain"}
	result := extractVariablesFromTemplate(templateStr)

	if len(result) != len(expected) {
		t.Fatalf("Expected %d variables, but got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected variable '%s', but got '%s'", v, result[i])
		}
	}
}

func extractVariablesFromTemplate(templateStr string) []string {
	tree, err := parse.Parse("templateName", templateStr, "{{", "}}")
	if err != nil {
		return nil
	}

	var variables []string
	for _, t := range tree {
		walkNodes(t.Root, func(n parse.Node) {
			if vn, ok := n.(*parse.FieldNode); ok {
				variables = append(variables, joinIdentifiers(vn.Ident))
			}
		})
	}

	return variables
}

func walkNodes(node parse.Node, fn func(n parse.Node)) {
	fn(node)
	switch n := node.(type) {
	case *parse.ListNode:
		for _, ln := range n.Nodes {
			walkNodes(ln, fn)
		}
	case *parse.ActionNode:
		walkNodes(n.Pipe, fn)
	case *parse.PipeNode:
		for _, cmd := range n.Cmds {
			walkNodes(cmd, fn)
		}
	case *parse.CommandNode:
		for _, arg := range n.Args {
			walkNodes(arg, fn)
		}
	}
}

func joinIdentifiers(ident []string) string {
	var result string
	for _, id := range ident {
		if result != "" {
			result += "."
		}
		result += id
	}
	return result
}
