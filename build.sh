#!/bin/bash
GOOS=linux GOARCH=arm GOARM=7 go build && rsync -av * pi:~/go-led-stream/
