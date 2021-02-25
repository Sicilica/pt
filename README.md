# pt

## feature ideas

- better summaries
- undo / history
- trends over time (how has categorized work changed over time)
- change summary output format automatically based on input length?
- `pt i` / `pt repl`, opens a repl in which you can run as many commands as you want
  without having to prepend them with `pt`
- `pt edit`, opens a cli-type interface to navigate timeline and make changes to single tasks, such as:
  - change start/stop time (and auto adjust adjacent time as well?)
  - delete
  - merge with adjacent tasks (can this be auto?)
  - split (or, create new one given time bounds?)
  - edit tags (delete is actually equivalent to removing all tags)
- `pt recover` for if you accidentally lose a task you wanted to keep on stack (e.g., via `in` or `drop` or `resume`)
- tracked / untracked time
- breakup summary into at least `pt log` (show actual tasks) and `pt breakdown` (show amount of time spent in different areas)
  - make time period parsing automatic (like PopTime, but PopInterval)
  - default to today (instead of required param)
- remove all time.Now() and store time on command to make simultaneous things simpler / reduce coordination
  - could use autoEnd and the like in more places if this existed
- use sqlx
- `pt update` to easily fetch+build latest version (or, just do `pt-update.sh` on path)
- feature to rotate backup encryption key
