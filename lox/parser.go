package lox

import "fmt"

type Parser struct {
	tokens  []Token
	current int

	hadError bool
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements, nil
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if val := recover(); val != nil {
			parsingError := val.(*LoxError)
			fmt.Println(parsingError.Error())
			p.synchronize()
			p.hadError = true
		}
	}()

	if p.match(CLASS) {
		return p.classDeclaration()
	}

	if p.match(FUN) {
		return p.function("function")
	}

	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) classDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "expect class name.")
	p.consume(LEFT_BRACE, "expect '{' before class body.")

	var methods []Function
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, *(p.function("method")).(*Function))
	}

	p.consume(RIGHT_BRACE, "expect '}' after class body.")
	return NewClass(name, Variable{}, methods)
}

func (p *Parser) function(kind string) Stmt {
	name := p.consume(IDENTIFIER, fmt.Sprintf("expect %v name.", kind))
	p.consume(LEFT_PAREN, fmt.Sprintf("expect '(' after %v name.", kind))
	var parameters []Token
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				p.reportError(p.peek(), "can't have more than 255 parameters.")
			}
			parameters = append(parameters, p.consume(IDENTIFIER, "expect parameter name."))
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "expect ')' after parameters.")

	p.consume(LEFT_BRACE, fmt.Sprintf("expect '{' before %v body.", kind))
	body := p.blockStatement()
	return NewFunction(name, parameters, body)
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "expect variable name.")

	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "expect ';' after variable declaration.")
	return NewVar(name, initializer)
}

func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	}

	if p.match(IF) {
		return p.ifStatement()
	}

	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(RETURN) {
		return p.returnStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	if p.match(LEFT_BRACE) {
		return NewBlock(p.blockStatement())
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "expect '(' after 'for'.")

	var initializer Stmt = nil
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr = nil
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "expect ';' after loop condition.")

	var increment Expr = nil
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = NewBlock([]Stmt{body, NewExpression(increment)})
	}

	if condition == nil {
		condition = NewLiteral(true)
	}
	body = NewWhile(condition, body)

	if initializer != nil {
		body = NewBlock([]Stmt{initializer, body})
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "expect '(' after if.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "expect ')' after condition.")

	thenBranch := p.statement()
	var elseBranch Stmt = nil
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return NewIf(condition, thenBranch, elseBranch)
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "expect '(' after while.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "expect ')' after condition.")
	body := p.statement()
	return NewWhile(condition, body)
}

func (p *Parser) blockStatement() []Stmt {
	var statements []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(RIGHT_BRACE, "expect '}' after block.")
	return statements
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "expect ';' after value.")
	return NewPrint(value)
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		value = p.expression()
	}
	p.consume(SEMICOLON, "expect ';' after return value.")
	return NewReturn(keyword, value)
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "expect ';' after expression.")
	return NewExpression(expr)
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if v, ok := expr.(*Variable); ok {
			name := v.Name
			return NewAssign(name, value)
		} else if g, ok := expr.(*Get); ok {
			return NewSet(g.Object, g.Name, value)
		}

		p.reportError(equals, "invalid assignment target")
	}

	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = NewLogical(expr, operator, right)
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = NewLogical(expr, operator, right)
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}
	return p.call()
}

func (p *Parser) finishCall(callee Expr) Expr {
	var arguments []Expr
	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				p.reportError(p.peek(), "can't have more than 255 arguments")
			}
			arguments = append(arguments, p.expression())
			if !p.match(COMMA) {
				break
			}
		}
	}

	paren := p.consume(RIGHT_PAREN, "expect ')' after arguments.")
	return NewCall(callee, paren, arguments)
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name := p.consume(IDENTIFIER, "expect property after '.'.")
			expr = NewGet(name, expr)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return NewLiteral(false)
	}
	if p.match(TRUE) {
		return NewLiteral(true)
	}
	if p.match(NIL) {
		return NewLiteral(nil)
	}
	if p.match(NUMBER, STRING) {
		return NewLiteral(p.previous().Literal)
	}
	if p.match(IDENTIFIER) {
		return NewVariable(p.previous())
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "expect ')' after expression.")
		return NewGrouping(expr)
	}
	panic(NewLoxError(p.peek(), "expected expression"))
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, t := range tokenTypes {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}
func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType TokenType, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic(NewLoxError(p.peek(), message))
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) reportError(token Token, message string) {
	fmt.Println(NewRuntimeError(token, message))
	p.hadError = true
}
