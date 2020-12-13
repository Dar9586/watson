// Package lexer provides a way to convert a byte sequence into a sequence of Watson's instructions and vice versa.
package lexer

import (
	"fmt"
	"io"

	"github.com/genkami/watson/pkg/vm"
)

// Mode is an important concept that is unique to Watson.
// It determines the correspondence between Vm's instructions and their representation.
type Mode int

const (
	A Mode = iota // A, S are modes of the lexer. See the specification for more details.
	S
)

var opTableA = map[byte]vm.Op{
	char("B"): vm.Inew,
	char("u"): vm.Iinc,
	char("b"): vm.Ishl,
	char("a"): vm.Iadd,
	char("A"): vm.Ineg,
	char("e"): vm.Isht,
	char("i"): vm.Itof,
	char("q"): vm.Finf,
	char("t"): vm.Fnan,
	char("p"): vm.Fneg,
	char("?"): vm.Snew,
	char("!"): vm.Sadd,
	char("~"): vm.Onew,
	char("M"): vm.Oadd,
	char("@"): vm.Anew,
	char("s"): vm.Aadd,
	char("z"): vm.Bnew,
	char("o"): vm.Bneg,
	char("."): vm.Nnew,
	char("*"): vm.Gdup,
	char("#"): vm.Gpop,
	char("%"): vm.Gswp,
}

var reversedTableA map[vm.Op]byte

var opTableS = map[byte]vm.Op{
	char("S"): vm.Inew,
	char("h"): vm.Iinc,
	char("a"): vm.Ishl,
	char("k"): vm.Iadd,
	char("r"): vm.Ineg,
	char("A"): vm.Isht,
	char("z"): vm.Itof,
	char("p"): vm.Finf,
	char("b"): vm.Fnan,
	char("u"): vm.Fneg,
	char("$"): vm.Snew,
	char("-"): vm.Sadd,
	char("+"): vm.Onew,
	char("g"): vm.Oadd,
	char("v"): vm.Anew,
	char("?"): vm.Aadd,
	char("^"): vm.Bnew,
	char("!"): vm.Bneg,
	char("y"): vm.Nnew,
	char("/"): vm.Gdup,
	char("c"): vm.Gpop,
	char(":"): vm.Gswp,
}

var reversedTableS map[vm.Op]byte

func init() {
	reversedTableA = make(map[vm.Op]byte)
	for k, v := range opTableA {
		reversedTableA[v] = k
	}
	reversedTableS = make(map[vm.Op]byte)
	for k, v := range opTableS {
		reversedTableS[v] = k
	}
}

func readOp(m Mode, b byte) (op vm.Op, ok bool) {
	op, ok = opTableA[b]
	return
}

func showOp(m Mode, op vm.Op) byte {
	if b, ok := reversedTableA[op]; ok {
		return b
	}
	panic(fmt.Errorf("unknown Op: %#v\n", op))
}

func char(s string) byte {
	return []byte(s)[0]
}

// Lexer is responsible for converting a Watson Representation into a sequence of vm.Ops.
type Lexer struct {
	r   io.Reader
	buf [1]byte
}

// Creates a new Lexer that reads Watson Representation from r.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{r: r}
}

// Returns the next Op.
// This returns io.EOF if it hits on the end of the input.
func (l *Lexer) Next() (vm.Op, error) {
	for {
		_, err := l.r.Read(l.buf[:])
		if err != nil {
			// Note that it returns io.EOF if the underlying Reader returns io.EOF.
			return 0, err
		}
		if op, ok := readOp(A, l.buf[0]); ok {
			return op, nil
		}
	}
}
