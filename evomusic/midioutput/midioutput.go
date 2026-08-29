package midioutput

import (
	"fmt"
	"math/rand/v2"
	"time" // Make sure time is imported

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

// MakeMidiFile selects a random Bach chorale, transposes its notes, sets a custom BPM, cuts it to 30s, and saves it.
func MakeMidiFile(title string, bpm uint16, transpose int8) error {
	randNum := rand.IntN(372)
	choraleFile := fmt.Sprintf("../motif-find/jsb_chorales_midi/bach_chorale_%d.mid", randNum)

	choraleSMF, err := smf.ReadFile(choraleFile)
	if err != nil {
		return fmt.Errorf("failed to read target midi file %s: %v", choraleFile, err)
	}

	newSMF := smf.New()
	newSMF.TimeFormat = choraleSMF.TimeFormat

	// Calculate maximum ticks for 30 seconds using the built-in Ticks() method
	var maxTicks uint32
	if metric, ok := choraleSMF.TimeFormat.(smf.MetricTicks); ok {
		// Pass the BPM and the desired duration
		maxTicks = metric.Ticks(float64(bpm), 30*time.Second)
	} else {
		// Fallback for SMPTE timecode formats
		maxTicks = uint32(bpm) * 480 / 2
	}

	for _, track := range choraleSMF.Tracks {
		var newTrack smf.Track
		var currentTicks uint32

		// Track active notes to turn them off when cutting the track
		activeNotes := make(map[[2]uint8]bool)

		for _, ev := range track {
			delta := ev.Delta
			msg := ev.Message

			// If this event pushes us past 30 seconds, close all notes and end the track
			if currentTicks+delta >= maxTicks {
				remainingDelta := maxTicks - currentTicks
				firstOff := true

				for chKey := range activeNotes {
					ch, k := chKey[0], chKey[1]
					d := remainingDelta
					if !firstOff {
						d = 0 // Only advance time for the first note-off
					}
					firstOff = false
					newTrack.Add(d, smf.Message(midi.NoteOff(ch, k)))
				}
				break // Stop processing this track
			}

			currentTicks += delta

			var ch, key, vel uint8
			var trackBPM float64

			switch {
			case msg.GetNoteOn(&ch, &key, &vel):
				newKey := clampKey(int(key) + int(transpose))
				msg = smf.Message(midi.NoteOn(ch, newKey, vel))

				if vel > 0 {
					activeNotes[[2]uint8{ch, newKey}] = true
				} else {
					delete(activeNotes, [2]uint8{ch, newKey})
				}

			case msg.GetNoteOff(&ch, &key, &vel):
				newKey := clampKey(int(key) + int(transpose))
				msg = smf.Message(midi.NoteOff(ch, newKey))
				delete(activeNotes, [2]uint8{ch, newKey})

			case msg.GetMetaTempo(&trackBPM):
				if bpm > 0 {
					msg = smf.MetaTempo(float64(bpm))
				}
			}

			newTrack.Add(delta, msg)
		}
		newSMF.Add(newTrack)
	}

	outputPath := fmt.Sprintf("%s.mid", title)
	err = newSMF.WriteFile(outputPath)
	if err != nil {
		return fmt.Errorf("failed to write midi file: %v", err)
	}

	fmt.Printf("Successfully transposed and saved %s -> %s\n", choraleFile, outputPath)
	return nil
}

func clampKey(k int) uint8 {
	if k < 0 {
		return 0
	}
	if k > 127 {
		return 127
	}
	return uint8(k)
}
