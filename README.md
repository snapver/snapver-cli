![demo](https://github.com/user-attachments/assets/778dc230-1deb-49af-985f-add0c53b37dd)

# `snapver-cli`

Snapver is a tool that automatically records file changes. It aims to leave at least a restorable snapshot during automated workflows.

## Introduction

Recently, I’ve seen more people using automated tools like Claude Code or Gemini CLI.
While Git can already handle those changes, I built Snapver to make it much easier to isolate, track, and manage them within those workflows.

> [!NOTE]  
> If you already manage your Git history with care, this is a “useless” tool.

## Guide

### Installation

```sh
# Homebrew
brew install snapver/tap/snapver-cli
```

### Usage

```
Usage: snapver-cli [command]

Available Commands:
  start       Start tracking changes
  stop        Stop tracking
  ok          Merge changes into default branch
  clear       Delete snapver branches and diff data.
  version     Show current version
```

### Guide

1. Before using tools like Claude Code or Gemini CLI to auto-modify your code, run `snapver-cli start`.
2. After the automated changes are complete, use `snapver-cli ok` to keep the changes, or `snapver-cli clear` to discard them.

> [!WARNING]  
> Only run `snapver-cli start` for automated workflows that need to be tracked. Since it records all changes in detail, do not manually edit files while it’s running.

## Goal

The ultimate goal of this tool is to help people who want to automate everything (even version control) or don't know Git well yet. Once they realize the importance of proper version control, they should throw this tool in the trash.

## Under the hood

Snapver is a Git wrapper. It isolates a working branch, tracks all changes, and helps you return to your original branch with simpler, more intuitive commands.

- `snapver-cli start`: Creates a new branch and begins tracking changes
- `snapver-cli clear`: Discards all changes and switches back to the default branch
- `snapver-cli ok`: Merges changes into the default branch and switches back

## License

MIT licensed. Free forever.
