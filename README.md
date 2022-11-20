# TFW
Something not important

```
Done!
Total m 4.547797m
Total success: 114
Total fail: 0

You can watch logs.
```

## Requirements

It is using `modernc.org/sqlite` with CGO.

* musl-gcc to static compile

```bash
go mod tidy -v
```

## Build

```bash
cd cmd
make
```
for build static