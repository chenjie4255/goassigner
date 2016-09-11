## GO Assigner

As we know, Go is not support assignment between two different structure, but sometimes it's necessary for some aggregation handler or something else, and the only way to actieve this is a hard coding, for example:

```
type Foo struct {
    ID int64
    Name string
    Email string
    ....
} 

type FooX struct {
    ID int64
    Name string
}

func AssignFooToFooX(dest *FooX, src Foo) {
    dest.ID = Foo.ID
    dest.Name = foo.Name
}

``` 

If we have plain of some FooX like struct, We need to write lots of hard code, which is inconvenient for me. Insprited by Go Generate (stringer, [joiner](github.com/bslatkin/joiner)), I wrote [go assigner](github.com/chenjie4255/goassigner) to slove this problem in a different way. now the code look like this:

```
type Foo struct {
    ID int64
    Name string
    Email string
    ....
} 

// @goassigner:Foo
type FooX struct {
    ID int64
    Name string
}

``` 

Just run the ```go generate```, a function ```func (s *FooX)Assign(src Foo)``` is generated automatically.


## Usage

```
go get github.com/chenjie4255/goassigner
go install github.com/chenjie4255/goassigner
 
```

remember set PATH to $GOPATH/bin

More usage see example folder.

thanks.
