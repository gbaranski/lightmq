package lightmq

func TestClientStore() {
	cs := NewClientStore()
	cs.Add(Client{
		ClientID:  "someClientID",
		IPAddress: nil,
	})

}
