package internal

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	NODE_ROOT = iota
	NODE_TOKEN
	NODE_INCLUDE_SIMPLE
	NODE_INCLUDE
	NODE_CONTENT
	NODE_FUNC
	NODE_CODEBLOCK
	NODE_TEXT
	NODE_DISPLAY
	NODE_DISPLAY_RAW
	NODE_DISPLAY_INT
	NODE_DISPLAY_INT64
	NODE_TOKEN_RAW
	NODE_YIELD
	NODE_IF
	NODE_FOR
	NODE_ELSE
	NODE_ENDIF
	NODE_END
)

type ast interface {
	GetRoot() ast
	addChild(child ast)
	GetChildren() []ast
	ToString() string
	GetToken() *Token
}

type astBase struct {
	parent   ast
	children []ast
	nodeType int
	token    *Token
}

func (b *astBase) GetRoot() ast {
	return b.parent.GetRoot()
}

func (b *astBase) GetToken() *Token {
	return b.token
}

func (b *astBase) ToString() string {
	return fmt.Sprintf("%d: %d: %s:%s", b.nodeType, len(b.children), b.token.Type, b.token.Literal)
}

func (b *astBase) GetChildren() []ast {
	return b.children
}

func (b *astBase) addChild(child ast) {
	b.children = append(b.children, child)
}

type rootAst struct {
	astBase
}

func (r *rootAst) GetRoot() ast {
	return r
}

type PError struct {
	lineNum int
	linePos int
	msg     string
}

type Parser struct {
	depth int

	lex       *Lexer
	root      ast
	current   ast
	errors    []PError
	imports   []*Token
	args      []*Token
	context   []*Token
	functions []ast
}

func newAst(parent ast, nodeType int, token *Token) *astBase {
	return &astBase{
		parent:   parent,
		children: make([]ast, 0),
		nodeType: nodeType,
		token:    token,
	}
}

func New(lex *Lexer) *Parser {
	return &Parser{
		lex:       lex,
		imports:   make([]*Token, 0),
		args:      make([]*Token, 0),
		context:   make([]*Token, 0),
		depth:     0,
		functions: make([]ast, 0),
		errors:    make([]PError, 0),
		root: &rootAst{
			astBase: *newAst(nil, NODE_ROOT, nil),
		},
	}
}

// Parse
// At a top level node
func (p *Parser) Parse() {
	p.parseNode(p.root, true, []TokenType{EOF})
}

func (p *Parser) contains(s []TokenType, str TokenType) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// parseNode --
// Parse the content and returns when specif tokens are found.
// The main loop will continue until EOF
// A for loop will continue until EOF or END node found
func (p *Parser) parseNode(node ast, isRoot bool, terminators []TokenType) *Token {

	// On input is the next token

	for {
		token := p.lex.NextToken()

		if p.contains(terminators, token.Type) {
			return token
		}

		switch token.Type {
		//case ENDBLOCK:
		//	if isRoot {
		//		p.addError(token, fmt.Sprintf("Unexpected %s found", token.Type))
		//	}
		//	p.addError(token, fmt.Sprintf("Unexpected %s found", token.Type))
		//	return token
		case FUNCTS:
			if !isRoot {
				p.addError(token, "functions are only allowed at root")
			}
			p.processFunction(node, token)
		case EOF:
			if !isRoot {
				p.addError(token, "Unexpected EOF, was there a missing @end ?")
			}
			return token
		case LITERAL:
			node.addChild(newAst(node, NODE_TOKEN, token))

		case ARG:
			if isRoot {
				p.args = append(p.args, token)
				p.verifySplit(token, 2)
			} else {
				p.rootRequiredError(token)
			}
		case CONTEXT:
			if isRoot {
				p.context = append(p.context, token)
				p.verifySplit(token, 2)
			} else {
				p.rootRequiredError(token)
			}
		case ATDisplay:
			node.addChild(newAst(node, NODE_DISPLAY, token))
		case ATDisplayInt:
			node.addChild(newAst(node, NODE_DISPLAY_INT, token))
		case ATDisplayInt64:
			node.addChild(newAst(node, NODE_DISPLAY_INT64, token))
		case ATDisplayUnsafe:
			node.addChild(newAst(node, NODE_DISPLAY_RAW, token))
		case IMPORT:
			if isRoot {
				p.imports = append(p.imports, token)
			} else {
				p.rootRequiredError(token)
			}
		case INCLUDE:
			if p.validateNoNewline(token) {
				node.addChild(newAst(node, NODE_INCLUDE_SIMPLE, token))
			}
		case EXTEND:
			if p.validateNoNewline(token) {
				p.processExtend(node, token)
			}
			// node.addChild(newAst(node, NODE_INCLUDE, p.processExtend(node, token)))
		case TEXT:
			p.processTextBlock(node, token)
		case STARTBLOCK:
			p.processCodeBlock(node, token)
		case YIELD:
			if p.validateNoNewline(token) {
				node.addChild(newAst(node, NODE_YIELD, token))
			}
		case IF:
			p.processIfStatement(node, token)

		case FOR:
			p.processForBlock(node, token)
		//case ELSE:
		//	node.addChild(newAst(node, NODE_ELSE, token))
		//case ENDIF:
		//	node.addChild(newAst(node, NODE_ENDIF, token))
		//case END:
		//	//node.addChild(newAst(node, NODE_END, token))
		//	//if isRoot {
		//	p.addError(token, fmt.Sprintf("Unexpected %s found", token.Type))
		//	//}
		//	return token

		default:
			p.addError(token, fmt.Sprintf("Parser error Unexpected: %s:%s", token.Type, token.Literal))
		}
	}

}

