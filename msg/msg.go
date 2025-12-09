package msg

type Putmsg struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}
type Getmsg struct {
	Key string `json:"Key"`
}
