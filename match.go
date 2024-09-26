package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

const TICK_RATE = 30 // number of ticks the server runs per second

type PMatchLabel struct {
	Open    int `json:"open"`
	Physics int `json:"physics"`
}

type PMatchState struct {
	// match lifecycle related
	playing         bool
	label           *PMatchLabel
	joinsInProgress int

	// user maps
	presences map[string]*runtime.Presence

	game *PGame
}

type PMatch struct{}

func newMatch(
	ctx context.Context,
	logger runtime.Logger,
	db *sql.DB,
	nk runtime.NakamaModule) (m runtime.Match, err error) {
	return &PMatch{}, nil
}

func (m *PMatch) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	state := &PMatchState{
		playing:         false,
		label:           &PMatchLabel{Open: 1, Physics: 1},
		joinsInProgress: 0,

		presences: make(map[string]*runtime.Presence),

		game: newPhysicsGame(),
	}

	label := ""
	labelBytes, err := json.Marshal(state.label)
	if err == nil {
		label = string(labelBytes)
	}

	return state, TICK_RATE, label
}

func (m *PMatch) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	state := state_.(*PMatchState)

	state.joinsInProgress++
	return state, true, ""
}

func (m *PMatch) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, presences []runtime.Presence) interface{} {
	state := state_.(*PMatchState)

	for _, p := range presences {
		state.joinsInProgress--
		id := p.GetUserId()
		state.presences[id] = &p
	}

	state.playing = true

	return state
}

func (m *PMatch) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, presences []runtime.Presence) interface{} {
	state := state_.(*PMatchState)

	for _, p := range presences {
		id := p.GetUserId()
		delete(state.presences, id)
	}

	if len(state.presences) == 0 {
		return nil
	}

	return state
}

func (m *PMatch) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, messages []runtime.MatchData) interface{} {
	state := state_.(*PMatchState)

	for _, msg := range messages {
		sender := msg.GetUserId()
		op := msg.GetOpCode()
		data := msg.GetData()

		logger.Debug("received: %s %d %s", sender, op, data)

		switch op {
		case OpForceBodies:
			// TODO
		default:
			logger.Debug("unsupported opcode: %d", op)
		}
	}

	state.game.update(1.0 / TICK_RATE)
	updateBodies := make(UpdateBodiesBody)
	b := state.game.b1
	p := b.Position()
	updateBodies["b1"] = UpdateBodyValue{
		Position: Pos{p.X, p.Y},
		Angle:    b.Angle(),
	}
	bcUpdate(dispatcher, &updateBodies, nil)

	return state
}

func (m *PMatch) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, data string) (interface{}, string) {
	state := state_.(*PMatchState)

	if data == "kill" {
		return nil, "killing match due to rpc signal"
	}

	return state, ""
}

func (m *PMatch) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state_ interface{}, graceSeconds int) interface{} {
	state := state_.(*PMatchState)

	return state
}
