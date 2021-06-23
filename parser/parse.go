package parser

func Parse(content string) error {
	return nil
}

// expression     → literal
//                | unary
//                | binary
//                | grouping ;

// literal        → NUMBER | STRING | "true" | "false" | "nil" ;
// grouping       → "(" expression ")" ;
// unary          → ( "-" | "!" ) expression ;
// binary         → expression operator expression ;
// operator       → "==" | "!=" | "<" | "<=" | ">" | ">="
//                | "+"  | "-"  | "*" | "/" ;

// All expression nodes implement the Expr interface.

type Node interface{}
type Expr interface {
	Node
	exprNode()
}

type BinaryExpr struct {
	Left, Right Expr
	Operator    Token
}

func (BinaryExpr) exprNode() {}
