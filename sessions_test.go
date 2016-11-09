// Copyright 2016 The Gem Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package sessions

import (
	"encoding/gob"
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

	cookie := &fasthttp.Cookie{}
	cookie.SetKey("session-key")

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
	// Check cookie.
	if !ctx.Response.Header.Cookie(cookie) {
		t.Fatalf("The cookie has not been sent to client.")
		t.FailNow()
	}
	if string(cookie.Domain()) != session.Options.Domain {
		t.Errorf("Expected cookie Domain: %s; Got %s", session.Options.Domain, cookie.Domain())
	}
	if string(cookie.Path()) != session.Options.Path {
		t.Errorf("Expected cookie Path: %s; Got %s", session.Options.Path, cookie.Path())
	}
	if cookie.HTTPOnly() != session.Options.HttpOnly {
		t.Errorf("Expected cookie HTTPOnly: %t; Got %t", session.Options.HttpOnly, cookie.HTTPOnly())
	}
	if cookie.Secure() != session.Options.Secure {
		t.Errorf("Expected cookie Secure: %t; Got %t", session.Options.Secure, cookie.Secure())
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
