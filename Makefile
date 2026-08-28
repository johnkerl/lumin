PREFIX=$(HOME)
INSTALLDIR=$(PREFIX)/bin

# ================================================================
# General-use targets

# This must remain the first target in this file, which is what 'make' with no
# arguments will run.
build:
	go build
	@echo "Build complete. The lumin executable is ./lumin (or .\lumin.exe on Windows)."

# DESTDIR is for package installs; nominally blank when this is run interactively.
# See also https://www.gnu.org/prep/standards/html_node/DESTDIR.html
install: build
	mkdir -p $(DESTDIR)/$(INSTALLDIR)
	cp lumin $(DESTDIR)/$(INSTALLDIR)/

# ----------------------------------------------------------------
# Formatting
# go fmt ./... finds experimental C files which we want to ignore.
fmt format:
	-go fmt ./cmd/...
	-go fmt ./pkg/...
	-go fmt ./regression_test.go

# ----------------------------------------------------------------
# Static analysis

# Needs first: go install honnef.co/go/tools/cmd/staticcheck@latest
# See also: https://staticcheck.io
staticcheck:
	staticcheck ./pkg/... ./cmd/mlr/...

# Needs first: https://golangci-lint.run/welcome/install/#local-installation
# Config lives in .golangci.yml; same invocation as .github/workflows/golangci-lint.yml
lint:
	golangci-lint run ./cmd/mlr ./pkg/...

lumin:
	go build

# ================================================================
# Go does its own dependency management, outside of make.
.PHONY: build fmt format staticcheck lint lumin
