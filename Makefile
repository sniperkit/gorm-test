
# add auto-detect
LOCAL_OS 			:= darwin
LOCAL_ARCH 			:= amd64

# golang deps
.PHONY: deps
deps:
	@glide install --strip-vendor

.PHONY: deps-rebuild
deps-rebuild:
	@if [ -f "glide.yaml" ]; then rm -f ./glide.*; fi
	@yes no | glide create
	@glide install --strip-vendor

# goss variables
GOSS_VERSION 		:= 0.3.6
GOSS_REPO_FULLNAME  := aelsabbahy/goss
GOSS_LOCAL_DIR 		:= ./shared/tools/goss
GOSS_BIN_DIR  		:= /usr/local/bin
GOSS_BIN_BASENAME   := goss
GOSS_BIN_PATH 		:= $(GOSS_BIN_DIR)/$(GOSS_BIN_BASENAME)

.PHONY: goss
# https://github.com/aelsabbahy/goss/tree/master/extras/dgoss#mac-osx
goss:
	@curl -L https://github.com/$(GOSS_REPO_FULLNAME)/releases/download/v$(GOSS_VERSION)/$(GOSS_BIN_BASENAME)-$(LOCAL_OS)-$(LOCAL_ARCH) -o $(GOSS_BIN_PATH)
	@chmod +rx $(GOSS_BIN_PATH)

# crane variables
CRANE_VERSION 		:= 3.4.2
CRANE_REPO_FULLNAME := michaelsauter/crane
CRANE_LOCAL_DIR 	:= ./shared/tools/crane
CRANE_BIN_DIR  		:= /usr/local/bin
CRANE_BIN_BASENAME  := crane
CRANE_BIN_PATH 		:= $(CRANE_BIN_DIR)/$(CRANE_BIN_BASENAME)

.PHONY: crane
crane:
	@echo "CRANE_PRO_KEY: $(CRANE_PRO_KEY)"
	@if [ "$(CRANE_PRO_KEY)" != "" ]; then \
	    mkdir -p $(CRANE_LOCAL_DIR) && \
		cd $(CRANE_LOCAL_DIR)&& \
		curl -sL https://raw.githubusercontent.com/$(CRANE_REPO_FULLNAME)/v$(CRANE_VERSION)/download.sh && \
		mv $(CRANE_BIN_BASENAME) $(CRANE_BIN_PATH) && \
	else \
		curl -Lo /usr/local/bin/crane "https://www.craneup.tech/downloads/mac/v$(CRANE_VERSION)?key=$(CRANE_PRO_KEY)" ;
	fi \
	chmod +x $(CRANE_BIN_PATH))

.PHONY: crane-build
crane-build:
	@go get -u github.com/michaelsauter/crane
	@go install github.com/michaelsauter/crane

# lift groups

## lift dev container (interactive)
.PHONY: lift-dev
lift-dev:
	@crane lift dev

## lift rmdb related containers
.PHONY: lift-rmdb
lift-rmdb:
	@crane lift app