package internal

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
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
	NODE_TOKEN_RAW
	NODE_YIELD
	NODE_IF
	NODE_FOR
	NODE_ELSE
	NODE_ENDIF
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
	depth        int
	includeDepth int
	lex          *Lexer
	root         ast
	current      ast
	errors       []PError
	imports      []*Token
	args         []*Token
	context      []*Token
	functions    []ast
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
		lex:          lex,
		imports:      make([]*Token, 0),
		args:         make([]*Token, 0),
		context:      make([]*Token, 0),
		depth:        0,
		includeDepth: 0,
		functions:    make([]ast, 0),
		errors:       make([]PError, 0),
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
		case ATDisplayUnsafe:
			node.addChild(newAst(node, NODE_DISPLAY_RAW, token))
		case IMPORT:
			if isRoot {
				p.imports = append(p.imports, token)
			} else {
				p.rootRequiredError(token)
			}
		case INCLUDE:
			node.addChild(newAst(node, NODE_INCLUDE_SIMPLE, token))
		case EXTEND:
			p.processInclude(node, token)
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
		case ENDBLOCK:
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

func (p *Parser) wStr(o io.Writer, msg string) *Parser {
	o.Write([]byte(msg))
	return p
}
func (p *Parser) wNL(o io.Writer) *Parser {
	o.Write([]byte("\n"))
	return p
}

// renderOutput
func (p *Parser) renderOutput(o io.Writer, packageName string, templateName string, opt *BlipOptions) {
	p.wStr(o, "package ").wStr(o, packageName).wNL(o)

	p.wStr(o, "// Do Not Edit\n")
	p.wStr(o, fmt.Sprintf("// Generated by Blip %s\n", Version))
	p.wStr(o, fmt.Sprintf("// %v\n", time.Now()))

	p.outputImports(o, opt)
	p.writeFuncts(o)

	p.writeMainFunction(o, templateName)
}

func (p *Parser) outputImports(o io.Writer, opt *BlipOptions) {

	imports := make(map[string]string)

	if !p.hasErrors() {
		imports["\""+opt.SupportBranch+"\""] = ""
		// imports["\"github.com/samlotti/blip/blipUtil\""] = ""
	}
	imports["\"context\""] = ""
	imports["\"io\""] = ""
	imports["\"fmt\""] = ""

	if !p.hasErrors() {
		for _, imp := range p.imports {
			trimmed := strings.Trim(imp.Literal, "\t")
			trimmed = strings.Trim(trimmed, " ")
			imports[strings.Trim(trimmed, " ")] = ""
		}
	}

	imports2 := make([]string, 0)
	for k, _ := range imports {
		imports2 = append(imports2, k)
	}
	sort.Strings(imports2)
	p.wStr(o, "\nimport (")
	for _, v := range imports2 {
		p.wStr(o, "\n\t").wStr(o, v)
	}
	p.wStr(o, "\n)\n\n")
}

func (p *Parser) writeFuncts(o io.Writer) {
	if p.hasErrors() {
		return
	}
	for _, fa := range p.functions {
		for _, ft := range fa.GetChildren() {
			p.wStr(o, fmt.Sprintf("// Function block from line: %d\n", ft.GetToken().Line))
			p.wStr(o, ft.GetToken().Literal)
		}

	}
}

func (p *Parser) hasErrors() bool {
	return len(p.errors) > 0
}

// convertTemplateNameToFunctionName
// convert   path.path.templateName
// tp        path.path.TemplateNameProcess
func (p *Parser) convertTemplateNameToFunctionName(templateName string) string {
	sects := strings.Split(templateName, ".")
	last := len(sects) - 1
	sects[last] = strings.Title(sects[last]) + "Process"
	return strings.Join(sects, ".")

}

func (p *Parser) writeMainFunction(o io.Writer, templateName string) {
	p.wStr(o, "\n\n").
		wStr(o, "func ").
		wStr(o, p.convertTemplateNameToFunctionName(templateName)).
		wStr(o, "( ")

	if !p.hasErrors() {
		p.writeArgVar(o)
	}

	p.wStr(o, "c context.Context, w io.Writer ) ")
	p.wStr(o, `(terror error) {
	var si = blipUtil.Instance()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Catch panic %s: %s\n", "`)
	p.wStr(o, p.convertTemplateNameToFunctionName(templateName))
	p.wStr(o, `", err)
			terror = fmt.Errorf("%v", err)
		}
	}()
	si.IncProcess()
`)

	if p.hasErrors() {
		p.wStr(o, "Errors found in transforming the template\n")
		for _, err := range p.errors {
			p.wStr(o, fmt.Sprintf("Error at %d: %s\n", err.lineNum, err.msg))
		}
		p.wStr(o, "\n")
	} else {
		p.writeContentVar(o)
		p.writeBody(p.root, 1, o)
	}

	p.wStr(o, "\treturn")
	p.wStr(o, "\n}")

}

