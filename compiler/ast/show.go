package ast

import (
	"bytes"
	"fmt"
)

func Show(node Node) string {
	switch node := node.(type) {
	case *Program:
		var b bytes.Buffer
		for _, stmt := range node.Functions {
			b.WriteString(Show(stmt))
		}
		return b.String()
	case *Identifier:
		return node.Token.Literal
	case *IntegerLiteral:
		return node.Token.Literal
	case *PrefixExpression:
		return fmt.Sprintf("(%s%s)", node.Operator, Show(node.Right))
	case *InfixExpression:
		return fmt.Sprintf("(%s %s %s)", Show(node.Left), node.Operator, Show(node.Right))
	case *AssignExpression:
		return fmt.Sprintf("(%s = %s)", Show(node.Left), Show(node.Value))
	case *CallExpression:
		var b bytes.Buffer
		for i, param := range node.Parameters {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(Show(param))
		}
		return fmt.Sprintf("(func-call %s[%s])", Show(node.Function), b.String())
	case *BlockStatement:
		var b bytes.Buffer
		for _, stmt := range node.Statements {
			b.WriteString(Show(stmt))
		}
		return b.String()
	case *ExpressionStatement:
		return Show(node.Expression)
	case *FunctionStatement:
		var b bytes.Buffer
		for i, param := range node.Parameters {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(Show(param))
			b.WriteString(": ")
			b.WriteString(node.Type.Params[i].String())
		}
		return fmt.Sprintf("(def-func %s[%s]: %s (%s))", Show(node.Ident), b.String(), node.Type.RetType, Show(node.Body))
	case *VarStatement:
		var b bytes.Buffer
		b.WriteString("(var ")
		b.WriteString(Show(node.Ident))
		if node.Type != nil {
			b.WriteString(": ")
			b.WriteString(node.Type.String())
		}
		if node.Value != nil {
			b.WriteString(" = ")
			b.WriteString(Show(node.Value))
		}
		b.WriteString(")")
		return b.String()
	case *ReturnStatement:
		if node.Value == nil {
			return "(return)"
		}
		return fmt.Sprintf("(return %s)", Show(node.Value))
	case *IfStatement:
		var b bytes.Buffer
		b.WriteString("(if ")
		b.WriteString(Show(node.Condition))
		b.WriteString("(")
		b.WriteString(Show(node.Consequence))
		b.WriteString(")")
		if node.Alternative != nil {
			b.WriteString("(")
			b.WriteString(Show(node.Alternative))
			b.WriteString(")")
		}
		b.WriteString(")")
		return b.String()
	case *WhileStatement:
		var b bytes.Buffer
		b.WriteString("(while ")
		b.WriteString(Show(node.Condition))
		b.WriteString("(")
		b.WriteString(Show(node.Body))
		b.WriteString(")")
		return b.String()
	case *ForStatement:
		var b bytes.Buffer
		b.WriteString("(for ")
		b.WriteString(Show(node.Init))
		b.WriteString("; ")
		b.WriteString(Show(node.Init))
		b.WriteString("; ")
		b.WriteString(Show(node.Condition))
		b.WriteString("; ")
		b.WriteString(Show(node.Post))
		b.WriteString("(")
		b.WriteString(Show(node.Body))
		b.WriteString("))")
		return b.String()
	}
	return fmt.Sprintf("unknown: %s", node.String())
}
