package midioutput

import (
	"fmt"
	"math/rand/v2"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

// MakeMidiFile selects a random Bach chorale, transposes its notes, sets a custom BPM, and saves it.
func MakeMidiFile(title string, bpm uint16, transpose int8) error {
	// 1. Pick a random Bach chorale
	randNum := rand.IntN(372)
	choraleFile := fmt.Sprintf("../motif-find/jsb_chorales_midi/bach_chorale_%d.mid", randNum)

	// Read the selected chorale
	choraleSMF, err := smf.ReadFile(choraleFile)
	if err != nil {
		return fmt.Errorf("failed to read target midi file %s: %v", choraleFile, err)
	}

	// 2. Setup the new SMF, inheriting the time format
	newSMF := smf.New()
	newSMF.TimeFormat = choraleSMF.TimeFormat

	// 3. Iterate through all tracks and events to shift pitch and set tempo
	for _, track := range choraleSMF.Tracks {
		var newTrack smf.Track

		for _, ev := range track {
			delta := ev.Delta
			msg := ev.Message

			var ch, key, vel uint8
			var trackBPM float64

			switch {
			// Intercept NoteOn messages to transpose pitch
			case msg.GetNoteOn(&ch, &key, &vel):
				newKey := int(key) + int(transpose)
				if newKey < 0 {
					newKey = 0
				}
				if newKey > 127 {
					newKey = 127
				}

				// Cast midi.Message to smf.Message
				msg = smf.Message(midi.NoteOn(ch, uint8(newKey), vel))

			// Intercept NoteOff messages to transpose pitch
			case msg.GetNoteOff(&ch, &key, &vel): // NoteOff in gomidi v2 gives ch, key, vel
				newKey := int(key) + int(transpose)
				if newKey < 0 {
					newKey = 0
				}
				if newKey > 127 {
					newKey = 127
				}

				// Cast midi.Message to smf.Message
				msg = smf.Message(midi.NoteOff(ch, uint8(newKey)))

			// Intercept Tempo changes to enforce our requested BPM
			case msg.GetMetaTempo(&trackBPM):
				if bpm > 0 {
					// smf.MetaTempo already returns an smf.Message, so no cast needed here
					msg = smf.MetaTempo(float64(bpm))
				}
			}

			// Add the modified (or original) event to our new track
			newTrack.Add(delta, msg)
		}

		// Add the processed track to the new SMF
		newSMF.Add(newTrack)
	}

	// 4. Write the modified tracks to the new output file
	outputPath := fmt.Sprintf("%s.mid", title)
	err = newSMF.WriteFile(outputPath)
	if err != nil {
		return fmt.Errorf("failed to write midi file: %v", err)
	}

	fmt.Printf("Successfully transposed and saved %s -> %s\n", choraleFile, outputPath)
	return nil
}
