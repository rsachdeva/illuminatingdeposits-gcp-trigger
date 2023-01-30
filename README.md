# Illuminating Deposits - GCP Trigger Applied to Cloud Functions

Illuminating Deposits Google Cloud Platform GCP based trigger of Resources applied to Cloud Functions

(Development is WIP)

Illuminating Deposits Project Logo: 

![Illuminating Deposits Project Logo](logo.png "Illuminating Deposits Project Logo")

## Overall google cloud architecture system design:
![google cloud architecture system design](GoogleCloudArchitectureSystemDesign.png "google cloud architecture system design")

Created from link:

[Google Cloud Architecture](https://googlecloudcheatsheet.withgoogle.com/architecture)

**---------------------------**
## Cloud Function: gcf_upload
### Make steps for gcf_upload:
Add alias tf=terraform in .zshrc or equivalent
Steps start from root of project folder
1. `cd gcf_upload`
2. `make init`
3. `make apply`
In the end on successful completion you will get something like:
`google_cloudfunctions2_function.illuminating_gcf_upload: Creation complete`
4. After you are done using this function and no longer need for any processing, `make destroy`

### output for gcf_upload:
Http trigger with curl for interestcal:

`
curl -m 70 -X POST https://illuminating-deposits-vzeropoint1-2fzqdixaqa-uc.a.run.app \
-H "Authorization: bearer $(gcloud auth print-identity-token)" \
-H "Content-Type: application/json" \
-d '{
"name": "Whats  up"
}'
`

Output:

`Hello, Whats  up!`

**---------------------------**

## Cloud function
Http trigger with curl for notifyslack:
curl -m 70 -X POST https://notifyslack-vzeropoint1-2fzqdixaqa-uc.a.run.app \
-H "Authorization: bearer $(gcloud auth print-identity-token)" \
-H "Content-Type: application/json" \
-d '{
"name": "Hello World"
}'

Output:

Slack channel notified 
with message:
Triggering Illuminating Calculation Wrap Up

# Version - Initial Placeholders
v0.0.8