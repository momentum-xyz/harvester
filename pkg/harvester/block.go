package harvester

type Event struct {
	Method string
	Data   string
}

type Extrinsic struct {
	Signer string
	Method string
	Args   map[string]interface{}
	Events []Event
}

type Block struct {
	Number          uint32      `json:"number"`
	AuthorID        string      `json:"authorId"`
	Finalized       bool        `json:"finalized"`
	ExtrinsicsCount int         `json:"extrinsicsCount"`
	Extrinsics      []Extrinsic `json:"extrinsics"`
	Chain           string      `json:"chain"`
}
