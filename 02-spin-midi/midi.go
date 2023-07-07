package main

import (
	"fmt"
	"sync/atomic"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type MidiMgr struct {
	in       drivers.In
	stop     func()
	knob1Val int
	shared   *midiMgrLockState
}

type midiMgrLockState struct {
	k1 int32 // atomic access only!
}

func NewMidiMgr(cfg Config) (*MidiMgr, error) {
	in, err := midi.FindInPort(cfg.MidiDevice)
	if err != nil {
		return nil, fmt.Errorf("FindInPort(%s): %w", cfg.MidiDevice, err)
	}

	shared := &midiMgrLockState{}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		var ch, controller, value uint8
		switch {
		case msg.GetControlChange(&ch, &controller, &value):
			if ch == uint8(cfg.Knob1Chan) && controller == uint8(cfg.Knob1Controller) {
				atomic.StoreInt32(&shared.k1, int32(value))
			}
		default:
			// ignore
		}
	})
	if err != nil {
		return nil, fmt.Errorf("ListenTo(%s): %w", cfg.MidiDevice, err)
	}

	return &MidiMgr{
		in:       in,
		stop:     stop,
		knob1Val: 0,
		shared:   shared,
	}, nil
}

func (m *MidiMgr) Close() {
	m.stop()
}

func (m *MidiMgr) Update() {
	m.knob1Val = int(atomic.LoadInt32(&m.shared.k1))
}

func (m *MidiMgr) Knob1() int {
	return m.knob1Val
}
