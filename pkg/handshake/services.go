package handshake

type ServiceFlag = uint64

const (
	FullNodeFlag ServiceFlag = 1 << iota
)
