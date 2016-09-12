package example

//go:generate goassigner -f=$GOFILE

var Foo int64

// @goassigner:User:github.com/chenjie4255/goassigner/example/model
type Child struct {
	ID   int64  `gas:"-"`
	Name string `gas:"-"`
	Age  int64  `gas:"-"`
}

// StructA is a text structure
// @goassigner:User:github.com/chenjie4255/goassigner/example/model
type UserBrief struct {
	ID   int64  `gas:"-"`
	Name string `gas:"-"`

	// @goassigner:VIPInfo:github.com/chenjie4255/goassigner/example/model
	VIP struct {
		VIPLevel int64 `gas:"-"`

		// @goassigner:VIPExInfo:github.com/chenjie4255/goassigner/example/model
		ExInfo struct {
			SumCost int64 `gas:"-"`
		}
	}
	Children []Child
}
