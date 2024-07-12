# Ginsu 🔪

> Cut through the noise.

*GIthub Notifications - Sift Unread* (or `ginsu` for short) is a small tool that will sift through your unread GitHub
notifications, and where an issue or PR is already closed, mark the notification as done. ✅

## Usage

Export a PAT with `repo` and `notifications` scopes to the `GITHUB_TOKEN` env var.

There is a `--owner-allowlist` flag that (if set) will limit which notifications are affected.

See `ginsu --help` for more details.
