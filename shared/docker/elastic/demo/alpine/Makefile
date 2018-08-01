.PHONY: build size tags test tar push run ssh stop circle

REPO=blacktop/es-data
ORG=blacktop
NAME=es-data
BUILD ?=$(shell cat LATEST)
LATEST ?=$(shell cat LATEST)

all: build size test

build: ## Build docker image
	cd $(BUILD); docker build --pull -t $(ORG)/$(NAME):$(BUILD) .

size: build ## Get built image size
ifeq "$(BUILD)" "$(LATEST)"
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(BUILD)| cut -d' ' -f1)-blue/' README.md
	sed -i.bu '/latest/ s/[0-9.]\{3,5\}MB/$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(BUILD))/' README.md
endif
	sed -i.bu '/$(BUILD)/ s/[0-9.]\{3,5\}MB/$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(BUILD))/' README.md

tags:
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(ORG)/$(NAME)

test: stop ## Test docker image
	@echo "===> Starting elasticsearch"
	@docker run --init -d --name elasticsearch -p 80:80 -p 9200:9200 blacktop/elastic-stack:geoip; sleep 15
	# @docker run --init -d --name elasticsearch -p 9200:9200 blacktop/elasticsearch:geoip; sleep 15
	# @echo "===> Starting kibana"
	# @docker run --init -d --name kibana --link elasticsearch -p 5601:5601 blacktop/kibana:6.0.0 sleep 5
	@echo "===> Adding es-data to DB"
	@docker run -d --name $(NAME) --link elasticsearch $(ORG)/$(NAME):$(BUILD); sleep 5
	@open -a Safari http://localhost
	@http localhost:9200/_cat/indices
	@http localhost:9200/filebeat-*/_count
	@docker logs -f es-data

tar: ## Export tar of docker image
	docker save $(ORG)/$(NAME):$(BUILD) -o $(NAME).tar

push: build ## Push docker image to docker registry
	@echo "===> Pushing $(ORG)/$(NAME):$(BUILD) to docker hub..."
	@docker push $(ORG)/$(NAME):$(BUILD)

run: stop ## Run docker container
	@docker run --init -d --name $(NAME) $(ORG)/$(NAME):$(BUILD)

ssh: ## SSH into docker image
	@docker run --init -it --rm --entrypoint=sh $(ORG)/$(NAME):$(BUILD)

stop: ## Kill running malice-engine docker containers
	@docker rm -f elasticsearch || true
	@docker rm -f kibana || true
	@docker rm -f $(NAME) || true

circle: ci-size ## Get docker image size from CircleCI
	@sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell cat .circleci/SIZE)-blue/' README.md
	@echo "===> Image size is: $(shell cat .circleci/SIZE)"

ci-build:
	@echo "===> Getting CircleCI build number"
	@http https://circleci.com/api/v1.1/project/github/${REPO} | jq '.[0].build_num' > .circleci/build_num

ci-size: ci-build
	@echo "===> Getting image build size from CircleCI"
	@http "$(shell http https://circleci.com/api/v1.1/project/github/${REPO}/$(shell cat .circleci/build_num)/artifacts circle-token==${CIRCLE_TOKEN} | jq '.[].url')" > .circleci/SIZE

clean: ## Clean docker image and stop all running containers
	docker-clean stop
	docker rmi $(ORG)/$(NAME):$(BUILD) || true

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

# curl -XPUT "http://localhost:9200/.kibana/index-pattern/filebeat-*" -H 'Content-Type: application/json' -d '{"title" : "filebeat-*",  "timeFieldName": "@timestamp"}'
# curl -XPUT http://localhost:9200/.kibana/config/4.5.0 -H 'Content-Type: application/json' -d '{"defaultIndex" : "filebeat-*""}'
