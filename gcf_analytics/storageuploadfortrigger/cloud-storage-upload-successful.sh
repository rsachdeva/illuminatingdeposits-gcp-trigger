# gsutil is a Python application that lets you access Cloud Storage from the command line
# gsutil cp <LOCAL_FILE_PATH> gs://<BUCKET_NAME>/<OBJECT_NAME>
# interestresponse-for-cli-upload.json is used for this command for direct upload of interest response to output bucklet to trigger analytics GCF
gsutil cp interestresponse-for-cli-upload.json gs://illuminating_upload_json_bucket_output/interestresponse.json