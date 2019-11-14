package readline

import (
	"bufio"
	"fmt"
	"io"
)

type ContextHelper interface {
	// Readline will pass the whole line and current offset to it
	// Completer need to pass all the candidates, and how long they shared the same characters in line
	// Example:
	//   [go, git, git-shell, grep]
	//   Do("g", 1) => ["o", "it", "it-shell", "rep"], 1
	//   Do("gi", 2) => ["t", "t-shell"], 2
	//   Do("git", 3) => ["", "-shell"], 3
	Do(line []rune, pos int) (help []rune, lines int)
}

type Helper struct{}

func (t *Helper) Do([]rune, int) ([]rune, int) {
	return nil, 0
}

type opHelper struct {
	w     io.Writer
	op    *Operation
	width int
}

func newOpHelper(w io.Writer, op *Operation, width int) *opHelper {
	return &opHelper{
		w:     w,
		op:    op,
		width: width,
	}
}

func (o *opHelper) OnHelp() bool {
	if o.width == 0 {
		return false
	}

	lineCnt := o.op.buf.CursorLineCount()

	help, lines := o.op.cfg.ContextHelp.Do(o.op.buf.Runes(), o.op.buf.idx)
	if help == nil || lines == 0 {
		return true
	}

	buf := bufio.NewWriter(o.w)
	buf.WriteString(string(help))
	fmt.Fprintf(buf, "\033[%dA\r", lineCnt-1+lines)
	fmt.Fprintf(buf, "\033[%dC", o.op.buf.idx+o.op.buf.PromptLen())
	buf.Flush()

	return true
}

func (o *opHelper) OnWidthChange(newWidth int) {
	o.width = newWidth
}
