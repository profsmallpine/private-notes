# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
live/templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=false

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	air \
	--build.cmd "go build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# start all watch processes in parallel.
dev: 
	make -j5 live/templ live/server