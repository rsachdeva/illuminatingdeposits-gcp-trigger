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
	#zip illuminating_gcf_interestcal_src.zip go.mod vendor/* gcf_interestcal.go
	zip illuminating_gcf_interestcal_src.zip go.mod *.go



.PHONY: move-to-terraform-folder
move-to-terraform-folder:
	mv illuminating_gcf_interestcal_src.zip deploy/terraform

.PHONY: cloud-successful-storage-upload
cloud-successful-storage-upload:
	cd storageuploadfortrigger && sh ./cloud-storage-upload-successful.sh

.PHONY: local-successful-storage-upload
# https://cloud.google.com/functions/docs/running/calling#cloudevent-function-curl-tabs-storage
local-successful-send-event-to-cloud-event-function:
	cd storageuploadfortrigger && sh ./local-send-event-to-cloud-event-function.sh

# check targets
.PHONY: unzip-check
unzip-check:
	unzip -p deploy/terraform/illuminating_gcf_interestcal_src.zip gcf_interestcal.go | less

.PHONY: mod
mod:
	go mod tidy && go mod vendor

.PHONY: build-check
build-check: mod
	go build .

# local gcf target
.PHONY: gcf-local
gcf-local:
	export FUNCTION_TARGET=InterestCalStorage && go run cmd_localgcf/main.go
