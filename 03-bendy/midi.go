package main

import (
	"fmt"
	"sync/atomic"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

const KNOB_COUNT = 4

type MidiMgr struct {
	in      drivers.In
	stop    func()
	knobVal [KNOB_COUNT]int
	shared  *midiMgrLockState
}

type midiMgrLockState struct {
	knob [KNOB_COUNT]atomic.Int32
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
			// TODO: once Go 1.21 comes out: `max := min(len(cfg.Knobs), KNOB_COUNT)`
			max := len(cfg.Knobs)
			if max > KNOB_COUNT {
				max = KNOB_COUNT
			}
			for i := 0; i < max; i++ {
				knob := cfg.Knobs[i]
				if ch == uint8(knob.Channel) && controller == uint8(knob.Controller) {
					shared.knob[i].Store(int32(value))
				}
			}
		default:
			// ignore
		}
	})
	if err != nil {
		return nil, fmt.Errorf("ListenTo(%s): %w", cfg.MidiDevice, err)
	}

	return &MidiMgr{
		in:     in,
		stop:   stop,
		shared: shared,
	}, nil
}

func (m *MidiMgr) Close() {
	m.stop()
}

func (m *MidiMgr) Update() {
	for i := range m.knobVal {
		m.knobVal[i] = int(m.shared.knob[i].Load())
	}
}

func (m *MidiMgr) Knob(n int) int {
	return m.knobVal[n]
}
