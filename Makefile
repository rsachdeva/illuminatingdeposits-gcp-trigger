#https://web.mit.edu/gnu/doc/html/make_5.html#SEC48 Recursive Use of make
.PHONY: init
init:
	$(MAKE) -C gcf_upload init && $(MAKE) -C gcf_interestcal init && $(MAKE) -C gcf_analytics init  && $(MAKE) -C gcf_notify init

.PHONY: destroy
destroy:
	$(MAKE) -C gcf_notify destroy && $(MAKE) -C gcf_analytics destroy && $(MAKE) -C gcf_interestcal destroy && $(MAKE) -C gcf_upload destroy

.PHONY: apply
apply:
	$(MAKE) -C gcf_upload apply && $(MAKE) -C gcf_interestcal apply && $(MAKE) -C gcf_analytics apply && $(MAKE) -C gcf_notify apply

fmt:
	$(MAKE) -C gcf_upload fmt && $(MAKE) -C gcf_interestcal fmt && $(MAKE) -C gcf_analytics fmt && $(MAKE) -C gcf_notify fmt

fmt-diff:
	$(MAKE) -C gcf_upload fmt-diff && $(MAKE) -C gcf_interestcal fmt-diff && $(MAKE) -C gcf_analytics fmt-diff && $(MAKE) -C gcf_notify fmt-diff

lint:
	$(MAKE) -C gcf_upload lint && $(MAKE) -C gcf_interestcal lint && $(MAKE) -C gcf_analytics lint && $(MAKE) -C gcf_notify lint

lint-fix:
	$(MAKE) -C gcf_upload lint-fix && $(MAKE) -C gcf_interestcal lint-fix && $(MAKE) -C gcf_analytics lint-fix && $(MAKE) -C gcf_notify lint-fix

## go fix is the Go modernizer tool
fix-diff:
	$(MAKE) -C gcf_upload fix-diff && $(MAKE) -C gcf_interestcal fix-diff && $(MAKE) -C gcf_analytics fix-diff && $(MAKE) -C gcf_notify fix-diff

fix:
	$(MAKE) -C gcf_upload fix && $(MAKE) -C gcf_interestcal fix && $(MAKE) -C gcf_analytics fix && $(MAKE) -C gcf_notify fix
