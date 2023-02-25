#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[${#BASH_SOURCE[@]} - 1]}")" &> /dev/null && pwd)"
cd "$script_dir"

# for each node (vulnerability) in the JSON
# - join git toplevel to vuln manifest path, minus go.sum
# - go to that directory
# - run 'go get -u <sec vuln pkg name>' et al
# - git commit

top_level=$(git rev-parse --show-toplevel)

gql_result=$(gh api graphql -f=query='
  query{
    repository(name: "golang-workbench", owner: "jlucktay") {
      vulnerabilityAlerts(last: 100, states: OPEN) {
        nodes {
          number
          vulnerableManifestPath
          securityVulnerability {
            severity
            package {
              name
            }
            vulnerableVersionRange
          }
        }
      }
    }
  }
')

list_length=$(jq '.data.repository.vulnerabilityAlerts.nodes | length' <<< "$gql_result")

for ((i = 0; i < list_length; i++)); do
  vuln_number=$(jq --exit-status --raw-output \
    ".data.repository.vulnerabilityAlerts.nodes[$i].number" <<< "$gql_result")

  vuln_path=$(jq --exit-status --raw-output \
    ".data.repository.vulnerabilityAlerts.nodes[$i].vulnerableManifestPath" <<< "$gql_result")

  vuln_pkg=$(jq --exit-status --raw-output \
    ".data.repository.vulnerabilityAlerts.nodes[$i].securityVulnerability.package.name" <<< "$gql_result")

  (
    set -x

    cd "$top_level/${vuln_path%/go.sum}"

    go mod edit -go=1.20
    go get -u -v "$vuln_pkg"
    go mod tidy
    go build ./...
    go clean -x

    git unstage
    git add -- ./go.{mod,sum}

    if ! git commit --message="build(${vuln_path%/go.sum}): upgrade Go to 1.20" > /dev/null; then
      echo "Vuln #$vuln_number affecting '$vuln_pkg' in '${vuln_path%/go.sum}' seems to be resolved already."
    fi
  )

  echo
done
