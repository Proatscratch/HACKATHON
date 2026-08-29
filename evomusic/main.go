package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"runtime"
	"sort"

	"github.com/proatscratch/hackathon/evomusic/midioutput"
)

type Organism struct {
	BPM          uint16
	Transpose    int8
	DeletionRate int
	Fitness      int
}

func randomOrganism() Organism {
	minBPM, maxBPM := uint16(40), uint16(120)
	bpm := minBPM + uint16(rand.IntN(int(maxBPM-minBPM+1)))
	transpose := int8(rand.IntN(36) - 12)
	return Organism{BPM: bpm, Transpose: transpose}
}

// playMidi invokes the default system MIDI player
func playMidi(filename string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("afplay", filename)
	case "windows":
		// 'start /wait' opens the default player and waits for it to finish
		cmd = exec.Command("cmd", "/c", "start", "/wait", filename)
	case "linux":
		// timidity is the most common CLI midi player for linux
		cmd = exec.Command("timidity", filename)
	default:
		fmt.Println("Audio playback not supported on this OS. Please open the file manually.")
		return
	}

	fmt.Printf("▶️ Playing %s...\n", filename)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Playback interrupted or player not found: %v\n", err)
	}
}

func main() {
	numberOrganisms := 4
	totalGenerations := 3
	population := make([]Organism, numberOrganisms)

	for i := range population {
		population[i] = randomOrganism()
	}

	for gen := 1; gen <= totalGenerations; gen++ {
		fmt.Printf("\n=========================================\n")
		fmt.Printf("          GENERATION %d\n", gen)
		fmt.Printf("=========================================\n")

		for i := range population {
			org := population[i]
			fileName := fmt.Sprintf("gen%d_org%d_bpm%d_t%d", gen, i+1, org.BPM, org.Transpose)

			err := midioutput.MakeMidiFile(fileName, org.BPM, org.Transpose)
			if err != nil {
				log.Fatal(err)
			}

			// Play the generated file so the user can evaluate it
			playMidi(fileName + ".mid")

			fmt.Printf("Organism %d [BPM: %d | Transpose: %d] -> Enter deletion rate (0-100, lower is better): ", i+1, org.BPM, org.Transpose)

			var deletionRate int
			_, err = fmt.Scanf("%d", &deletionRate)
			if err != nil {
				log.Fatalf("Invalid input received: %v", err)
				os.Exit(1)
			}

			population[i].DeletionRate = deletionRate
			population[i].Fitness = 1000 - deletionRate
		}

		sort.Slice(population, func(i, j int) bool {
			return population[i].Fitness > population[j].Fitness
		})

		fmt.Println("\n--- Selection Results (Ranked) ---")
		for i, organism := range population {
			fmt.Printf("Rank %d: BPM %d | Transpose %d | Fitness: %d\n", i+1, organism.BPM, organism.Transpose, organism.Fitness)
		}

		if gen == totalGenerations {
			fmt.Printf("\n=========================================\n")
			fmt.Printf("Evolution Complete! Optimal -> BPM: %d | Transpose: %d\n", population[0].BPM, population[0].Transpose)
			fmt.Printf("=========================================\n")
			break
		}

		fmt.Println("\n--- Breeding Generation", gen+1, "---")
		nextGen := make([]Organism, numberOrganisms)

		nextGen[0] = Organism{BPM: population[0].BPM, Transpose: population[0].Transpose}
		fmt.Printf("Slot 1 (Elitism): Saved BPM %d | Transpose %d\n", nextGen[0].BPM, nextGen[0].Transpose)

		parent1 := population[0]
		parent2 := population[1]

		for i := 1; i < numberOrganisms; i++ {
			var childBPM uint16
			var childTranspose int8

			if rand.N(2) == 0 {
				childBPM = parent1.BPM
			} else {
				childBPM = parent2.BPM
			}

			if rand.N(2) == 0 {
				childTranspose = parent1.Transpose
			} else {
				childTranspose = parent2.Transpose
			}

			if rand.N(100) < 25 {
				newBPM := int32(childBPM) + int32(rand.IntN(17)-8)
				if newBPM < 40 {
					newBPM = 40
				}
				if newBPM > 160 {
					newBPM = 160
				}
				childBPM = uint16(newBPM)
			}

			if rand.N(100) < 25 {
				newTranspose := childTranspose + int8(rand.IntN(7)-3)
				if newTranspose < -24 {
					newTranspose = -24
				}
				if newTranspose > 24 {
					newTranspose = 24
				}
				childTranspose = newTranspose
			}

			fmt.Printf("Slot %d (Child): BPM %d | Transpose %d\n", i+1, childBPM, childTranspose)
			nextGen[i] = Organism{BPM: childBPM, Transpose: childTranspose}
		}
		population = nextGen
	}
}