func (p *Parser) addSlashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\n':
			buf.WriteString("\\n")
		case '\t':
			buf.WriteString("\\t")
		case '\'':
			buf.WriteRune('\\')
			buf.WriteRune(char)
		case '"':
			buf.WriteRune('\\')
			buf.WriteRune(char)
		default:
			buf.WriteRune(char)
		}

	}
	return buf.String()
}

func (p *Parser) getTabsDepth(depth int) string {
	var tabs = "\t\t\t\t\t\t\t\t"[0:depth]
	return tabs
}

func (p *Parser) writeBody(node ast, depth int, o io.Writer) {

	for _, ast := range node.GetChildren() {

		base := ast.(*astBase)
		var tabs = p.getTabsDepth(depth)

		p.wStr(o, fmt.Sprintf("%s// Line: %d\n", tabs, base.token.Line))

		switch base.nodeType {
		case NODE_TOKEN_RAW:
			p.wStr(o, base.token.Literal)
		case NODE_TOKEN:
			p.wStr(o, fmt.Sprintf("%ssi.Write(w, []byte(\"%s\"))\n", tabs, p.addSlashes(base.token.Literal)))
		case NODE_DISPLAY:
			// si.WriteStr(w, game.Opponent)
			p.wStr(o, fmt.Sprintf("%ssi.WriteStrSafe(w, %s)\n", tabs, p.addSlashes(base.token.Literal)))
		// f.Sp "si.Write(w, indexpage1)")
		case NODE_DISPLAY_RAW:
			// si.WriteStr(w, game.Opponent)
			p.wStr(o, fmt.Sprintf("%ssi.WriteStr(w, %s)\n", tabs, p.addSlashes(base.token.Literal)))
		// f.Sp "si.Write(w, indexpage1)")
		case NODE_INCLUDE_SIMPLE:
			p.WriteNodeSimpleCall(o, base, depth, "c")
			// @include Base @
		case NODE_INCLUDE:
			p.includeDepth += 1
			// var c1 context.Context
			p.wStr(o, fmt.Sprintf("%svar %s = context.WithValue(c, \"__Blip__\", 1)\n", tabs, p.contextVarName()))
			// Do After
		case NODE_CODEBLOCK:
			p.wStr(o, fmt.Sprintf("%s// Code block follows\n", tabs))

		case NODE_CONTENT:
			// var f1 = func() {
			p.wStr(o, fmt.Sprintf("%svar contentF%d = func() (terror error) {\n", tabs, p.includeDepth))

		case NODE_YIELD:
			// 	si.CallCtxFunc(c, "myJavascript")
			p.wStr(o, fmt.Sprintf("%sterror = si.CallCtxFunc(c, \"%s\")\n", tabs, p.trimAll(base.token.Literal)))
			p.wStr(o, fmt.Sprintf("%sif terror != nil { return }\n", tabs))

		case NODE_IF:
			p.wStr(o, fmt.Sprintf("%sif %s {\n", tabs, p.trimAll(base.token.Literal)))
		case NODE_FOR:
			p.WriteNodeForStatement(o, base, depth)
		case NODE_ELSE:
			p.wStr(o, fmt.Sprintf("%s} else {\n", tabs))
		case NODE_ENDIF:
			p.wStr(o, fmt.Sprintf("%s}\n", tabs))
		case NODE_ENDFOR:
			p.wStr(o, fmt.Sprintf("%s}\n", tabs))

		default:
			// Force compiler error.. should never happen unless I add new nodes and forget to implement
			p.wStr(o, fmt.Sprintf("%s Unsupported Node type: %d: '%s'  -> %d\n", tabs, base.nodeType, p.addSlashes(base.token.Literal), len(base.children)))

		}

		p.writeBody(base, depth+1, o)

		switch base.nodeType {
		case NODE_INCLUDE:
			p.WriteNodeSimpleCall(o, base, depth, p.contextVarName())
			p.includeDepth -= 1
		case NODE_CONTENT:
			// var f1 = func() {
			p.wStr(o, fmt.Sprintf("%sreturn\n", p.getTabsDepth(depth+1)))
			p.wStr(o, fmt.Sprintf("%s}\n", tabs))
			p.wStr(o, fmt.Sprintf("%s%s = context.WithValue(%s, ", tabs, p.contextVarName(), p.contextVarName()))
			p.wStr(o, fmt.Sprintf("\"%s\", contentF%d)", p.trimAll(base.token.Literal), p.includeDepth))

		}

	}
}

