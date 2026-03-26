# gobox
This is a small go project that makes use of the [gin](https://gin-gonic.com/) package to make a simple file hosting site. Heavily inspired by [catbox.moe](https://catbox.moe).

# Disclaimer
This project should not be used in a production, or professional environment, as it lacks A LOT of security measures, and is extremely simple.

# Features
- Simple file uploading
- Unique file names by using the file hash instead of the filename
- File serving using regexp to filter invalid filenames or insecurities.
- Low resource usage

# Usage
You can download the project from the Versions tab, simply run it with `go run main.go` , or build it yourself.
# Flags
```-p [PORT] The port the server will run in. 8080 by default.```
```-ip [IP] The IP address the server will run in. localhost by default.```