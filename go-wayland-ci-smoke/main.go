//go:build linux
// +build linux

// Copyright 2026 Neurlang project

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// Command go-wayland-ci-smoke is a minimal headless Wayland client used by CI
// (see .github/workflows/ci.yml) to prove two things against a real, running
// Wayland compositor (a headless weston in CI):
//
//  1. This Go Wayland client library can actually connect to a compositor,
//     bind the required globals, create a toplevel window and talk the real
//     wire protocol end to end (not just compile).
//  2. Window.SetAppID sends a real xdg_toplevel.set_app_id request. Run the
//     compositor with WAYLAND_DEBUG=1 and grep its log for "set_app_id" to
//     confirm the request was actually received.
//
// It intentionally avoids any rendering (no SHM/EGL buffers, no cairo, no
// xkbcommon keymap loading) to stay light enough for a CI runner: creating
// the toplevel is enough to get an *xdg.Toplevel to call SetAppID on.
package main

import (
	"fmt"
	"os"

	"github.com/neurlang/wayland/wlclient"
	"github.com/neurlang/wayland/window"
)

func main() {
	appID := "org.unxed.wayland.ci-smoke"
	if len(os.Args) > 1 && os.Args[1] != "" {
		appID = os.Args[1]
	}

	d, err := window.DisplayCreate(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "go-wayland-ci-smoke: DisplayCreate failed:", err)
		os.Exit(1)
	}

	w := window.Create(d)
	if w == nil {
		fmt.Fprintln(os.Stderr, "go-wayland-ci-smoke: window.Create failed")
		os.Exit(1)
	}

	w.SetAppID(appID)

	// Force the request above onto the wire and wait for the server to
	// process it before we exit and close the socket.
	if err := wlclient.DisplayRoundtrip(d.Display); err != nil {
		fmt.Fprintln(os.Stderr, "go-wayland-ci-smoke: roundtrip failed:", err)
		os.Exit(1)
	}

	fmt.Println("go-wayland-ci-smoke: connected to compositor and sent SetAppID(" + appID + ")")
}
