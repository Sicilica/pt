# pt

## Installation

Make sure you have Go installed, then just run `make` and add the resulting `bin/pt` to your PATH.
That's it!

## Getting Started

`pt` use tags to track blocks of time called tasks.
To start a new task, just use `pt in` with some tags that describe what you're doing:

```sh
# Starts a task with tag "slack"
pt in slack
# Starts a task with tags "cooking" and "lunch"
pt in cooking lunch
#  Starts a wask with tag "bug-123"
pt in bug-123
```

When you finish what you're working on, you can just start another task right away with `pt in`.
This will automatically close the current task and start a new one.
If you don't know what you'll be doing next, you can use `pt out` instead:
```sh
# Start working on a ticket
pt in bug-123
# Stop working on the ticket to attend a meeting
pt in meeting
# Finish the meeting and decide what to do next
pt out
```

If you ever forget what you were working on, use `pt status`:
```sh
# What task is active right now?
pt status
> doing [bug-123] for 5m since 10:20
```

Later, you can use `pt summary` to recall your tracked tasks:
```sh
# List today's recorded tasks
pt summary today
> 8:30 - 10m [email]
> 8:40 - 1h15m [bug-123]
> 9:55 - 20m [meeting]
> 10:15 - 1h [bug-123]
> -----
> total: 2h45m

# List tasks for the whole week
pt summary this week

# How much time have I spent in meetings this month?
pt summary this month meeting
```

That's it! For a full list of commands, use `pt help`, or you can keep reading to learn about some of the more powerful things you can do.

## Advanced Usage

### Forgetting to Clock In

If you forgot to start a task, you can use `pt rollback`:

```sh
# Started working on foobar at 15:00
pt rollback 15:00 foobar
# Started working on foobar 5 minutes ago
pt rollback 5m foobar
```

`pt rollback` is just like you had run `pt in` at the specified time.

Additionally, if you used `pt out` but then realized you actually meant to start a new task, you can leave out the time to automatically rollback to the time that the last task ended:

```sh
# While working on task A...
pt in A

# We actually stopped working on task A 5 minutes ago
pt out 5m
# And, we started working on task B at that exact time
pt rollback B
```

### Forgetting to Clock Out

If you forgot to end a task, you can provide a time to `pt out`:

```sh
# End active task at 15:00
pt out 15:00
# End active task 5 minutes ago
pt out 5m 
```

Or, if you want to start another task in its place, use `pt rollback`:

```sh
# End active task 5 minutes ago and replace with "foobar"
pt rollback 5m foobar
```

### Editing the Active Task

If you realize you forgot to set the correct tags on your task, you can edit these tags with `pt add` and `pt remove`:
```sh
# Start working on "interview"
pt in interview
# Now, tags will be both "interview" and "meeting"
pt add meeting
# On second thought, it wasn't an "interview" after all
pt remove interview

# If we check now, our task only has "meeting"
pt summary
```

_Hint: Unlike `pt split`, `pt add` and `pt remove` edit the active task, so changes apply retroactively to the time when the task first started. If you want to track the time when these tags changed, see [Multitasking](#multitasking)._

### Tag Groups

While there is no limit on how many tags you can apply to each task, it can be burdensome to type them all out if you want to use a lot of tags.
Tag groups can help with this:
```sh
# Interviews should always count as meetings
pt describe interview meeting
# Lunch is a type of break
pt describe lunch break

# This task is only tagged with "interview", but it will show up when we look at time spent in "meeting" tasks
pt in interview
```

For tags that shouldn't belong to a group, there is no need to call `pt describe`.
You can always start using tags right away.
And if you later decide a tag _should_ belong to a group, you can add one with `pt define` and it will apply retroactively, too.

### Multitasking

Use `pt split` when you want to start a new copy of the current task.
This can be useful for narrowing down the exact thing you're working on, or for working on multiple things at the same time.

_Hint: Unlike `pt add` and `pt remove`, `pt split` creates a new task and records new in/out times._

```sh
# Working on "bug-123"
pt in bug-123
# Splits the task; new tags are "bug-123" and "fix-tests"
pt split add fix-tests

# Tags are "cooking" and "dinner"
pt in cooking dinner
# Finished cooking, but still eating
pt split remove cooking
```

An optional start time can be provided, similar to `rollback`:
```sh
# Stopped doing "research" 5 minutes ago, but continuing to work on other tags
pt split 5m remove research
```

### Context Switching

Sometimes you need to temporarily stop working on a task, but you don't want to have to type the full task definition out again later.
All of these tags can help with that:
```sh
# Stop the active task for now, but save it for later
# (similar to `pt out`)
pt pause
# Suspend the active task, then start a new one with tag "meeting"
# (similar to `pt in`)
pt swap meeting

# Continue working on the most recently suspended task
pt resume
# Resume the suspended task at index #1
# (use `pt status` to check your currently suspended tasks)
pt resume 1
# Suspend the active task, and resume your most recently suspended task
pt swap
# Suspend the active task, and resume task #2
pt swap 2

# Actually, suspended task #1 isn't relevant anymore
pt drop 1
```

### Backup and Sync

TODO...

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
    - may actually need PopDate instead of PopTime, and have inferred 00:00 and 23:59 for the start/end
  - default to today (instead of required param)
- remove all time.Now() and store time on command to make simultaneous things simpler / reduce coordination
  - could use autoEnd and the like in more places if this existed
- use sqlx
- `pt update` to easily fetch+build latest version (or, just do `pt-update.sh` on path)
- feature to rotate backup encryption key
- `pt history` to show actual tasks (similar to `pt log`?), in a dynamic way (no need to specify a time period)
- `pt reopen`; like rollback, but always uses the last closed time, and keeps the exact same tags without you having to specify. (I have mixed feelings on this...)

Observation: with better editing support (esp. the ability to quickly add a bunch of already-ended tasks), I would be more willing to keep high tracking overall. Currently, if I fail to track some things, then I know I can't start tracking my new thing without making it impossible to add my old ones - and I also know that adding tracking for all my old things will take a non-trivial amount of effort/time.

(The other factor that has reduced tracking has been the multi-computer setup with not-always-robust syncing, i.e. "I don't really want to boot up the other comp and make it sync right now" syndrome. Luckily, it running slow as balls on my mac has fixed this kek)
