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
zip: mod
	#zip illuminating_gcf_upload_src.zip go.mod vendor/* gcf_upload.go
	zip illuminating_gcf_upload_src.zip go.mod gcf_upload.go

.PHONY: mod
mod:
	go mod tidy && go mod vendor

.PHONY: move-to-terraform-folder
move-to-terraform-folder:
	mv illuminating_gcf_upload_src.zip deploy/terraform


.PHONY: successful-http-request
successful-http-request:
	sh httprequestfortrigger/curl-httpreq-successful.sh

.PHONY: not-successful-http-request
not-successful-http-request:
	sh httprequestfortrigger/curl-httpreq-not-successful.sh

.PHONY: unzip-check
unzip-check:
	unzip -p deploy/terraform/illuminating_gcf_upload_src.zip gcf_upload.go | less

.PHONY: build-check
build-check:
	go build .
