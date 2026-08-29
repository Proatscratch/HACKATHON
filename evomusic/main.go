package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"sort"
	"time"

	"github.com/proatscratch/hackathon/evomusic/midioutput"
)

// Organism represents an individual in our population with both BPM and Transpose traits
type Organism struct {
	BPM          uint16
	Transpose    int8
	DeletionRate int
	Fitness      int
}

// randomOrganism generates completely random starting traits
func randomOrganism() Organism {
	minBPM, maxBPM := uint16(40), uint16(120)
	bpm := minBPM + uint16(rand.IntN(int(maxBPM-minBPM+1)))

	// Random transpose between -12 and +12
	transpose := int8(rand.IntN(36) - 12)

	return Organism{BPM: bpm, Transpose: transpose}
}

func main() {
	numberOrganisms := 4
	totalGenerations := 3
	population := make([]Organism, numberOrganisms)

	// Initialize Generation 1
	for i := range population {
		population[i] = randomOrganism()
	}

	for gen := 1; gen <= totalGenerations; gen++ {
		fmt.Printf("\n=========================================\n")
		fmt.Printf("          GENERATION %d\n", gen)
		fmt.Printf("=========================================\n")

		// 1. Generate MIDI files and collect user feedback
		for i := range population {
			org := population[i]

			// Generate the file for this specific organism so the user can listen to it
			fileName := fmt.Sprintf("gen%d_org%d_bpm%d_t%d", gen, i+1, org.BPM, org.Transpose)
			err := midioutput.MakeMidiFile(fileName, org.BPM, org.Transpose)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Organism %d [BPM: %d | Transpose: %d] -> Enter deletion rate (0-100, lower is better): ", i+1, org.BPM, org.Transpose)

			time.Sleep(200 * time.Millisecond)

			var deletionRate int
			_, err = fmt.Scanf("%d", &deletionRate)
			if err != nil {
				log.Fatalf("Invalid input received: %v", err)
				os.Exit(1)
			}

			population[i].DeletionRate = deletionRate
			population[i].Fitness = 1000 - deletionRate
		}

		// 2. Sort by fitness descending
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

		// 3. Breed the next generation
		fmt.Println("\n--- Breeding Generation", gen+1, "---")
		nextGen := make([]Organism, numberOrganisms)

		// Elitism
		nextGen[0] = Organism{BPM: population[0].BPM, Transpose: population[0].Transpose}
		fmt.Printf("Slot 1 (Elitism): Saved BPM %d | Transpose %d\n", nextGen[0].BPM, nextGen[0].Transpose)

		parent1 := population[0]
		parent2 := population[1]

		for i := 1; i < numberOrganisms; i++ {
			var childBPM uint16
			var childTranspose int8

			// Uniform Crossover for BPM
			if rand.N(2) == 0 {
				childBPM = parent1.BPM
			} else {
				childBPM = parent2.BPM
			}

			// Uniform Crossover for Transpose (independent from BPM)
			if rand.N(2) == 0 {
				childTranspose = parent1.Transpose
			} else {
				childTranspose = parent2.Transpose
			}

			// Mutation: BPM
			if rand.N(100) < 25 {
				mutationAmount := int16(rand.IntN(17) - 8)
				newBPM := int32(childBPM) + int32(mutationAmount)
				if newBPM < 40 {
					newBPM = 40
				}
				if newBPM > 160 {
					newBPM = 160
				}
				childBPM = uint16(newBPM)
			}

			// Mutation: Transpose
			if rand.N(100) < 25 {
				mutationAmount := int8(rand.IntN(7) - 3) // Mutate by -3 to +3 semitones
				newTranspose := childTranspose + mutationAmount

				// Keep bounds safe (e.g., +/- 2 octaves max)
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
