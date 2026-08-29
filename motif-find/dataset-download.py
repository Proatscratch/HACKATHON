from music21 import corpus
import os

# Create a folder to store the MIDI files
output_dir = "jsb_chorales_midi"
os.makedirs(output_dir, exist_ok=True)

print("Extracting JSB Chorales...")

successful = 0
failed = 0

# corpus.chorales.Iterator() contains all ~380+ Bach chorales
for i, chorale in enumerate(corpus.chorales.Iterator()):
    filename = f"bach_chorale_{i+1}.mid"
    filepath = os.path.join(output_dir, filename)
    
    try:
        # Attempt to export the chorale as a MIDI file
        chorale.write('midi', fp=filepath)
        successful += 1
    except Exception as e:
        # If a repeat ExpanderException occurs, skip it and continue
        print(f"Skipped chorale {i+1} due to complex repeats.")
        failed += 1
    
print(f"\nExtraction complete. Successfully saved: {successful}. Skipped: {failed}.")