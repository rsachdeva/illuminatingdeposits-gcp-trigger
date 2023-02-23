# https://cloud.google.com/functions/docs/running/calling#cloudevent-function-curl-tabs-pubsub
curl localhost:8080 \
      -X POST \
      -H "Content-Type: application/json" \
      -H "ce-id: 123451234512345" \
      -H "ce-specversion: 1.0" \
      -H "ce-time: 2020-01-02T12:34:56.789Z" \
      -H "ce-type: google.cloud.pubsub.topic.v1.messagePublished" \
      -H "ce-source: //pubsub.googleapis.com/projects/illuminatingdeposits-gcp/topics/deltaanalyticstopic" \
      -d '{
            "message": {
              "data": "eyJoaWdoZXN0X2RlbHRhX2RlcG9zaXRzX2J5X2RhdGUiOiJbZGVsdGEgZGVwb3NpdHNCeURhdGVdXG5bMTQ3LjM3IFdlZCwwMi8yMi8yMF1cblsxODcwLjU3IFdlZCwwMi8yMS8yMF1cblsxODcwLjU3IFdlZCwwMi8yMC8yMF1cblsxNDcuMzcgV2VkLDAyLzE5LzIwXVxuWzE4NzAuNTcgV2VkLDAyLzE4LzIwXVxuIn0=",
              "attributes": {
                 "googclient_schemaencoding":"JSON",
                 "googclient_schemarevisionid":"1fd566bd"
              }
            },
            "subscription": "projects/illuminatingdeposits-gcp/subscriptions/deltaanalyticstestsubscription"
          }'