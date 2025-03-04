# Define build target path
BUILD_DIR := bin
BINARY := danta-auto-tool

# Build the project
.PHONY: build
build:
	bash scripts/build.sh "$(BUILD_DIR)" "$(BINARY)"

# Run the project
.PHONY: run
run:
	bash scripts/run.sh "$(BUILD_DIR)/$(BINARY)"

# Clean the program output
.PHONY: clean-output
clean-output:
	bash scripts/clean_output.sh

# Clean the build
.PHONY: clean-build
clean-build:
	bash scripts/clean_build.sh "$(BUILD_DIR)/$(BINARY)"

# Clean both output and build
.PHONY: clean
clean: clean-output clean-build
