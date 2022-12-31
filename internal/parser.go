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
	NODE_ENDFOR
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
	p.parseNode(p.root, true)
}

func (p *Parser) parseNode(node ast, isRoot bool) {

	// On input if the next token

	for {
		token := p.lex.NextToken()
		switch token.Type {
		case ENDBLOCK:
			if isRoot {
				p.addError(token, fmt.Sprintf("Unexpected %s found", token.Type))
			}
			return
		case FUNCTS:
			if !isRoot {
				p.addError(token, "functions are only allowed at root")
			}
			p.processFunction(node, token)
		case EOF:
			if !isRoot {
				p.addError(token, "Unexpected EOF, was there a missing @} ?")
			}
			return
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
				p.processInclude(node, token)
			}
			// node.addChild(newAst(node, NODE_INCLUDE, p.processInclude(node, token)))
		case STARTBLOCK:
			p.processCodeBlock(node, token)
		case YIELD:
			if p.validateNoNewline(token) {
				node.addChild(newAst(node, NODE_YIELD, token))
			}
		case IF:
			if p.validateNoNewline(token) {
				node.addChild(newAst(node, NODE_IF, token))
			}
		case FOR:
			if p.validateNoNewline(token) {
				if p.validateFor(token) {
					node.addChild(newAst(node, NODE_FOR, token))
				}
			}
		case ELSE:
			node.addChild(newAst(node, NODE_ELSE, token))
		case ENDIF:
			node.addChild(newAst(node, NODE_ENDIF, token))
		case ENDFOR:
			node.addChild(newAst(node, NODE_ENDIF, token))
		case END:
			node.addChild(newAst(node, NODE_END, token))

		default:
			p.addError(token, fmt.Sprintf("Parser error Unexpected: %s:%s", token.Type, token.Literal))
		}
	}

}

// processCodeBlock
// Code blocks are placed inline with literal output
func (p *Parser) processCodeBlock(parent ast, token *Token) {
	child := newAst(parent, NODE_CODEBLOCK, token)
	parent.addChild(child)

	for {
		token := p.lex.NextToken()
		switch token.Type {

		case LITERAL:
			child.addChild(newAst(child, NODE_TOKEN_RAW, token))

		case ENDBLOCK:
			return

		case EOF:
			p.addError(token, "missing end of code block, unexpected EOF")
			return

		default:
			p.addError(token, "Unexpected inside @{")

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

		case ENDBLOCK:
			return

		case EOF:
			p.addError(token, "missing end of @func, unexpected EOF")
			return

		default:
			p.addError(token, "Unexpected inside @func")

		}
	}

}

// processInclude
// This node type has content, but must only be @content entries.
// Non-blank literals are ignored
func (p *Parser) processInclude(parent ast, token *Token) {
	// Process until end
	// node.addChild(newAst(node, NODE_INCLUDE, p.processInclude(node, token)))
	child := newAst(parent, NODE_INCLUDE, token)
	parent.addChild(child)

	for {
		token := p.lex.NextToken()
		switch token.Type {

		case LITERAL:
			if p.IsLiteralWhiteSpace(token.Literal) {
				// ok
			} else {
				p.addError(token, "Include blocks content must be embedded in a content block (@content)")
			}

		case CONTENT:
			p.processContent(child, token)
		case END:
			return

		case EOF:
			p.addError(token, "Expected end of include, unexpected EOF")
			return
		}
	}

}

func (p *Parser) processContent(parent ast, token *Token) {
	// Process until end
	// node.addChild(newAst(node, NODE_INCLUDE, p.processInclude(node, token)))
	child := newAst(parent, NODE_CONTENT, token)
	parent.addChild(child)

	p.parseNode(child, false)

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

func (p *Parser) validateFor(token *Token) bool {
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
