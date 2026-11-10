#!/bin/env sh

GOOS=linux GOARCH=arm64 go build -o rpi/boot/inshrimp .

cd frontend
npm run build
cd ..
cp frontend/dist rpi/boot

