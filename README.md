# TFW

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/2868bf70612b463cbf4f511abe4b2020)](https://app.codacy.com/gh/raifpy/tfw?utm_source=github.com&utm_medium=referral&utm_content=raifpy/tfw&utm_campaign=Badge_Grade_Settings)

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