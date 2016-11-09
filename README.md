Gem Sessions
============
Gem Sessions is a sessions package for [fasthttp](https://github.com/valyala/fasthttp), it provides cookie and filesystem sessions
and infrastructure for custom session backends.

This project inspired by [gorilla sessions](https://github.com/gorilla/sessions) and [gorilla context](https://github.com/gorilla/context),
their LICENSES can be found in LICENSE file.

## Features

* Simple API: use it as an easy way to set signed (and optionally
  encrypted) cookies.
* Built-in backends to store sessions in cookies or the filesystem.
* Flash messages: session values that last until read.
* Convenient way to switch session persistency (aka "remember me") and set
  other attributes.
* Mechanism to rotate authentication and encryption keys.
* Multiple sessions per request, even using different backends.
* Interfaces and infrastructure for custom session backends: sessions from
  different stores can be retrieved and batch-saved using a common API.

## Usage

### Install
```
go get github.com/go-gem/sessions
```

### Example
```go
package main

import (
	"fmt"
	"time"
	"log"

	"github.com/go-gem/sessions"
	"github.com/valyala/fasthttp"
)

var (
	store sessions.Store
)

func handler(ctx *fasthttp.RequestCtx) {
	// Get session from store.
	session, err := store.Get(ctx, "GOSESSION")
	if err != nil {
        log.Printf("Failed to get session: %s\n", err.Error())
        return
	}

	// Save session.
	defer session.Save(ctx)

	if string(ctx.Path()) == "/set" {
		name := string(ctx.FormValue("name"))
		if len(name) > 0 {
			session.Values["name"] = name
			ctx.SetBodyString(fmt.Sprintf("Name has been set as: %s\n", session.Values["name"]))
		} else {
			ctx.SetBodyString("No name specified.")
		}
		return
	}

	if len(ctx.QueryArgs().Peek("a")) > 0 {
		time.Sleep(time.Second * time.Duration(5))
	}

	if name, ok := session.Values["name"].(string); ok {
		ctx.SetBodyString(fmt.Sprintf("Name: %s\n", name))
		return
	}

	ctx.SetContentType("text/html charset:utf-8")
	ctx.SetBodyString(`
	You should navigate to
	<a href="http://127.0.0.1:8080/set?name=Gem" target="_blank">http://127.0.0.1:8080/set?name=Gem</a>
	to set specified name.
	`)
}

func main() {
	store = sessions.NewCookieStore([]byte("something-very-secret"))
	fasthttp.ListenAndServe(":8080", sessions.ClearHandler(handler))
}
```

First we initialize a session store calling `NewCookieStore()` and passing a
secret key used to authenticate the session. Inside the handler, we call
`store.Get()` to retrieve an existing session or a new one. Then we set some
session values in session.Values, which is a `map[interface{}]interface{}`.
And finally we call `session.Save()` to save the session in the response.

**Important Note**: application **must** to call `sessions.Clear` at the end of a request lifetime.
An easy way to do this is to wrap your handler with `sessions.ClearHandler`.


## Store Implementations

Other implementations of the `sessions.Store` interface:

1. Pending to add.


## License

MIT licensed. See the LICENSE file for details.
