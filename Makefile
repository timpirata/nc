# Makefile for nc - The random numerus clausus generator.

BUILD_ROOT := ./build
NC := $(BUILD_ROOT)/nc 
NC_VERSION := $(shell git describe --tags | cut -dv -f2)
LDFLAGS := -X main.AppVersion=$(NC_VERSION)
DOCKER_IMAGE := schnoddelbotz/nc
DOCKER_TAGS := -t $(DOCKER_IMAGE):latest -t $(DOCKER_IMAGE):$(NC_VERSION)
GOCURL := $(BUILD_ROOT)/gocurl

# If I wasn't into semantic versioning, I'd have this one called
# Basic Application Example, 4. You.

# HTML output dependency MathJax, to render TeX formulae.
# to update, set version number and run `make update-mathjax`.
# Then commit to git, to avoid JS build dependencies for us.
MATHJAX_VERSION := 3.2.2
MATHJAX_URL := https://github.com/mathjax/MathJax/archive/refs/tags/$(MATHJAX_VERSION).zip
MATHJAX_ZIP := $(BUILD_ROOT)/$(MATHJAX_VERSION).zip
MATHJAX_DIR := $(BUILD_ROOT)/MathJax
MATHJAX_CC  := tex-chtml-full-speech.js
MATHJAX_TGT := ./output/templates/js/$(MATHJAX_CC)

nc: $(NC) ESC

$(NC): ESC cmd/nc/*.go output/*.go quiz/*.go go.mod
	go generate ./...
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(NC) ./cmd/nc

ESC:
	test -n "$(shell which esc)" || go install github.com/mjibson/esc

# TODO / TESTS 

tests/output/5-additions.html: nc 
	$(NC) -a 5 -f html -o $< 

tests/output/3-multiplications-2-divisions.html: nc 
	$(NC) -m 3 -d 2 -f html -o $<

tests/output/5-each.pdf: nc 
	# should produce PDF with 5 quiz/exam questions per type
	$(NC) -A 5 -f pdf -o $<

clean:
	rm -f $(NC) $(GOCURL) $(MATHJAX_ZIP) build/_docker_image
	rm -rf $(MATHJAX_DIR) 

serve:
	# run NC as local webserver, serving quizzes on port 7898
	$(NC) -s -p 7898

deploy:
	# should deploy as cloudfunction 
	gcloud deploy

#
# Docker 

build/_docker_image:
	docker build $(DOCKER_TAGS) --progress plain .
	touch build/_docker_image

docker-image-uncached:
	docker build $(DOCKER_TAGS) --progress plain --no-cache .

docker-image: build/_docker_image
	
docker-ratzeputz:
	docker system prune -af

docker-serve: docker-image
	docker run -p2001:2001 $(DOCKER_IMAGE):latest ./nc -serve

ratzeputz: clean docker-ratzeputz

#
# Dependencies 

# instead of relying on curl or wget, we quickly build or own Go downloader,
# which accidently also includes unzip. Nice, no? Called bootstrapping ;-*
gocurl: cmd/gocurl/main.go 
	go build -o $(GOCURL) ./cmd/gocurl 

# https://github.com/mathjax/MathJax 
# https://docs.mathjax.org/en/latest/web/components/combined.html
update-mathjax: gocurl
	mkdir -p $(BUILD_ROOT)
	$(GOCURL) -remote $(MATHJAX_URL) -local $(MATHJAX_ZIP) -unzipTo $(MATHJAX_DIR)
	cp $(MATHJAX_DIR)/MathJax-$(MATHJAX_VERSION)/es5/$(MATHJAX_CC) $(MATHJAX_TGT) 
	mkdir -p output/templates/js/output/chtml
	cp -r $(MATHJAX_DIR)/MathJax-$(MATHJAX_VERSION)/es5/output/chtml/fonts output/templates/js/output/chtml
