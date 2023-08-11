package sdsnowflake

import (
	"sync"
	"time"

	"github.com/gaorx/stardust5/sderr"
	"github.com/gaorx/stardust5/sdlocal"
)

const (
	nodeBits        = 10
	stepBits        = 12
	nodeMax         = -1 ^ (-1 << nodeBits)
	stepMask  int64 = -1 ^ (-1 << stepBits)
	timeShift uint8 = nodeBits + stepBits
	nodeShift uint8 = stepBits
)

var Epoch int64 = 1288834974657

type Node struct {
	mux  sync.Mutex
	time int64
	node int64
	step int64
}

func New(node int64) (*Node, error) {
	if node < 0 || node > nodeMax {
		return nil, sderr.New("node number must be between 0 and 1023")
	}
	return &Node{
		time: 0,
		node: node,
		step: 0,
	}, nil
}

func NewFromIP() (*Node, error) {
	// IP后10位作为node
	ip, err := sdlocal.IP(sdlocal.Is4(), sdlocal.IsPrivate())
	if err != nil {
		return nil, sderr.Wrap(err, "sdsnowflake get local ip error")
	}
	var node int64 = 0
	if ip != nil {
		h := int64([]byte(ip)[2]) & int64(0x03) // 0b00000011
		l := int64([]byte(ip)[3])
		node = (h << 1) | l
	}
	return New(node)
}

func (n *Node) Generate() int64 {
	n.mux.Lock()
	defer n.mux.Unlock()

	now := time.Now().UnixNano() / 1000000
	if n.time == now {
		n.step = (n.step + 1) & stepMask
		if n.step == 0 {
			for now <= n.time {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		n.step = 0
	}
	n.time = now
	return (now-Epoch)<<timeShift | (n.node << nodeShift) | (n.step)
}
