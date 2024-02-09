package client

import (
	"bytes"
	"testing"
	"time"

	"github.com/joaovictorsl/dcache"
	"github.com/joaovictorsl/dcache/core/cache"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const (
	s1Addr = "127.0.0.1:3000"
	s2Addr = "127.0.0.1:3001"
	s3Addr = "127.0.0.1:3002"
	s4Addr = "127.0.0.1:3003"
	s5Addr = "127.0.0.1:3004"
)

var client *DCacheClient

func TestMain(m *testing.M) {
	addresses := []string{s1Addr, s2Addr, s3Addr, s4Addr, s5Addr}
	for _, addr := range addresses {
		s := dcache.NewServer(addr, cache.NewSimple())
		go s.Start()
	}

	// Give some time so servers are up
	time.Sleep(200 * time.Millisecond)
	client = New(addresses...)

	m.Run()
}

func TestNew(t *testing.T) {
	addrList := []string{s1Addr, s2Addr}
	c := New(addrList...)

	// Connection keys can be returned in any given order,
	// slices.Equal takes into consideration the order to
	// determine wether slices are equal or nor. In our case
	// the order doesn't matter. Therefore, we sort it so the test won't fail
	connsKeys := maps.Keys(c.conns)
	slices.Sort(connsKeys)

	if !slices.Equal(connsKeys, addrList) {
		t.Errorf("expected (%v), got (%s)", addrList, connsKeys)
	} else if c.mu == nil {
		t.Errorf("expected client mutex to be non-nil, it was nil")
	}

	for _, dconn := range c.conns {
		if dconn.active {
			t.Errorf("expected dconn to not be active, it was active")
		} else if dconn.conn != nil {
			t.Errorf("expected dconn conn to be nil, it was not")
		} else if dconn.mu == nil {
			t.Errorf("expected dconn mutex to be non-nil, it was nil")
		}
	}
}

func TestConnect(t *testing.T) {
	client.Connect(2, 2*time.Second)

	for _, dconn := range client.conns {
		if !dconn.active {
			t.Errorf("expected dconn to be active, it was not active")
		} else if dconn.conn == nil {
			t.Errorf("expected dconn conn to be non-nil, it was nil")
		}
	}
}

func TestCommands(t *testing.T) {
	client.Connect(2, 2*time.Second)

	cases := [][]byte{
		[]byte("Foo"), []byte("Bar"),
		[]byte("Some other key"), []byte("B"),
		[]byte("Tamtamtam"), []byte("D"),
		[]byte("Plant"), []byte("F"),
		[]byte("grasping"), []byte("H"),
		[]byte("example"), []byte("J"),
		[]byte("Kikikiki funny"), []byte("L"),
		[]byte("My car is yellow"), []byte("N"),
		[]byte("oh no, our table is broken"), []byte("P"),
		[]byte("something else"), []byte("S"),
	}

	for i := 0; i < len(cases); i += 2 {
		k := cases[i]
		v := cases[i+1]

		err := client.Set(string(k), v, 10000)
		if err != nil {
			t.Errorf("no error was expected on SET operation, but got: %s", err)
		}

		res, ok, err := client.Get(string(k))
		if err != nil {
			t.Errorf("no error was expected on GET operation, but got: %s", err)
		} else if !ok {
			t.Errorf("expected key %s to be found, but wasn't", string(k))
		} else if !bytes.Equal(res, v) {
			t.Errorf("expected GET command on %s key to return %s but got: %v | %s", string(k), string(v), res, string(res))
		}
	}
}

func TestClose(t *testing.T) {
	client.Connect(2, 2*time.Second)
	client.End()

	// Connection keys can be returned in any given order,
	// slices.Equal takes into consideration the order to
	// determine wether slices are equal or nor. In our case
	// the order doesn't matter. Therefore, we sort it so the test won't fail
	connsKeys := maps.Keys(client.conns)

	if len(connsKeys) != 0 {
		t.Errorf("expected connection keys to be empty, got (%s)", connsKeys)
	} else if !client.done {
		t.Error("expected client to be done, but it wasn't")
	}
}
