# Illuminating Deposits - GCP Trigger Applied to Cloud Functions

Illuminating Deposits Google Cloud Platform GCP based trigger of Resources applied to Cloud Functions

(Development is WIP)

<p align="center">
<img src="./logo.png" alt="Illuminating Deposits Project Logo" title="Illuminating Deposits Project Logo" />
</p>


# output for http trigger initial-v0.1
Http trigger with curl:

`
curl -m 70 -X POST https://illuminating-deposits-vzeropoint1-2fzqdixaqa-uc.a.run.app \
-H "Authorization: bearer $(gcloud auth print-identity-token)" \
-H "Content-Type: application/json" \
-d '{
"name": "Whats  up"
}
`

Output:

`Hello, Whats  up!`

# Version
v0.0.5