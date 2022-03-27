package features_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/hedhyw/bdd-resizer-example/internal/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runTestServer(tb testing.TB) (addr string) {
	tb.Helper()

	s := server.New()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(tb, err)

	tb.Cleanup(func() {
		cerr := ln.Close()
		assert.NoError(tb, cerr)
	})

	go func() {
		serr := s.Serve(ln)
		if serr != nil && !errors.Is(serr, http.ErrServerClosed) {
			assert.NoError(tb, err)
		}
	}()

	return ln.Addr().String()
}

type testHelper struct {
	addr string
}

func newTestHelper(tb testing.TB) *testHelper {
	addr := runTestServer(tb)

	_, port, err := net.SplitHostPort(addr)
	require.NoError(tb, err, addr)
	addr = "http://127.0.0.1:" + port

	return &testHelper{
		addr: addr,
	}
}

func (th testHelper) CallAPIGet(tb testing.TB, path string, query url.Values) (resp *http.Response) {
	tb.Helper()

	u, err := url.Parse(th.addr)
	require.NoError(tb, err, th.addr)

	u.Path = path
	u.RawQuery = query.Encode()

	tb.Logf("calling GET %s", u)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	require.NoError(tb, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(tb, err)

	tb.Cleanup(func() {
		err := resp.Body.Close()
		assert.NoError(tb, err)
	})

	return resp
}
