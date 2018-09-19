# Design

- Recursive crawl func
- Get `monzo.com` and child links
  - Store `monzo.com` in map, with all child links as a sub-element slice
  - `go func crawl()` on each child link
  - Rinse, repeat
  - Won't double up because the map will keep track of where we've been
- Function to spit out map in a nice-looking, with parent and child URLs
  - Add some kind of CLI filter, to parse out URLs and see links?
