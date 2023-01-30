.PHONY: tf-prepare-apply
tf-prepare-apply: tf-prepare apply

.PHONY: tf-prepare
tf-prepare: clean zip move-to-terraform-folder

.PHONY: init
init:
	cd deploy/terraform && terraform init

.PHONY: destroy
destroy:
	cd deploy/terraform && terraform destroy -auto-approve

.PHONY: apply
apply:
	cd deploy/terraform && terraform apply -auto-approve

.PHONY: clean
clean:
	rm -f *.zip
	rm -f deploy/terraform/*.zip

.PHONY: zip
zip: mod
	zip illuminating_gcf_upload_src.zip go.mod vendor/* gcf_upload.go

.PHONY: mod
mod:
	go mod tidy && go mod vendor

.PHONY: move-to-terraform-folder
move-to-terraform-folder:
	mv illuminating_gcf_upload_src.zip deploy/terraform

#only for checking

.PHONY: unzip-check
unzip-check:
	unzip -p deploy/terraform/illuminating_gcf_upload_src.zip gcf_upload.go | less

.PHONY: build-check
build-check:
	go build .
