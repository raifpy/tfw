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

It is using `github.com/mattn/go-sqlite3` with CGO.

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


[![Codacy Badge](https://app.codacy.com/project/badge/Grade/485fab2c584a4474974faac8af0a589d)](https://www.codacy.com/gh/raifpy/tfw/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=raifpy/tfw&amp;utm_campaign=Badge_Grade)