# dotfiles

The aim of this `github/dotfiles` project is to:

- parse this [dotfiles GitHub](https://dotfiles.github.io) site
- collect all links to GitHub repos
- see which language is most-used in each repo
- return the repos that are primarily Go-based

Internal design follows the pipelines pattern.

## Notes

I know I'm taking the long way round by doing things this way, but it's all for a good cause: learning.

The GitHub API does offer [searching by language](https://help.github.com/articles/searching-repositories/#search-by-language) and the [GET on /repos](https://developer.github.com/v3/repos/#get) does offer the primary language at the top level without drilling down into `/languages` but where's the ~~education~~ fun in not doing it yourself? :)
