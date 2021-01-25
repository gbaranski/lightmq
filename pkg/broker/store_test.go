package broker

import "testing"

func TestClientStore(t *testing.T) {
	cs := NewClientStore()
	c1 := Client{
		ClientID:  "SomeClientID",
		IPAddress: nil,
	}
	c2 := Client{
		ClientID:  "SomeDifferentClientID",
		IPAddress: nil,
	}

	if err := cs.Add(c1); err != nil {
		t.Fatalf("fail add c1 %s", err.Error())
	}
	if err := cs.Add(c2); err != nil {
		t.Fatalf("fail add c2 %s", err.Error())
	}

	c, ok := cs.Get(c1.ClientID)
	if !ok {
		t.Fatalf("fail retreive c1")
	}
	if c.ClientID != c1.ClientID {
		t.Fatalf("invalid c1 clientID")
	}

	c, ok = cs.Get(c2.ClientID)
	if !ok {
		t.Fatalf("fail retreive c2")
	}
	if c.ClientID != c2.ClientID {
		t.Fatalf("invalid c2 clientID")
	}
}

func TestClientStoreNotPresentUser(t *testing.T) {
	cs := NewClientStore()
	c1 := Client{
		ClientID:  "SomeClientID",
		IPAddress: nil,
	}

	c2 := Client{
		ClientID:  "SomeDifferentClientID",
		IPAddress: nil,
	}

	if err := cs.Add(c1); err != nil {
		t.Fatalf("fail add c1 %s", err.Error())
	}

	c, ok := cs.Get(c1.ClientID)
	if !ok {
		t.Fatalf("fail retreive c1")
	}
	if c.ClientID != c1.ClientID {
		t.Fatalf("invalid c1 clientID")
	}

	c, ok = cs.Get(c2.ClientID)
	if ok {
		t.Fatalf("expected to fail on c2")
	}
	if c.ClientID == c2.ClientID {
		t.Fatalf("invalid c2 clientID")
	}
}

func TestClientStoreDuplicate(t *testing.T) {
	cs := NewClientStore()
	c1 := Client{
		ClientID:  "SomeClientID",
		IPAddress: nil,
	}

	if err := cs.Add(c1); err != nil {
		t.Fatalf("fail add c1 %s", err.Error())
	}

	if err := cs.Add(c1); err == nil {
		t.Fatalf("expected error, received nil while adding dup")
	}
}
