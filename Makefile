
#
# runs tests for all subpackages.
#
.PHONY: test
test:
	go test ./...

#
# Builds gl, a simple tool that can be used to run lisp files.
#
.PHONY: bin/gl
bin/gl:
	go build -o bin/gl ./cmds/gl

#
# Executes all files in the examples directory.
#
.PHONY: run-gl-examples
run-gl-examples: bin/gl
	find examples -maxdepth 1 -type f -name '*.l' | xargs -I{} bin/gl {}
