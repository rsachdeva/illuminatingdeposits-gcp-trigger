# Illuminating Deposits - Google Cloud Platform (GCP) trigger of Resources applied to Golang (Go) Google Cloud Functions (GCF 2nd Generation)

(Development is WIP)

Illuminating Deposits Project Logo: 

![Illuminating Deposits Project Logo](logo.png "Illuminating Deposits Project Logo")

## Overall google cloud architecture system design:
![google cloud architecture system design](GoogleCloudArchitectureSystemDesign.png "google cloud architecture system design")

Created from link:

[Google Cloud Architecture](https://googlecloudcheatsheet.withgoogle.com/architecture)

Please see GREEN in System Diagram. Those are Golang (Go) based Google Clod Functions (GCF) that have initial versions completed.
System design can change as the project progresses.

**---------------------------**
## Cloud Function Deploy: gcf_upload
gcf_upload google cloud function is triggered by HTTP request and submits HTTP POST request body to
google cloud storage bucket for further processing.

Terraform is used to deploy gcf_upload to google cloud.
The service account used for deployment needs keys to be generated and stored in the deploy/terraform folder
with the file being renamed to keys.json.
The gcf_upload is written in Go and uses google cloud storage and google cloud functions.

### Make steps for gcf_upload deployment in cloud:
Add alias tf=terraform in .zshrc or equivalent
Steps start from root of project folder
1. `cd gcf_upload`
2. `make init`
3. `make apply`
In the end on successful creation you will get something like:
`google_cloudfunctions2_function.illuminating_gcf_upload: Creation complete`
4. Dependent Google cloud functions should be destroyed first. So if you are also working on
gcf_interestcal also for example (see below) then see destruction in gcf_interestcal.
If you are only playing around with this Google cloud function,
after you are done using this function and no longer need for any processing, `make destroy`
In the end on successful destruction you will get something like:
`google_storage_bucket.illuminating_gcf_upload_bucket: Destruction complete`
 Also, See Related [Multiple Cloud Functions Deploy: gcf_interestcal && gcf_upload for more details](#multiple-cloud-functions-deploy-gcf_interestcal--gcf_upload)

#### Testing gcf_upload in cloud:
1. `make cloud-incorrect-json`
2. `make cloud-not-successful-http-request-accounttypemissing`
3. `make cloud-not-successful-http-request-banknamemissing`
4. `make cloud-successful-http-request`

### Make steps for gcf_upload deployment locally:
Steps start from root of project folder
1. `cd gcf_upload`
2. `make gcf-local`

#### Testing gcf_upload locally:
1. `cd gcf_upload`
2. `make local-incorrect-json`
3. `make local-not-successful-http-request-accounttypemissing`
4. `make local-not-successful-http-request-banknamemissing`
5. `make local-successful-http-request`


### Seeing logs in console on the web
1. Go to Log tab of Google Cloud Function
2. Click Log Explorer
3. At the bottom Click Edit Time to see logs for the last 15 minutes 
(even this can be edited to say last 10 minutes)
Per the retention period these logs delete after 30 days. There is no charge for the 30-day period.
([https://cloud.google.com/logging#section-7](https://storage.cloud.google.com/illuminating_upload_json_bucket/inputrequest.json))
4. Click Refresh button to see the logs 

To see logs using gcloud command line:
* `gcloud config set project illuminatingdeposits-gcp`
* `gcloud auth login `
* `sudo pip3 install grpcio`
* `export CLOUDSDK_PYTHON_SITEPACKAGES=1`
* `gcloud alpha logging tail --format="default(timestamp,text_payload)"`
**---------------------------**


**---------------------------**
## Cloud Function Deploy: gcf_interestcal
gcf_interstcal google cloud function is triggered by upload to illuminating_upload_json_bucket_input by gcf_upload cloud
function  and starts calculating Delta for each bank account including overall for all deposits. It then also
uploads the calculated  Deltas with deposit details to illuminating_upload_json_bucket_output google cloud storage bucket
for further processing. 30 days Interest for a deposit is called Delta here. So if a deposit is for 10 years the Delta
is still calculated for 30 days Interest using the Interest for 10 years to get an approx. value of 30 days Interest.

Terraform is used to deploy gcf_interestcsal to google cloud.
The service account used for deployment needs keys to be generated and stored in the deploy/terraform folder
with the file being renamed to keys.json.
The gcf_interest is written in Go and uses google cloud storage and google cloud functions.

### Make steps for gcf_interestcal deployment in cloud:
Add alias tf=terraform in .zshrc or equivalent
Steps start from root of project folder
1. `cd gcf_interestcal`
2. `make init`
3. `make apply`
   In the end on successful creation you will get something like:
   `google_cloudfunctions2_function.illuminating_gcf_interestcal: Creation complete`
4. Then delete the function that this function depends on so do the following:
    4a: `cd ../gcf_upload`
    4b: `make destroy` In the end on successful destruction you will get something like:
   `google_storage_bucket.illuminating_gcf_upload_bucket: Destruction complete`
5. After you are done using this function and no longer need for any processing, `make destroy`
   In the end on successful destruction you will get something like:
   `google_storage_bucket.illuminating_gcf_interestcal_bucket: Destruction complete`
Also, See Related [Multiple Cloud Functions Deploy: gcf_interestcal && gcf_upload for more details](#multiple-cloud-functions-deploy-gcf_interestcal--gcf_upload)

#### Testing gcf_interest in cloud:
`make cloud-successful-storage-upload`
This copes the file locally to cloud storage

### Make steps for gcf_interestcal deployment locally:
Steps start from root of project folder
1. `cd gcf_interestcal`
2. `make gcf-local`

#### Testing gcf_interestcal locally:
`make local-successful-send-event-to-cloud-event-function`
This only triggers the function locally and does not copy the file to cloud storage.
You can send san event to CloudEvent function locally with already existing file in cloud storage
**---------------------------**

## Multiple Cloud Functions 

### Multiple Cloud Functions Deploy: gcf_interestcal && gcf_upload
Add alias tf=terraform in .zshrc or equivalent
Steps start from root of project folder at `illuminatingdeposits-gcp-trigger`

1. `make init`
2. `make apply`
3. `make destroy`

#### Integration Testing gcf_upload and gcf_interestcal in cloud:
SInce everything gets triggered with gcf_upload, we will test gcf_upload and gcf_interestcal together using gcf_upload http request submission:
1. `cd gcf_upload`
2. `make local-successful-http-request`

# Version
v0.2.0
