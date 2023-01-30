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
zip:
	zip illuminating_gcf_upload_src.zip go.mod gcf_upload.go vendor

.PHONY: move-to-terraform-folder
move-to-terraform-folder:
	mv illuminating-gosource.zip deploy/terraform

#only for checking

.PHONY: unzip-check
unzip-check:
	unzip -p deploy/terraform/illuminating-gosource.zip hello.go | less

.PHONY: build-check
build-check:
	go build .
