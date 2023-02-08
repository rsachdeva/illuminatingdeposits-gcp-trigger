.PHONY: init
init:
	$(MAKE) -C gcf_upload init && $(MAKE) -C gcf_interestcal init

.PHONY: destroy
destroy:
	$(MAKE) -C gcf_interestcal destroy && $(MAKE) -C gcf_upload destroy

.PHONY: apply
apply:
	$(MAKE) -C gcf_upload apply && $(MAKE) -C gcf_interestcal apply