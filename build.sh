#!/bin/bash
GOOS=linux GOARCH=arm GOARM=7 go build && rsync -av * pi:~/go-led-stream/ && ssh pi 'sudo systemctl restart led-stream'
