import os
from collections import Counter
from music21 import converter, interval, note, stream

def get_interval_sequence(part):
    intervals = []
    notes = list(part.recurse().getElementsByClass(note.Note))
    for i in range(len(notes) - 1):
        inter = interval.Interval(noteStart=notes[i], noteEnd=notes[i+1])
        intervals.append(inter.semitones)
    return intervals

def get_ngrams(sequence, n):
    ngrams = []
    for i in range(len(sequence) - n + 1):
        ngrams.append(tuple(sequence[i:i+n]))
    return ngrams

def analyze_dataset_motifs(midi_directory, n_gram_size=4, top_k=500):
    all_ngrams = []
    for filename in os.listdir(midi_directory):
        if filename.endswith(".mid") or filename.endswith(".midi"):
            filepath = os.path.join(midi_directory, filename)
            try:
                score = converter.parse(filepath)
                for part in score.parts:
                    intervals = get_interval_sequence(part)
                    part_ngrams = get_ngrams(intervals, n_gram_size - 1)
                    all_ngrams.extend(part_ngrams)
            except Exception:
                pass # Silently skip errors for a cleaner console

    motif_counts = Counter(all_ngrams)
    return motif_counts.most_common(top_k)

def save_motif_as_midi(interval_sequence, filepath, start_pitch='C4'):
    """Converts a list of intervals back into notes and saves as MIDI."""
    motif_stream = stream.Stream()
    
    # Create the starting note
    current_note = note.Note(start_pitch)
    motif_stream.append(current_note)
    
    # Generate the rest of the notes based on the intervals
    for semitones in interval_sequence:
        next_note = current_note.transpose(semitones)
        motif_stream.append(next_note)
        current_note = next_note
        
    # Export to file
    motif_stream.write('midi', fp=filepath)

# --- Run the Analysis ---
dataset_path = "./jsb_chorales_midi"

# 1. Create a folder to put the new motif MIDI files
output_dir = "./extracted_top_motifs"
os.makedirs(output_dir, exist_ok=True)

# 2. Analyze the dataset
print("Analyzing dataset... this may take a moment.")
top_motifs = analyze_dataset_motifs(dataset_path, n_gram_size=5, top_k=100)

# 3. Export the top motifs
print(f"Exporting top {len(top_motifs)} motifs to '{output_dir}'...")
for rank, (motif_intervals, count) in enumerate(top_motifs, 1):
    filename = f"rank_{rank}_count_{count}.mid"
    filepath = os.path.join(output_dir, filename)
    
    save_motif_as_midi(motif_intervals, filepath)
    print(f"Saved #{rank} (Appeared {count} times): {filename}")