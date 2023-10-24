# Utilities for galqiwi's personal archive

### Features:

- `garc limbo` personal temporary file storage, requires remote directory with ssh access
- `garc update` self update
- `garc config edit` simple self-config editor, respects `EDITOR` env

### Principles
- Static binary -- it should be runnable on all x86 linux installations. One should compule this code with `CGO_ENABLED=0`.
- Minimal assumptions about os -- this code should prefer goland dependencies over os dependencies. For example, it should embed ssh client over using `ssh` command that can be unavailable.

### Why?

This project has two purposes:
- It is a single entrypoint for all scripts and automation I need for my archival needs
- It is an accurate representation of how I can code, which is useful for my portfolio
