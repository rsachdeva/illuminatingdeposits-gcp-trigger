.PHONY: init
init:
	cd deploy/terraform && terraform init

.PHONY: destroy
destroy:
	cd deploy/terraform && terraform destroy -auto-approve

.PHONY: apply
apply: tf-prepare
	cd deploy/terraform && terraform apply -auto-approve

.PHONY: tf-prepare
tf-prepare: clean zip move-to-terraform-folder

.PHONY: clean
clean:
	rm -f *.zip
	rm -f deploy/terraform/*.zip

.PHONY: zip
# If you are using Go 1.16 or higher, you can leave the vendor directory as-is. https://cloud.google.com/functions/docs/writing/specifying-dependencies-go
zip:
	#zip illuminating_gcf_upload_src.zip go.mod vendor/* gcf_upload.go
	zip illuminating_gcf_upload_src.zip go.mod gcf_upload.go



.PHONY: move-to-terraform-folder
move-to-terraform-folder:
	mv illuminating_gcf_upload_src.zip deploy/terraform

# curl targets
LOCAL_URL := http://localhost:8080
.PHONY: local-successful-http-request
local-successful-http-request:
	URL=$(LOCAL_URL) sh httprequestfortrigger/env-curl-httpreq-successful.sh

CLOUD_URL := https://illuminating-gcf-upload-2fzqdixaqa-uc.a.run.app
.PHONY: cloud-successful-http-request
cloud-successful-http-request:
	URL=$(CLOUD_URL) sh httprequestfortrigger/env-curl-httpreq-successful.sh


LOCAL_URL := http://localhost:8080
.PHONY: local-not-successful-http-request-banknamemissing
local-not-successful-http-request-banknamemissing:
	URL=$(LOCAL_URL) sh httprequestfortrigger/env-curl-httpreq-not-successful-banknamemissing.sh

CLOUD_URL := https://illuminating-gcf-upload-2fzqdixaqa-uc.a.run.app
.PHONY: cloud-not-successful-http-request-banknamemissing
cloud-not-successful-http-request-banknamemissing:
	URL=$(CLOUD_URL) sh httprequestfortrigger/env-curl-httpreq-not-successful-banknamemissing.sh


LOCAL_URL := http://localhost:8080
.PHONY: local-not-successful-http-request-accounttypemissing
local-not-successful-http-request-accounttypemissing:
	URL=$(LOCAL_URL) sh httprequestfortrigger/env-curl-httpreq-not-successful-accounttypemissing.sh

CLOUD_URL := https://illuminating-gcf-upload-2fzqdixaqa-uc.a.run.app
.PHONY: cloud-not-successful-http-request-accounttypemissing
cloud-not-successful-http-request-accounttypemissing:
	URL=$(CLOUD_URL) sh httprequestfortrigger/env-curl-httpreq-not-successful-accounttypemissing.sh

LOCAL_URL := http://localhost:8080
.PHONY: local-incorrect-json
local-incorrect-json:
	URL=$(LOCAL_URL) sh httprequestfortrigger/env-curl-httpreq-incorrect-json.sh

CLOUD_URL := https://illuminating-gcf-upload-2fzqdixaqa-uc.a.run.app
.PHONY: cloud-incorrect-json
cloud-incorrect-json:
	URL=$(CLOUD_URL) sh httprequestfortrigger/env-curl-httpreq-incorrect-json.sh

# check targets
.PHONY: unzip-check
unzip-check:
	unzip -p deploy/terraform/illuminating_gcf_upload_src.zip gcf_upload.go | less

.PHONY: mod
mod:
	go mod tidy && go mod vendor

.PHONY: build-check
build-check: mod
	go build .

# local gcf target
.PHONY: gcf-local
gcf-local:
	export FUNCTION_TARGET=UploadHTTP && go run cmd_localgcf/main.go