// processIfStatement
// If Else End
func (p *Parser) processIfStatement(parent ast, token *Token) {
	if !p.validateNoNewline(token) {
		return
	}

	child := newAst(parent, NODE_IF, token)
	parent.addChild(child)

	endToken := p.parseNode(child, false, []TokenType{END, ELSE})
	if endToken.Type == ELSE {
		child.addChild(newAst(child, NODE_ELSE, endToken))
		endToken = p.parseNode(child, false, []TokenType{END})
	}
	if endToken.Type == END {
		child.addChild(newAst(child, NODE_ENDIF, endToken))
	} else {
		p.addError(token, fmt.Sprintf("Expected %s found %s at line %d, unterminated @for block ", END, endToken.Type, endToken.Line))
	}
}

// processForBlock
// @for x in list ...
// @end
func (p *Parser) processForBlock(parent ast, token *Token) {
	if !p.validateNoNewline(token) || !p.validateForStatementCommand(token) {
		return
	}
	child := newAst(parent, NODE_FOR, token)
	parent.addChild(child)
	endToken := p.parseNode(child, false, []TokenType{END})
	if endToken.Type != END {
		p.addError(token, fmt.Sprintf("Expected %s found %s at line %d, unterminated @for block ", END, endToken.Type, endToken.Line))
	} else {
		child.addChild(newAst(child, NODE_END, token))
	}
}

// processCodeBlock
// Code blocks are placed inline with literal output
func (p *Parser) processCodeBlock(parent ast, cbtoken *Token) {
	child := newAst(parent, NODE_CODEBLOCK, cbtoken)
	parent.addChild(child)
	//endToken := p.parseNode(child, false, []TokenType{ENDBLOCK}, true)
	//if endToken.Type != ENDBLOCK {
	//	p.addError(token, fmt.Sprintf("Expected %s found %s at line %d, unterminated @{ block ", ENDBLOCK, endToken.Type, endToken.Line))
	//}
	for {
		token := p.lex.NextToken()
		switch token.Type {
		case LITERAL:
			child.addChild(newAst(child, NODE_TOKEN_RAW, token))
		case END:
			return
		case EOF:
			// cbtoken to show line of statr code block
			p.addError(cbtoken, "missing end of code block, unexpected EOF")
			return
		default:
			p.addError(token, "Unexpected inside @code")

		}
	}
}

// processTextBlock
// Text blocks are raw output of the text.
func (p *Parser) processTextBlock(parent ast, cbtoken *Token) {
	child := newAst(parent, NODE_TEXT, cbtoken)
	parent.addChild(child)
	for {
		token := p.lex.NextToken()
		switch token.Type {
		case LITERAL:
			child.addChild(newAst(child, NODE_TOKEN, token))
		case END:
			return
		case EOF:
			p.addError(cbtoken, "missing end of text block, unexpected EOF")
			return
		default:
			p.addError(token, "Unexpected inside @text")

		}
	}

}

// processFunction
// Collection of literals until end bloc
// functions are placed outside the main code
func (p *Parser) processFunction(parent ast, token *Token) {
	child := newAst(parent, NODE_FUNC, token)
	// parent.addChild(child)
	p.functions = append(p.functions, child)

	for {
		token := p.lex.NextToken()
		switch token.Type {

		case LITERAL:
			child.addChild(newAst(child, NODE_TOKEN, token))

		case END:
			return

		case EOF:
			p.addError(token, "missing end of @func, unexpected EOF")
			return

		default:
			p.addError(token, "Unexpected inside @func")

		}
	}

}