func (p *Parser) writeArgVar(o io.Writer) {
	for idx, c := range p.args {
		if idx > 0 {
			p.wStr(o, ", ")
		}
		// 	user := c.Value("user").(manual.User)
		split := p.splitString(c.Literal, 2)
		p.wStr(o, split[0])
		p.wStr(o, " ")
		p.wStr(o, split[1])
	}

	if len(p.args) > 0 {
		p.wStr(o, ", ")
	}
}

func (p *Parser) writeContentVar(o io.Writer) {

	for _, c := range p.context {

		eq := " := "
		var addNullCheck = false
		if strings.Contains(c.Literal, "=") {
			eq = " = "
			// ex: @context errors []string = make([]string,0)
			p.wStr(o, "\tvar ")
			p.wStr(o, c.Literal)
			p.wStr(o, "\n")
			addNullCheck = true

		}
		// 	user := c.Value("user").(manual.User)
		split := strings.Split(strings.Trim(c.Literal, " "), " ")

		if addNullCheck {
			// if c.Value("errors") != nil {
			p.wStr(o, "\tif c.Value(\"")
			p.wStr(o, split[0])
			p.wStr(o, "\") != nil {")
			p.wNL(o)
			p.wStr(o, "\t")
		}

		p.wStr(o, "\t")
		p.wStr(o, split[0])
		p.wStr(o, eq)
		p.wStr(o, "c.Value(\"")
		p.wStr(o, split[0])
		p.wStr(o, "\").(")
		p.wStr(o, split[1])
		p.wStr(o, ")")
		p.wNL(o)

		if addNullCheck {
			// if c.Value("errors") != nil {
			p.wStr(o, "\t}")
			p.wNL(o)
		}

	}
}

func (p *Parser) verifySplit(token *Token, count int) {
	split := strings.SplitN(strings.Trim(token.Literal, " "), " ", count)
	if len(split) != count {
		p.addError(token, fmt.Sprintf("Expected %d values split by spaces in %s", count, token.Literal))
	}
}

func (p *Parser) splitString(lit string, max int) []string {
	return strings.SplitN(strings.Trim(lit, " "), " ", max)
}

func (p *Parser) WriteNodeSimpleCall(o io.Writer, node ast, depth int, contextName string) {
	// @include Base @
	p.wNL(o)

	splits := p.splitString(node.GetToken().Literal, 2)
	funName := p.convertTemplateNameToFunctionName(splits[0])
	p.wStr(o, p.getTabsDepth(depth))
	p.wStr(o, "terror = ")
	p.wStr(o, funName).wStr(o, "(")
	for idx, s := range splits {
		if idx == 0 {
			continue
		}
		if idx > 1 {
			p.wStr(o, ", ")
		}
		p.wStr(o, s)
	}
	if len(splits) > 1 {
		p.wStr(o, ", ")
	}
	p.wStr(o, fmt.Sprintf("%s, w)\n", contextName))

	p.wStr(o, p.getTabsDepth(depth))
	p.wStr(o, "if terror != nil { return }\n")

}

func (p *Parser) contextVarName() string {
	return fmt.Sprintf("ctxL%d", p.includeDepth)
}

func (p *Parser) trimAll(literal string) string {
	return strings.Trim(literal, " ")

}

// validateNoNewline
func (p *Parser) validateNoNewline(token *Token) bool {
	if strings.Contains(token.Literal, "\n") {
		p.addError(token, "It looks like this command was not closed properly")
		return false
	}
	return true

}

func (p *Parser) splitFor(token *Token) []string {
	return strings.Split(strings.Trim(token.Literal, " "), " ")
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

func (p *Parser) WriteNodeForStatement(o io.Writer, base *astBase, depth int) {
	sects := p.splitFor(base.token)
	p.wStr(o, fmt.Sprintf("%sfor idx, %s := range %s { _ = idx\n", p.getTabsDepth(depth), sects[0], sects[2]))

}
