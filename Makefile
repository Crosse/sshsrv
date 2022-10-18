NAME	= sshsrv
PACKAGE = github.com/Crosse/$(NAME)

default: release

define build
	@env GOOS=$(1) GOARCH=$(2) make release/$(NAME)-$(1)-$(2)$(3)
endef

release/$(NAME)-%:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o "$@" $(PACKAGE)

release/$(NAME)-darwin-universal: release/$(NAME)-darwin-amd64 release/$(NAME)-darwin-arm64
	$(RM) "$@"
	lipo -create -o "$@" $^

.PHONY: release
release:
	mkdir -p release
	$(call build,linux,arm)
	$(call build,linux,amd64)
	$(call build,linux,arm64)

	$(call build,darwin,amd64)
	$(call build,darwin,arm64)
	@make release/$(NAME)-darwin-universal

	$(call build,openbsd,arm)
	$(call build,openbsd,amd64)
	$(call build,openbsd,arm64)

	$(call build,freebsd,arm)
	$(call build,freebsd,amd64)
	$(call build,freebsd,arm64)

	$(call build,windows,amd64,.exe)
	$(call build,windows,arm64,.exe)

.PHONY: zip
zip: release
	find release -type f ! -name '*.zip' -execdir zip -9 "{}.zip" "{}" \;

.PHONY: clean
clean:
	$(RM) -r release
