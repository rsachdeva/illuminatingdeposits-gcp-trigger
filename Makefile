#https://web.mit.edu/gnu/doc/html/make_5.html#SEC48 Recursive Use of make
.PHONY: init
init:
	$(MAKE) -C gcf_upload init && $(MAKE) -C gcf_interestcal init && $(MAKE) -C gcf_analytics init

.PHONY: destroy
destroy:
	$(MAKE) -C gcf_analytics destroy && $(MAKE) -C gcf_interestcal destroy && $(MAKE) -C gcf_upload destroy

.PHONY: apply
apply:
	$(MAKE) -C gcf_upload apply && $(MAKE) -C gcf_interestcal apply && $(MAKE) -C gcf_analytics apply