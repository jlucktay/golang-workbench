# Redis: Zero to Master in 30 minutes

From [here](https://www.openmymind.net/2011/11/8/Redis-Zero-To-Master-In-30-Minutes-Part-1/) and
[here](https://www.openmymind.net/2011/11/8/Redis-Zero-To-Master-In-30-Minutes-Part-2/).

## Order of operations

**NOTE:** failure at any point will simply exit with a status of zero

1. Looks for the required `OPENMYMIND_RSS_URL` environment variable.
1. Connects to Redis at `localhost:6379`.
1. Fetches the RSS feed from the URL set in `OPENMYMIND_RSS_URL`.
1. Decodes the XML in the RSS feed.
1. Ranges across every item found in the feed:
    1. Converts the entire item to JSON.
    1. Converts the publication date to a UNIX timestamp.
    1. Sets the value of Redis key `item:<item guid>` to the entire JSON payload.
    1. Adds the above key to a sorted set where the score is the publication date.
1. Gets the keys of the three most recent items from the sorted set.
1. For each of these keys:
    1. Gets the JSON payload back out of Redis.
    1. Decodes the JSON.
    1. Prints the item's title and publication date.
