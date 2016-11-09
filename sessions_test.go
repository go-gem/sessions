// Copyright 2016 The Gem Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package sessions

import (
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/valyala/fasthttp"
)

func init() {
	gob.Register(flashMessage{})
}

type flashMessage struct {
	Type    int
	Message string
}

func TestFlashes(t *testing.T) {
	var err error
	var session *Session
	var flashes []interface{}

	store := NewCookieStore([]byte("secret-key"))
	store.Options = &Options{
		Path:     "/user",
		Domain:   "golang.org",
		Secure:   true,
		HttpOnly: true,
	}

	// Round 1
	ctx := &fasthttp.RequestCtx{}
	// Get a session by an empty cookie name.
	if _, err = store.Get(ctx, ""); err.Error() != "sessions: invalid character in cookie name: " {
		t.Fatalf("Expected error due to invalid cookie name")
	}
	// Get a session.
	if session, err = store.Get(ctx, "session-key"); err != nil {
		t.Fatalf("Error getting session: %v", err)
	}
	// Get a flash.
	flashes = session.Flashes()
	if len(flashes) != 0 {
		t.Errorf("Expected empty flashes; Got %v", flashes)
	}
	// Add some flashes.
	session.AddFlash("foo")
	session.AddFlash("bar")
	// Custom key.
	session.AddFlash("baz", "custom_key")
	// Save.
	if err = Save(ctx); err != nil {
		t.Fatalf("Error saving session: %v", err)
	}
	// Get session by an invalid cookie name.
	if _, err = store.Get(ctx, "session:key"); err.Error() != "sessions: invalid character in cookie name: session:key" {
		t.Fatalf("Expected error due to invalid cookie name")
	}

	// Round 2
	// Get a session.
	if session, err = store.Get(ctx, "session-key"); err != nil {
		t.Fatalf("Error getting session: %v", err)
	}
	// Check all saved values.
	flashes = session.Flashes()
	if len(flashes) != 2 {
		t.Fatalf("Expected flashes; Got %v", flashes)
	}
	if flashes[0] != "foo" || flashes[1] != "bar" {
		t.Errorf("Expected foo,bar; Got %v", flashes)
	}
	flashes = session.Flashes()
	if len(flashes) != 0 {
		t.Errorf("Expected dumped flashes; Got %v", flashes)
	}
	// Custom key.
	flashes = session.Flashes("custom_key")
	if len(flashes) != 1 {
		t.Errorf("Expected flashes; Got %v", flashes)
	} else if flashes[0] != "baz" {
		t.Errorf("Expected baz; Got %v", flashes)
	}
	flashes = session.Flashes("custom_key")
	if len(flashes) != 0 {
		t.Errorf("Expected dumped flashes; Got %v", flashes)
	}
}

func TestFlashes2(t *testing.T) {
	var err error
	var session *Session
	var flashes []interface{}

	cookie := &fasthttp.Cookie{}
	cookie.SetKey("session-key")

	store := NewCookieStore([]byte("secret-key"))
	store.Options = &Options{
		Path:     "/user",
		Domain:   "golang.org",
		Secure:   true,
		HttpOnly: true,
	}

	ctx := &fasthttp.RequestCtx{}

	// Get a session.
	if session, err = store.Get(ctx, "session-key"); err != nil {
		t.Fatalf("Error getting session: %v", err)
	}

	// Save.
	if err = Save(ctx); err != nil {
		t.Fatalf("Error saving session: %v", err)
	}
	// Check cookie.
	if !ctx.Response.Header.Cookie(cookie) {
		t.Fatalf("The cookie has not been sent to client.")
		t.FailNow()
	}
	if err = checkCookieOptions(cookie, store.Options); err != nil {
		t.Fatal(err)
	}

	// Round 3
	// Get a session.
	if session, err = store.Get(ctx, "session-key"); err != nil {
		t.Fatalf("Error getting session: %v", err)
	}
	// Get a flash.
	flashes = session.Flashes()
	if len(flashes) != 0 {
		t.Errorf("Expected empty flashes; Got %v", flashes)
	}
	// Add some flashes.
	session.AddFlash(flashMessage{42, "foo"})
	// Save.
	if err = Save(ctx); err != nil {
		t.Fatalf("Error saving session: %v", err)
	}

	// Round 4
	// Get a session.
	if session, err = store.Get(ctx, "session-key"); err != nil {
		t.Fatalf("Error getting session: %v", err)
	}
	// Check all saved values.
	flashes = session.Flashes()
	if len(flashes) != 1 {
		t.Fatalf("Expected flashes; Got %v", flashes)
	}
	custom := flashes[0].(flashMessage)
	if custom.Type != 42 || custom.Message != "foo" {
		t.Errorf("Expected %#v, got %#v", flashMessage{42, "foo"}, custom)
	}
}

func checkCookieOptions(cookie *fasthttp.Cookie, options *Options) error {
	if string(cookie.Domain()) != options.Domain {
		return fmt.Errorf("Expected cookie Domain: %s; Got %s", options.Domain, cookie.Domain())
	}
	if string(cookie.Path()) != options.Path {
		return fmt.Errorf("Expected cookie Path: %s; Got %s", options.Path, cookie.Path())
	}
	if cookie.HTTPOnly() != options.HttpOnly {
		return fmt.Errorf("Expected cookie HTTPOnly: %t; Got %t", options.HttpOnly, cookie.HTTPOnly())
	}
	if cookie.Secure() != options.Secure {
		return fmt.Errorf("Expected cookie Secure: %t; Got %t", options.Secure, cookie.Secure())
	}

	return nil
}
