package example

import (
	"github.com/chenjie4255/goassigner/example/st"
)

//go:generate goassigner -f=$GOFILE

var Foo int64

// @goassigner:StructB:github.com/chenjie4255/goassigner/example/st
type ArrayFoo struct {
	Foo1 string `gas:"-"`
}

// StructA is a text structure
// @goassigner:StructB:github.com/chenjie4255/goassigner/example/st
type StructA struct {
	Foo1 string `gas:"-"`
	Foo2 string `gas:"-"`

	// @goassigner:StructB:github.com/chenjie4255/goassigner/example/st
	Fxx struct {
		Foo1 string `gas:"-"`
		Foo2 string `gas:"-"`

		// @goassigner:StructB:github.com/chenjie4255/goassigner/example/st
		Fxxx struct {
			Foo1 string `gas:"-"`
		}
	}
	FA []ArrayFoo

	Txx            st.CustomT `gas:"-"`
	st.EmbedStruct `gas:"-"`  // not support for now...
}
