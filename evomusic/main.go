package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/proatscratch/hackathon/evomusic/midioutput"
)

type Organism struct {
	BPM          uint16
	Transpose    int8
	DeletionRate float64
	Fitness      float64
}

func randomOrganism() Organism {
	minBPM, maxBPM := uint16(40), uint16(120)
	bpm := minBPM + uint16(rand.IntN(int(maxBPM-minBPM+1)))
	transpose := int8(rand.IntN(36) - 12)

	return Organism{
		BPM:       bpm,
		Transpose: transpose,
	}
}

func main() {
	numberOrganisms := 4
	totalGenerations := 3

	population := make([]Organism, numberOrganisms)

	editorExecutable := "./text_editor_app"

	// --------------------------------------------------
	// Create initial population
	// --------------------------------------------------

	for i := range population {
		population[i] = randomOrganism()
	}

	// --------------------------------------------------
	// START EDITOR ONCE
	// --------------------------------------------------

	fmt.Println("Starting text editor...")

	editorCmd := exec.Command(editorExecutable)

	stdout, err := editorCmd.StdoutPipe()
	if err != nil {
		log.Fatalf("failed to create editor stdout pipe: %v", err)
	}

	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Start(); err != nil {
		log.Fatalf("failed to start editor: %v", err)
	}

	fmt.Println("Editor started.")

	// --------------------------------------------------
	// CONTINUOUSLY READ EDITOR OUTPUT
	// --------------------------------------------------

	rateCh := make(chan float64)

	go func() {
		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			fmt.Printf("Editor output: %s\n", line)

			rate, err := strconv.ParseFloat(line, 64)
			if err != nil {
				// Ignore anything that isn't a deletion rate.
				continue
			}

			// Send every deletion rate to the main loop.
			rateCh <- rate
		}

		if err := scanner.Err(); err != nil {
			log.Printf("error reading editor stdout: %v", err)
		}

		close(rateCh)
	}()

	// --------------------------------------------------
	// EVOLUTION
	// --------------------------------------------------

	for gen := 1; gen <= totalGenerations; gen++ {

		fmt.Printf("\n=========================================\n")
		fmt.Printf("          GENERATION %d\n", gen)
		fmt.Printf("=========================================\n")

		for i := range population {
			org := population[i]

			fileName := fmt.Sprintf(
				"gen%d_org%d_bpm%d_t%d",
				gen,
				i+1,
				org.BPM,
				org.Transpose,
			)

			midiFileName := fileName + ".mid"

			// Generate MIDI.
			err := midioutput.MakeMidiFile(
				fileName,
				org.BPM,
				org.Transpose,
			)
			if err != nil {
				log.Fatal(err)
			}

			// Start FluidSynth for this organism.
			synthCmd := exec.Command(
				"fluidsynth",
				"-i",
				"soundfont.SF2",
				midiFileName,
			)

			if err := synthCmd.Start(); err != nil {
				log.Printf("failed to play midi: %v", err)
			}

			fmt.Printf(
				"Organism %d [BPM: %d | Transpose: %d]\n",
				i+1,
				org.BPM,
				org.Transpose,
			)

			fmt.Println("Waiting for deletion rate...")

			// --------------------------------------------------
			// Wait for the next deletion rate from the editor.
			// The editor itself remains running.
			// --------------------------------------------------

			deletionRate, ok := <-rateCh
			if !ok {
				log.Fatal("editor stopped producing output")
			}

			fmt.Printf(
				"Received deletion rate: %.2f\n",
				deletionRate,
			)

			population[i].DeletionRate = deletionRate
			population[i].Fitness = 1000.0 - deletionRate

			// Stop FluidSynth for this organism.
			if synthCmd.Process != nil {
				if err := synthCmd.Process.Kill(); err != nil {
					log.Printf("failed to stop fluidsynth: %v", err)
				}

				_ = synthCmd.Wait()
			}
		}

		// --------------------------------------------------
		// SELECTION
		// --------------------------------------------------

		sort.Slice(population, func(i, j int) bool {
			return population[i].Fitness > population[j].Fitness
		})

		fmt.Println("\n--- Selection Results (Ranked) ---")

		for i, organism := range population {
			fmt.Printf(
				"Rank %d: BPM %d | Transpose %d | Fitness: %.2f\n",
				i+1,
				organism.BPM,
				organism.Transpose,
				organism.Fitness,
			)
		}

		if gen == totalGenerations {
			fmt.Printf("\n=========================================\n")
			fmt.Printf(
				"Evolution Complete! Optimal -> BPM: %d | Transpose: %d\n",
				population[0].BPM,
				population[0].Transpose,
			)
			fmt.Printf("=========================================\n")
			break
		}

		// --------------------------------------------------
		// BREEDING
		// --------------------------------------------------

		fmt.Println("\n--- Breeding Generation", gen+1, "---")

		nextGen := make([]Organism, numberOrganisms)

		// Elitism.
		nextGen[0] = Organism{
			BPM:       population[0].BPM,
			Transpose: population[0].Transpose,
		}

		fmt.Printf(
			"Slot 1 (Elitism): Saved BPM %d | Transpose %d\n",
			nextGen[0].BPM,
			nextGen[0].Transpose,
		)

		parent1 := population[0]
		parent2 := population[1]

		for i := 1; i < numberOrganisms; i++ {

			var childBPM uint16
			var childTranspose int8

			// Crossover BPM.
			if rand.N(2) == 0 {
				childBPM = parent1.BPM
			} else {
				childBPM = parent2.BPM
			}

			// Crossover transpose.
			if rand.N(2) == 0 {
				childTranspose = parent1.Transpose
			} else {
				childTranspose = parent2.Transpose
			}

			// BPM mutation.
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

			// Transpose mutation.
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

			fmt.Printf(
				"Slot %d (Child): BPM %d | Transpose %d\n",
				i+1,
				childBPM,
				childTranspose,
			)

			nextGen[i] = Organism{
				BPM:       childBPM,
				Transpose: childTranspose,
			}
		}

		population = nextGen
	}

	// --------------------------------------------------
	// CLEANUP
	// --------------------------------------------------

	fmt.Println("Evolution complete. Closing editor...")

	if editorCmd.Process != nil {
		if err := editorCmd.Process.Kill(); err != nil {
			log.Printf("failed to stop editor: %v", err)
		}

		_ = editorCmd.Wait()
	}
}
