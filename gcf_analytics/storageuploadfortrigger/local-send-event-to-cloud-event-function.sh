# https://cloud.google.com/functions/docs/running/calling#cloudevent-function-curl-tabs-storage
curl localhost:8080 \
  -X POST \
  -H "Content-Type: application/json" \
  -H "ce-id: 123451234512345" \
  -H "ce-specversion: 1.0" \
  -H "ce-time: 2020-01-02T12:34:56.789Z" \
  -H "ce-type: google.cloud.storage.object.v1.finalized" \
  -H "ce-source: //storage.googleapis.com/projects/_/buckets/illuminating_upload_json_bucket_output" \
  -H "ce-subject: objects/interestresponse.json" \
  -d '{
        "bucket": "illuminating_upload_json_bucket_output",
        "contentType": "application/json",
        "kind": "storage#object",
        "name": "interestresponse.json",
        "storageClass": "STANDARD"
      }'