# Design

- Recursive crawl func
- Get `monzo.com` and child links
  - Store `monzo.com` in map, with all child links as a sub-element slice
  - `go func crawl()` on each child link
  - Rinse, repeat
  - Won't double up because the map will keep track of where we've been
- Function to spit out map in a nice-looking format, with parent and child URLs

## Nice to have

- Add some kind of CLI filter, to filter URLs and see only their links
  - Can achieve this with the current JSON output using some `jq` wizardry
- Double check `crawl()` for potential stack overflow(s), as it is recursive
- Tidier `main()` func, as a lot of it is still dedicated to log/output setup
- Stream JSON output realtime, rather than one big marshal operation at the end
- Nail down some more edge cases in the ever-growing `convertURL()` func
  - Alternatively, source an external library to handle URL parsing
