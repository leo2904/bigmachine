// Copyright 2018 GRAIL, Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package bigmachine

import (
	"context"
	"io"
	"net/http"
	"time"
)

// A System implements a set of methods to set up a bigmachine and
// start new machines. Systems are also responsible for providing an
// HTTP client that can be used to communicate between machines
// and drivers.
type System interface {
	// Name is the name of this system. It is used to multiplex multiple
	// system implementations, and thus should be unique among
	// systems.
	Name() string
	// Init is called when the bigmachine starts up in order to
	// initialize the system implementation. If an error is returned,
	// the Bigmachine fails.
	Init(*B) error
	// Main is called to start  a machine. The system is expected to
	// take over the process; the bigmachine fails if main returns (and
	// if it does, it should always return with an error).
	Main() error
	// HTTPClient returns an HTTP client that can be used to communicate
	// from drivers to machines as well as between machines.
	HTTPClient() *http.Client
	// ListenAndServe serves the provided handler on an HTTP server that
	// is reachable from other instances in the bigmachine cluster. If addr
	// is the empty string, the default cluster address is used.
	ListenAndServe(addr string, handle http.Handler) error
	// Start launches up to n new machines.  The returned machines can
	// be in Unstarted state, but should eventually become available.
	Start(ctx context.Context, n int) ([]*Machine, error)
	// Exit is called to terminate a machine with the provided exit code.
	Exit(int)
	// Shutdown is called on graceful driver exit. It's should be used to
	// perform system tear down. It is not guaranteed to be called.
	Shutdown()
	// Maxprocs returns the maximum number of processors per machine,
	// as configured. Returns 0 if is a dynamic value.
	Maxprocs() int
	// KeepaliveConfig returns the various keepalive timeouts that should
	// be used to maintain keepalives for machines started by this system.
	KeepaliveConfig() (period, timeout, rpcTimeout time.Duration)
	// Tail should tail (follow) machine m's log output to the provided writer.
	Tail(ctx context.Context, w io.Writer, m *Machine) error
}
