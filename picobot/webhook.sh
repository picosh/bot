#!/usr/bin/env bash

script='mutation {
  createWebhook(config: {
    url: \"https://picobot.erock.io/push\"
    events: [REPO_UPDATE]
    query: \"\"\"
        query {
            webhook {
                uuid
                event
                date
                ... on RepositoryEvent {
                    repository {
                        name,
                        revparse_single(revspec: \"HEAD\") {
                            id,
                            shortId,
                            author { name },
                            committer { name },
                            message
                        }
                    }
                }
            }
        }
    \"\"\"
  }) { id }
}'

# to confirm it worked
# script='query {
# 	userWebhooks { results { id, url } }
# }'

# delete webhook
# script='mutation {
#     deleteWebhook(id: 13) { id }
# }'

script="$(echo $script)" # the query should be a one-liner, without newlines
curl -i \
    -H 'Content-Type: application/json' \
    --oauth2-bearer "$PAT" \
    -X POST -d "{ \"query\": \"$script\"}" \
    https://git.sr.ht/query
