package main

type Input struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Indexed bool   `json:"indexed"`
}

type Event struct {
	Name      string  `json:"name"`
	Signature string  `json:"signature"`
	ID        string  `json:"id"`
	Inputs    []Input `json:"inputs"`
}

type Log struct {
	Address     string `json:"address"`
	Topic0      string `json:"topic0"`
	Topic1      string `json:"topic1"`
	Topic2      string `json:"topic2"`
	Topic3      string `json:"topic3"`
	Data        string `json:"data"`
	BlockNumber uint64 `json:"blockNumber"`
	TxHash      string `json:"txHash"`
	TxIndex     uint   `json:"txIndex"`
	Index       uint   `json:"index"`
}
