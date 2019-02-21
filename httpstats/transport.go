package httpstats

import (
	"net/http"
	"time"

	"context"
	"github.com/segmentio/stats"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
// Also, I completely stole this.
type contextKey struct {
	name string
}

// String is Stringer implementation
func (k *contextKey) String() string {
	return "stats context value " + k.name
}

// contextKeyReqTags is contextKey for tags
var (
	contextKeyReqTags = &contextKey{
		name: "segmentio_httpstats_req_tags",
	}
)

// NewTransport wraps t to produce metrics on the default engine for every request
// sent and every response received.
func NewTransport(t http.RoundTripper) http.RoundTripper {
	return NewTransportWith(stats.DefaultEngine, t)
}

// NewTransportWith wraps t to produce metrics on eng for every request sent and
// every response received.
func NewTransportWith(eng *stats.Engine, t http.RoundTripper) http.RoundTripper {
	return &transport{
		transport: t,
		eng:       eng,
	}
}

type transport struct {
	transport http.RoundTripper
	eng       *stats.Engine
}

// RequestWithContext returns a shallow copy of req with its context changed with this provided tags
// so the they can be used later during the RoundTrip in the metrics recording.
// The provided ctx must be non-nil.
func RequestWithContext(req *http.Request, tags ...stats.Tag) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, contextKeyReqTags, tags)
	return req.WithContext(ctx)
}

// RoundTrip implements http.RoundTripper
func (t *transport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	start := time.Now()
	rtrip := t.transport
	eng := t.eng

	if rtrip == nil {
		rtrip = http.DefaultTransport
	}

	if tags, ok := req.Context().Value(contextKeyReqTags).([]stats.Tag); ok {
		eng = eng.WithTags(tags...)
	}

	if req.Body == nil {
		req.Body = &nullBody{}
	}

	m := &metrics{}

	req.Body = &requestBody{
		eng:     eng,
		req:     req,
		metrics: m,
		body:    req.Body,
		op:      "write",
	}

	res, err = rtrip.RoundTrip(req)
	// safe guard, the transport should have done it already
	req.Body.Close() // nolint

	if err != nil {
		m.observeError(time.Now().Sub(start))
		eng.ReportAt(start, m)
	} else {
		res.Body = &responseBody{
			eng:     eng,
			res:     res,
			metrics: m,
			body:    res.Body,
			op:      "read",
			start:   start,
		}
	}

	return
}
