package midioutput

import (
	"fmt"
	"log"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/gm"
	"gitlab.com/gomidi/midi/v2/smf"

	"github.com/proatscratch/hackathon/evomusic/noteevent"
)

const Rest int8 = 127 // sentinel value represents a rest in the semitone notation

func midiNotes(rootNote int8, semitones []int8) []noteevent.NoteEvent {
	currentNote := int(rootNote)
	var events []noteevent.NoteEvent

	for _, step := range semitones {
		if step == Rest {
			events = append(events, noteevent.NoteEvent{
				IsRest: true,
			})
		} else {
			currentNote += int(step)
			events = append(events, noteevent.NoteEvent{
				Note:   int8(currentNote),
				IsRest: false,
			})
		}
	}

	return events
}

func MakeMidiFile(rootNote uint8, semitones []int8, bpm uint16) {
	// 1. Setup the MIDI file template (960 ticks per quarter note)
	s := smf.New()
	s.TimeFormat = smf.MetricTicks(960)

	var tr smf.Track
	tr.Add(0, smf.MetaTempo(float64(bpm)))
	tr.Add(0, smf.MetaTrackSequenceName("Interval Melody"))
	tr.Add(0, midi.ProgramChange(0, uint8(gm.Instr_Violin))) // Acoustic Grand Piano

	// 2. Generate the events using your helper function
	events := midiNotes(int8(rootNote), semitones)
	noteDuration := uint32(960) // 1 quarter note length
	pendingDelta := uint32(0)   // Tracks time accumulation for rests

	// 3. Loop through the generated events and add them to the track
	for _, event := range events {
		if event.IsRest {
			// Accumulate rest duration instead of playing a note
			pendingDelta += noteDuration
			continue
		}

		// Note On combines the initial delta (0 or accumulated rest time)
		tr.Add(pendingDelta, midi.NoteOn(0, uint8(event.Note), 90))

		// Note Off after duration ticks
		tr.Add(noteDuration, midi.NoteOff(0, uint8(event.Note)))

		// Reset the pending delta after playing a valid note
		pendingDelta = 0
	}

	// 4. Close and write the file
	tr.Close(0)
	s.Add(tr)

	err := s.WriteFile("melody.mid")
	if err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	fmt.Println("Successfully generated melody.mid!")
}
