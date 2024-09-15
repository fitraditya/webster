package node

type Delegate struct {
	messages chan []byte
}

func NewDelegate() *Delegate {
	return &Delegate{
		messages: make(chan []byte),
	}
}

func (d *Delegate) NotifyMsg(msg []byte) {
	d.messages <- msg
}

func (d *Delegate) NodeMeta(limit int) []byte {
	return []byte("")
}

func (d *Delegate) LocalState(join bool) []byte {
	return []byte("")
}

func (d *Delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (d *Delegate) MergeRemoteState(buf []byte, join bool) {
	//
}
