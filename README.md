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
## Cloud Function Deploy: gcf_upload

### Make steps for gcf_upload deployment in cloud:
Add alias tf=terraform in .zshrc or equivalent
Steps start from root of project folder
1. `cd gcf_upload`
2. `make init`
3. `make apply`
In the end on successful creation you will get something like:
`google_cloudfunctions2_function.illuminating_gcf_upload: Creation complete`
4. After you are done using this function and no longer need for any processing, `make destroy`
In the end on successful destruction you will get something like:
`google_storage_bucket.illuminating_gcf_upload_bucket: Destruction complete`

#### Testing gcf_upload in cloud:
1. make cloud-incorrect-json
2. make cloud-not-successful-http-request-accounttypemissing
3. make cloud-not-successful-http-request-banknamemissing
4. make cloud-successful-http-request

### Make steps for gcf_upload deployment locally:
make gcf-local

#### Testing gcf_upload locally:
1. make local-incorrect-json
2. make local-not-successful-http-request-accounttypemissing
3. make local-not-successful-http-request-banknamemissing
4. make local-successful-http-request


### Seeing logs in console on the web
1. Go to Log tab of Google Cloud Function
2. Click Log Explorer
3. At the bottom Click Edit Time to see logs for the last 15 minutes 
(even this can be edited to say last 10 minutes)
Per the retention period these logs delete after 30 days. There is no charge for the 30-day period.
(https://cloud.google.com/logging#section-7)
4. Click Refresh button to see the logs 

To see logs using gcloud command line:
* gcloud config set project illuminatingdeposits-gcp
* gcloud auth login 
* sudo pip3 install grpcio
* export CLOUDSDK_PYTHON_SITEPACKAGES=1
* gcloud alpha logging tail --format="default(timestamp,text_payload)"
**---------------------------**

# Version - Initial Placeholders
v0.0.9