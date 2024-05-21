package astcontext

import (
	"errors"
	"sort"
)

type Decls []*Decl

func (f Decls) NextFuncShift(offset, shift int) (*Decl, error) {
	return f.nextFuncShift(offset, shift)
}

func (f Decls) PrevFuncShift(offset, shift int) (*Decl, error) {
	return f.prevFuncShift(offset, shift)
}

func (f Decls) nextFuncShift(offset, shift int) (*Decl, error) {
	if shift < 0 {
		return nil, errors.New("shift can't be negative")
	}

	// find nearest next function
	nextIndex := sort.Search(len(f), func(i int) bool {
		return f[i].DeclPos.Offset > offset
	})

	if nextIndex >= len(f) {
		return nil, errors.New("no functions found")
	}

	fn := f[nextIndex]

	// if our position is inside the doc, increase the shift by one to pick up
	// the next function. This assumes that people editing a doc of a func want
	// to pick up the next function instead of the current function.
	if fn.Doc != nil && fn.Doc.IsValid() {
		if fn.Doc.Offset <= offset && offset < fn.DeclPos.Offset {
			shift++
		}
	}

	if nextIndex+shift >= len(f) {
		return nil, errors.New("no functions found")
	}

	return f[nextIndex+shift], nil
}

func (f Decls) prevFuncShift(offset, shift int) (*Decl, error) {
	if shift < 0 {
		return nil, errors.New("shift can't be negative")
	}

	// start from the reverse to get the prev function
	f.Reserve()

	prevIndex := sort.Search(len(f), func(i int) bool {
		return f[i].DeclPos.Offset < offset
	})

	if prevIndex >= len(f) {
		return nil, errors.New("no functions found")
	}

	fn := f[prevIndex]

	if fn.NamePos != nil && fn.NamePos.Offset >= offset {
		shift++
	}

	if prevIndex+shift >= len(f) {
		return nil, errors.New("no functions found")
	}

	return f[prevIndex+shift], nil
}

func (f Decls) Len() int      { return len(f) }
func (f Decls) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f Decls) Less(i, j int) bool {
	return f[i].DeclPos.Offset < f[j].DeclPos.Offset
}

// Reserve reserves the Function data
func (f Decls) Reserve() {
	for start, end := 0, f.Len()-1; start < end; {
		f.Swap(start, end)
		start++
		end--
	}
}
