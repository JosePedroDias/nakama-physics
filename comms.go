package main

import (
	"encoding/json"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
)

const LOG_OUTGOING = false

type OpCode int

// game opcodes must be POSITIVE integers
const (
	// outgoing (starting in 100)
	OpUpdateBodies OpCode = iota + 100

	// incoming (starting in 200)
	OpForceBodies = iota + 199
)

type Pos = [2]float64
type UpdateBodyValue = struct {
	Position Pos     `json:"p"`
	Angle    float64 `json:"a"`
}

// outgoing
type UpdateBodiesBody = map[string]UpdateBodyValue

// incoming
type ForceBodiesBody = map[string]Pos

////

func getJustSender(state *PMatchState, userId string) []runtime.Presence {
	destinations := make([]runtime.Presence, 0)
	destinations = append(destinations, *state.presences[userId])
	return destinations
}

//lint:ignore U1000 optional method
func getAllButSender(state *PMatchState, userId string) []runtime.Presence {
	destinations := make([]runtime.Presence, 0)

	for k, v := range state.presences {
		if k != userId {
			destinations = append(destinations, *v)
		}
	}

	return destinations
}

////

func bcUpdate(dispatcher runtime.MatchDispatcher, updateBodies *UpdateBodiesBody, destinations []runtime.Presence) {
	data, err := json.Marshal(*updateBodies)
	if LOG_OUTGOING {
		fmt.Printf("%d %s\n", OpUpdateBodies, data)
	}
	if err == nil {
		dispatcher.BroadcastMessage(int64(OpUpdateBodies), data, destinations, nil, true)
	}
}
