#!/usr/bin/env bash
set -euo pipefail

# for each node (vulnerability) in the JSON
# - join git toplevel to vuln manifest path, minus go.mod
# - go to that directory
# - run 'go get -u <sec vuln pkg name>' et al
# - git commit

if ! top_level=$(git rev-parse --show-toplevel); then
  echo >&2 "Not running from inside a git repo."
  exit 1
fi

origin_url=$(git remote get-url origin)

current_project=$(cut -d '/' -f 5 <<< "$origin_url")
current_project=${current_project/%'.git'/''}

project_owner=$(cut -d '/' -f 4 <<< "$origin_url")

gql_result=$(gh api graphql -f=query="
  query{
    repository(name: \"$current_project\", owner: \"$project_owner\") {
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
          securityAdvisory {
            notificationsPermalink
          }
        }
      }
    }
  }
")

list_length=$(jq '.data.repository.vulnerabilityAlerts.nodes | length' <<< "$gql_result")

for ((i = 0; i < list_length; i++)); do
  vuln_number=$(jq --exit-status --raw-output \
    ".data.repository.vulnerabilityAlerts.nodes[$i].number" <<< "$gql_result")

  vuln_path=$(jq --exit-status --raw-output \
    ".data.repository.vulnerabilityAlerts.nodes[$i].vulnerableManifestPath" <<< "$gql_result")

  vuln_pkg=$(jq --exit-status --raw-output \
    ".data.repository.vulnerabilityAlerts.nodes[$i].securityVulnerability.package.name" <<< "$gql_result")

  vuln_link=$(jq --exit-status --raw-output \
    ".data.repository.vulnerabilityAlerts.nodes[$i].securityAdvisory.notificationsPermalink" <<< "$gql_result")

  tlvp="$top_level/${vuln_path/%'go.mod'/''}"

  this_module=$(git -C "$tlvp" rev-parse --show-prefix)
  this_module=${this_module%/}

  if [[ -n $this_module ]]; then
    this_module="($this_module)"
  fi

  (
    set -x

    cd "$tlvp"

    go get -u -v "$vuln_pkg"
    go mod tidy
    go build ./...
    go clean -x

    git unstage
    git add ./go.{mod,sum}

    if ! git commit \
      --message="build$this_module: address Dependabot vuln $vuln_number" \
      --message="URL: $vuln_link" \
      > /dev/null; then
      echo "Vuln #$vuln_number affecting '$vuln_pkg' in '${vuln_path%/go.sum}' seems to be resolved already."
    fi
  )

  echo
done
