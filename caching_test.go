package cache

import (
	"errors"
	"testing"
	"time"
)

// fakeCache is an in-memory ICache used to exercise the package-level helpers
// without a real backend.
type fakeCache struct {
	store  map[string]interface{}
	closed bool
	failOn string // key that should force every op to return errBoom
}

var errBoom = errors.New("boom")

func newFake() *fakeCache { return &fakeCache{store: map[string]interface{}{}} }

func (f *fakeCache) Set(key string, value interface{}, _ time.Duration) error {
	if key == f.failOn {
		return errBoom
	}
	f.store[key] = value
	return nil
}

func (f *fakeCache) Get(key string) (interface{}, error) {
	v, ok := f.store[key]
	if !ok {
		return nil, ErrCacheMiss
	}
	return v, nil
}

func (f *fakeCache) Del(key string) error {
	delete(f.store, key)
	return nil
}

func (f *fakeCache) Close() error {
	f.closed = true
	return nil
}

// reset clears the global driver between tests.
func reset() { cache = nil }

func TestHelpersReturnErrNotRegistered(t *testing.T) {
	reset()

	if err := Set("k", "v", 0); !errors.Is(err, ErrNotRegistered) {
		t.Fatalf("Set: want ErrNotRegistered, got %v", err)
	}
	if _, err := Get("k"); !errors.Is(err, ErrNotRegistered) {
		t.Fatalf("Get: want ErrNotRegistered, got %v", err)
	}
	if err := Del("k"); !errors.Is(err, ErrNotRegistered) {
		t.Fatalf("Del: want ErrNotRegistered, got %v", err)
	}
	if err := Close(); !errors.Is(err, ErrNotRegistered) {
		t.Fatalf("Close: want ErrNotRegistered, got %v", err)
	}
}

func TestSetGetDelRoundTrip(t *testing.T) {
	reset()
	f := newFake()
	Register(f)

	if err := Set("greeting", "hello", time.Minute); err != nil {
		t.Fatalf("Set returned error: %v", err)
	}

	got, err := Get("greeting")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got != "hello" {
		t.Fatalf("Get: want %q, got %q", "hello", got)
	}

	if err := Del("greeting"); err != nil {
		t.Fatalf("Del returned error: %v", err)
	}

	if _, err := Get("greeting"); !errors.Is(err, ErrCacheMiss) {
		t.Fatalf("Get after Del: want ErrCacheMiss, got %v", err)
	}
}

func TestGetMissReturnsErrCacheMiss(t *testing.T) {
	reset()
	Register(newFake())

	if _, err := Get("absent"); !errors.Is(err, ErrCacheMiss) {
		t.Fatalf("want ErrCacheMiss, got %v", err)
	}
}

func TestManagerReturnsRegisteredDriver(t *testing.T) {
	reset()
	if Manager() != nil {
		t.Fatalf("Manager should be nil before Register")
	}

	f := newFake()
	Register(f)
	if Manager() != f {
		t.Fatalf("Manager did not return the registered driver")
	}
}

func TestCloseInvokesICloser(t *testing.T) {
	reset()
	f := newFake()
	Register(f)

	if err := Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	if !f.closed {
		t.Fatalf("Close did not call the driver's Close")
	}
}

func TestCloseNoopWithoutICloser(t *testing.T) {
	reset()
	Register(noCloser{})

	if err := Close(); err != nil {
		t.Fatalf("Close should be a no-op, got %v", err)
	}
}

// noCloser implements ICache but not ICloser.
type noCloser struct{}

func (noCloser) Set(string, interface{}, time.Duration) error { return nil }
func (noCloser) Get(string) (interface{}, error)              { return nil, ErrCacheMiss }
func (noCloser) Del(string) error                             { return nil }

func TestKeyNamespacing(t *testing.T) {
	t.Setenv("APP_CODE", "myapp")
	if got, want := Key("session"), "myapp:session"; got != want {
		t.Fatalf("Key: want %q, got %q", want, got)
	}
}

func TestSetPropagatesDriverError(t *testing.T) {
	reset()
	f := newFake()
	f.failOn = "bad"
	Register(f)

	if err := Set("bad", "v", 0); !errors.Is(err, errBoom) {
		t.Fatalf("want errBoom, got %v", err)
	}
}
