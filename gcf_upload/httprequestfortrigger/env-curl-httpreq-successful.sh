# https://stackoverflow.com/questions/13341955/how-to-pass-a-variable-in-a-curl-command-in-shell-scripting
curl -m 70 -X POST "${URL}" \
  -H "Authorization: bearer $(gcloud auth print-identity-token)" \
  -H "Content-Type: application/json" \
  -d @sample-interest-valid.json
