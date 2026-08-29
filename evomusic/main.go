package main

import (
	"fmt"
	"math/rand"

	"github.com/proatscratch/hackathon/evomusic/midioutput"
	"github.com/proatscratch/hackathon/evomusic/motif"
	"gitlab.com/gomidi/midi/v2"
)

func randomMotif(motifs []motif.Motif) motif.Motif {
	randomNumber := rand.Intn(len(motifs) + 1)
	return motifs[randomNumber]
}

func main() {
	primaryMotif := randomMotif(motif.Motifs)
	secondaryMotif := randomMotif(motif.Motifs)

	song := append(primaryMotif.MotifContent, secondaryMotif.MotifContent...)
	song = append(song, 127)
	song = append(song, primaryMotif.MotifContent...)

	fmt.Println(song)

	midioutput.MakeMidiFile(midi.C(5).Value(), song)
}
