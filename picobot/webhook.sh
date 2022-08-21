#!/usr/bin/env bash

script='mutation {
  createWebhook(config: {
    url: "https://picobot.erock.io/push"
    events: [REPO_UPDATE]
    query: """
        query {
            webhook {
                uuid
                event
                date
                ... on RepositoryEvent {
                    repository {
                        name,
                        revparse_single(revspec: "HEAD") {
                            shortId,
                            author { name },
                            committer { name },
                            message
                        }
                    }
                }
            }
        }
    """
  }) { id }
}'
script="$(echo $script)"   # the query should be a one-liner, without newlines
curl -i \
    -H 'Content-Type: application/json' \
    --oauth2-bearer "$PAT" \
    -X POST -d "{ \"query\": \"$script\"}" \
    https://git.sr.ht/query
