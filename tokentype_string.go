// Code generated by "stringer -type=TokenType token_type.go"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LEFT_PAREN-0]
	_ = x[RIGHT_PAREN-1]
	_ = x[LEFT_BRACE-2]
	_ = x[RIGHT_BRACE-3]
	_ = x[COMMA-4]
	_ = x[DOT-5]
	_ = x[MINUS-6]
	_ = x[PLUS-7]
	_ = x[SEMICOLON-8]
	_ = x[SLASH-9]
	_ = x[STAR-10]
	_ = x[BANG-11]
	_ = x[BANG_EQUAL-12]
	_ = x[EQUAL-13]
	_ = x[EQUAL_EQUAL-14]
	_ = x[GREATER-15]
	_ = x[GREATER_EQUAL-16]
	_ = x[LESS-17]
	_ = x[LESS_EQUAL-18]
	_ = x[IDENTIFIER-19]
	_ = x[STRING-20]
	_ = x[NUMBER-21]
	_ = x[AND-22]
	_ = x[CLASS-23]
	_ = x[ELSE-24]
	_ = x[FALSE-25]
	_ = x[FUN-26]
	_ = x[FOR-27]
	_ = x[IF-28]
	_ = x[NIL-29]
	_ = x[OR-30]
	_ = x[PRINT-31]
	_ = x[RETURN-32]
	_ = x[SUPER-33]
	_ = x[THIS-34]
	_ = x[TRUE-35]
	_ = x[VAR-36]
	_ = x[WHILE-37]
	_ = x[EOF-38]
}

const _TokenType_name = "LEFT_PARENRIGHT_PARENLEFT_BRACERIGHT_BRACECOMMADOTMINUSPLUSSEMICOLONSLASHSTARBANGBANG_EQUALEQUALEQUAL_EQUALGREATERGREATER_EQUALLESSLESS_EQUALIDENTIFIERSTRINGNUMBERANDCLASSELSEFALSEFUNFORIFNILORPRINTRETURNSUPERTHISTRUEVARWHILEEOF"

var _TokenType_index = [...]uint8{0, 10, 21, 31, 42, 47, 50, 55, 59, 68, 73, 77, 81, 91, 96, 107, 114, 127, 131, 141, 151, 157, 163, 166, 171, 175, 180, 183, 186, 188, 191, 193, 198, 204, 209, 213, 217, 220, 225, 228}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
