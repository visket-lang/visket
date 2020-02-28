package ast

import (
	"bytes"
	"fmt"
)

func Show(node Node) string {
	switch node := node.(type) {
	case *Program:
		var b bytes.Buffer
		for _, stmt := range node.Structs {
			b.WriteString(Show(stmt))
		}
		for _, stmt := range node.Functions {
			b.WriteString(Show(stmt))
		}
		return b.String()
	case *Identifier:
		return node.Token.Literal
	case *IntegerLiteral:
		return node.Token.Literal
	case *FloatLiteral:
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
		return fmt.Sprintf("(func-call %s(%s))", Show(node.Function), b.String())
	case *IndexExpression:
		return fmt.Sprintf("(%s[%s])", Show(node.Left), Show(node.Index))
	case *NewExpression:
		return fmt.Sprintf("(new %s)", node.Ident.Token.Literal)
	case *LoadMemberExpression:
		return fmt.Sprintf("(%s.%s)", Show(node.Left), node.MemberIdent.Token.Literal)
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
		for i, p := range node.Sig.Params {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(Show(p))
		}
		return fmt.Sprintf("(def-func %s(%s): %s (%s))", Show(node.Sig.Ident), b.String(), node.Sig.RetType.Token.Literal, Show(node.Body))
	case *Param:
		return fmt.Sprintf("%s: %s", Show(node.Ident), Show(node.Type))
	case *VarStatement:
		var b bytes.Buffer
		b.WriteString("(var ")
		b.WriteString(Show(node.Ident))
		if node.Type != nil {
			b.WriteString(": ")
			b.WriteString(node.Type.Token.Literal)
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
		b.WriteString("))")
		return b.String()
	case *ForStatement:
		var b bytes.Buffer
		b.WriteString("(for ")
		b.WriteString(Show(node.Init))
		b.WriteString("; ")
		b.WriteString(Show(node.Condition))
		b.WriteString("; ")
		b.WriteString(Show(node.Post))
		b.WriteString("(")
		b.WriteString(Show(node.Body))
		b.WriteString("))")
		return b.String()
	case *StructStatement:
		var b bytes.Buffer
		b.WriteString("(struct ")
		b.WriteString(Show(node.Ident))
		b.WriteString("(")
		for i, m := range node.Members {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(fmt.Sprintf("%s: %s", Show(m.Ident), Show(m.Type)))
		}
		b.WriteString("))")
		return b.String()
	case *Type:
		var buf bytes.Buffer

		if node.IsArray {
			buf.WriteString(fmt.Sprintf("[%d]", node.Len))
		}

		buf.WriteString(node.Token.Literal)
		return buf.String()
	}
	return fmt.Sprintf("unknown: %s", node)
}