// processExtend
// This node type has content, but must only be @content entries.
// Non-blank literals are ignored
func (p *Parser) processExtend(parent ast, token *Token) {
	// Process until end
	// node.addChild(newAst(node, NODE_INCLUDE, p.processExtend(node, token)))
	child := newAst(parent, NODE_INCLUDE, token)
	parent.addChild(child)

	trimErrors := false
	for {
		token2 := p.lex.NextToken()
		switch token2.Type {

		case LITERAL:
			if p.IsLiteralWhiteSpace(token2.Literal) {
				// ok
			} else {
				p.addError(token, "Include blocks content must be embedded in a content block (@content)")
			}

		case CONTENT:
			p.processContent(child, token2)

		case END:
			return

		case EOF:
			p.addError(token, "Expected end of include, unexpected EOF")
			return

		default:
			if !trimErrors {
				p.addError(token, fmt.Sprintf("Unexpected %s at line %d", token2.Type, token2.Line))
				trimErrors = true
			}

		}
	}

}

func (p *Parser) processContent(parent ast, token *Token) {
	// Process until end
	// node.addChild(newAst(node, NODE_INCLUDE, p.processExtend(node, token)))
	child := newAst(parent, NODE_CONTENT, token)
	parent.addChild(child)

	if strings.Contains(token.Literal, "@") {
		p.addError(token, fmt.Sprintf("@content has invalid characters, @ not expected"))
	}

	endNode := p.parseNode(child, false, []TokenType{END})
	if endNode.Type != END {
		p.addError(token, "Content not terminated, expected @end")
	}

}

func (p *Parser) rootRequiredError(token *Token) {
	p.errors = append(p.errors, PError{
		lineNum: token.Line,
		linePos: token.Pos,
		msg:     fmt.Sprintf("%s is only allowed at root level", token.Type),
	})
}

func (p *Parser) addError(token *Token, msg string) {
	p.errors = append(p.errors, PError{
		lineNum: token.Line,
		linePos: token.Pos,
		msg:     fmt.Sprintf("%s : %s", token.Type, msg),
	})
}

func (p *Parser) IsLiteralWhiteSpace(literal string) bool {
	var pos = strings.IndexFunc(literal, func(r rune) bool {
		return !unicode.IsSpace(r)
	})
	return pos == -1
}

func (p *Parser) Dump() {
	fmt.Printf("================ Parse Results ============= \n")
	p.dumpNode(p.root, 0)
}

func (p *Parser) dumpNode(node ast, depth int) {
	if len(node.GetChildren()) == 0 {
		return
	}
	for idx, child := range node.GetChildren() {
		var strs = "                       "[0 : depth*2]
		fmt.Printf("%02d:%d: | %s%s\n", depth, idx, strs, child.ToString())
		p.dumpNode(child, depth+1)
	}

}

func (p *Parser) hasErrors() bool {
	return len(p.errors) > 0
}

func (p *Parser) verifySplit(token *Token, count int) {
	split := strings.SplitN(strings.Trim(token.Literal, " "), " ", count)
	if len(split) != count {
		p.addError(token, fmt.Sprintf("Expected %d values split by spaces in %s", count, token.Literal))
	}
}

// validateNoNewline
func (p *Parser) validateNoNewline(token *Token) bool {

	if strings.Contains(token.Literal, "@") {
		p.addError(token, fmt.Sprintf("Invalid character '@' found in %s", token.Literal))
		return false
	}

	if strings.Contains(token.Literal, "\n") {
		p.addError(token, "It looks like this command was not closed properly")
		return false
	}
	return true

}

func (p *Parser) validateForStatementCommand(token *Token) bool {
	sects := p.splitFor(token)
	if len(sects) != 3 {
		p.addError(token, fmt.Sprintf("@for expected:  `variable in list` found %d components, note remove any extra spaces", len(sects)))
		return false
	}
	if sects[1] != "in" {
		p.addError(token, fmt.Sprintf("@for expected:  `variable in list` found %s instead of  `in`", sects[1]))
		return false
	}
	return true
}

func (p *Parser) splitFor(token *Token) []string {
	return strings.Split(strings.Trim(token.Literal, " "), " ")
}
