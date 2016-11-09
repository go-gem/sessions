// Copyright 2016 The Gem Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package sessions

import (
	"encoding/base64"
	"testing"

	"github.com/valyala/fasthttp"
)

// Test for GH-8 for CookieStore
func TestGH8CookieStore(t *testing.T) {
	originalPath := "/"
	store := NewCookieStore()
	store.Options.Path = originalPath
	ctx := &fasthttp.RequestCtx{}

	session, err := store.New(ctx, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	store.Options.Path = "/foo"
	if session.Options.Path != originalPath {
		t.Fatalf("bad session path: got %q, want %q", session.Options.Path, originalPath)
	}
}

// Test for GH-8 for FilesystemStore
func TestGH8FilesystemStore(t *testing.T) {
	originalPath := "/"
	store := NewFilesystemStore("")
	store.Options.Path = originalPath

	ctx := &fasthttp.RequestCtx{}

	session, err := store.New(ctx, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	store.Options.Path = "/foo"
	if session.Options.Path != originalPath {
		t.Fatalf("bad session path: got %q, want %q", session.Options.Path, originalPath)
	}
}

// Test for GH-2.
func TestGH2MaxLength(t *testing.T) {
	store := NewFilesystemStore("", []byte("some key"))
	ctx := &fasthttp.RequestCtx{}

	session, err := store.New(ctx, "my session")
	session.Values["big"] = make([]byte, base64.StdEncoding.DecodedLen(4096*2))
	err = session.Save(ctx)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	store.MaxLength(4096 * 3) // A bit more than the value size to account for encoding overhead.
	err = session.Save(ctx)
	if err != nil {
		t.Fatal("failed to Save:", err)
	}
}

// Test delete filesystem store with max-age: -1
func TestGH8FilesystemStoreDelete(t *testing.T) {
	store := NewFilesystemStore("", []byte("some key"))
	ctx := &fasthttp.RequestCtx{}

	session, err := store.New(ctx, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	err = session.Save(ctx)
	if err != nil {
		t.Fatal("failed to save session", err)
	}

	session.Options.MaxAge = -1
	err = session.Save(ctx)
	if err != nil {
		t.Fatal("failed to delete session", err)
	}
}

// Test delete filesystem store with max-age: 0
func TestGH8FilesystemStoreDelete2(t *testing.T) {
	store := NewFilesystemStore("", []byte("some key"))
	ctx := &fasthttp.RequestCtx{}

	session, err := store.New(ctx, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	err = session.Save(ctx)
	if err != nil {
		t.Fatal("failed to save session", err)
	}

	session.Options.MaxAge = 0
	err = session.Save(ctx)
	if err != nil {
		t.Fatal("failed to delete session", err)
	}
}
