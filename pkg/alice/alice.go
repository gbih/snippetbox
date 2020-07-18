package alice

import (
	"fmt"
	"net/http"
)

// Constraining the elements to have this Middleware-pattern function signature
type Constructor func(http.Handler) http.Handler

// Enable slice type of the Middleware functions
type Chain struct {
	constructors []Constructor
}

// Create object of type Chain slice, to hold our functions.
// We create this by creating a Constructor slice, then appending our functions to it.
// To enable writing this, this needs to be a variadic function
// standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
// otherwise, we can only pass one argument
func New(constructors ...Constructor) Chain {
	tmp := Chain{append(([]Constructor(nil)), constructors...)}
	fmt.Printf("ALICE/New: %#v\n", tmp)
	return tmp
}

func (c Chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := range c.constructors {
		h = c.constructors[len(c.constructors)-1-i](h) // var h http.Handler
		fmt.Print(len(c.constructors) - 1 - i)
		fmt.Printf("ALICE/Then: %#v\n", h)
	}
	return h
}

func (c Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return c.Then(nil)
	}
	fmt.Printf("ALICE/Then: %#v\n", c.Then(fn))
	return c.Then(fn)
}

// Variadic function
func (c Chain) Append(constructors ...Constructor) Chain {
	newCons := make([]Constructor, 0, len(c.constructors)+len(constructors))
	fmt.Println("newCons 1", newCons)
	newCons = append(newCons, c.constructors...)
	fmt.Println("newCons 2", newCons)
	newCons = append(newCons, constructors...)
	fmt.Println("newCons 3", newCons)

	fmt.Println("Chain{newCons}", Chain{newCons})

	return Chain{newCons}
}

func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.constructors...)
}
