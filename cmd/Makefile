static:
	env CC=musl-gcc CGO_ENABLED=1 go build -ldflags="-w -s -linkmode external -extldflags '-static'" --trimpath -v -x -o static
