package entity

type Header struct {
	Seq    int64 `json:"seq"`
	OpType int   `json:"op_type"`
}
