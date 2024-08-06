package main

import (
	"github.com/blocksignalio/core"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func adaptEvents(xs map[string]abi.Event) []Event {
	ys := make([]Event, 0, len(xs))
	for _, e := range xs {
		inputs := make([]Input, len(e.Inputs))
		for i, a := range e.Inputs {
			inputs[i] = Input{
				Name:    a.Name,
				Type:    a.Type.String(),
				Indexed: a.Indexed,
			}
		}
		ys = append(ys, Event{
			Name:      e.RawName,
			Signature: e.Sig,
			ID:        e.ID.Hex(),
			Inputs:    inputs,
		})
	}
	return ys
}

func adaptLogs(xs []core.Log) []Log {
	ys := make([]Log, len(xs))
	for i, x := range xs {
		ys[i] = Log{
			Address:     x.Address,
			Topic0:      x.Topic0,
			Topic1:      x.Topic1,
			Topic2:      x.Topic2,
			Topic3:      x.Topic3,
			Data:        x.Data,
			BlockNumber: x.BlockNumber,
			TxHash:      x.TxHash,
			TxIndex:     x.TxIndex,
			Index:       x.Index,
		}
	}
	return ys
}
