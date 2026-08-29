Build something to help students have healthy habits
# HACKATHON
ONE OF A KIND!
Design Document
NOTE: plan must be NOT capitalized.

ACTS (Formulated by Anish Deeduvanu)
You may implement behavior within your own discretion.
You may not delete a bunch (LOOKING AT YOU VISMAI)  of code OR ADD OR REMOVE behavior without 2/3d votes.
	
What counts as implementing behavior:
Writing code for that behavior
Removing code that causes issues with that behavior, so long as you make sure other behavior behaves the same.
What doesn’t count as implementing behavior:
Performing what would count as implementing a bunch of behaviors at the same time in an unnatural way. (such as rewriting code that doesn’t need to be rewritten)
		Code that does not implement voted on behavior.
Any Policy not listed inside this document requires a 2/3 vote to be banned! Otherwise it is given automatically to everyone.
Team Synthesis Requirements (Formulated by Vismai Nair)
Brainstorming concludes at most 20 minutes after theme drop
5 minutes for thinking, then convene and share ideas from worst to best
Every 15 mins we all convene and share progress from what we think we did most to what we think we did least. (In addition, Behavior addition will be done here: Anish).

GitHub structure (Formulated by Vismai Nair)
Monorepo structure: All the code is in one repo regardless of whether it interacts or not.
Atomic commits under the Conventional Commits structure.
Features are like feat(thing that ur fixing): what you did
Docs are like docs(thing): what the doc is
Initialization goes under init:
Refactor goes under refactor(thing): thing
Performance optimizations go under perf(thing): thing
EVERY CHANGE NEEDS A GIT COMMIT LOOKING AT YOU ANISH
Example: 
git commit -m “feat(evol-alg): evolutionary algorithm for music”
git commit -m “docs(sim): document the sim”
git commit -m “refactor(evol-alg): extract line of code to vars”
Vismai creates the github page
Push every commit to remote origin
Essentially, the workflow:

mkdir (projectname)
cd (projectname)
git remote add origin https://github.com/vismainair/(projectname)
git commit -m “feat(evol-alg): evolutionary alg for music”
git branch -M main
git push -u origin main
git pull
Note I (Vismai) will probably force you all to install taskfile to do this for you \\ (L1br@ryUs3r)

Plan (To be done in hackathon)
All Behavior voting results should be done here.




# evomusic

Evomusic is the product and experience. This repository contains the open-source code for the Evomusic project.

## Overview

Evomusic is an adaptive study companion built to help students maintain focus during long study sessions. It responds to changes in attention, stress, and performance by adjusting the sound environment to support calmer, more productive work.

The website is secondary to the core idea: it exists to introduce, entertain, and explain the Evomusic experience. The real focus is the concept of adaptive music for study flow.

## Purpose

This project communicates the Evomusic vision through a minimal, editorial interface and a music-driven visual language. The goal is to entertain the audience with the product concept first, while the website simply presents it.

## Credits

- Vismai Nair
- Eklavya Saini
- Anish Deeduvanu

## Notes

- This is an open-source front-end project built with plain C++, Go, and Python.
- The website is a presentation layer for Evomusic, not the main product itself.
- This project is intended to run on macOS only.
- You must build from source; there are no releases
