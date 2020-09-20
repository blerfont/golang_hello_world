#!/bin/bash
PROJECT_FOLDER=$(pwd)
APP_NAME=$(grep "module " "./go.mod" | cut -d " " -f 2 | tr -d '[:space:]')
VERSION=$(grep "version" "./version.go" | cut -d "=" -f 2 | tr -d '"' | tr -d '[:space:]')
DATE=$(date +%s)
TOOLCHAIN=$1

if [ -z "$TOOLCHAIN" ];then
	echo ""
	echo "Please provide a valid toolchain."
	echo "./build.sh armv6/armv7/aarch64/x86_64"
	echo ""
	exit 1
fi

[ -f "$PROJECT_FOLDER/$APP_NAME" ] && rm -f "$PROJECT_FOLDER/$APP_NAME"
[ -f "$PROJECT_FOLDER/$APP_NAME-$VERSION-$TOOLCHAIN.zip" ] && rm -f "$PROJECT_FOLDER/$APP_NAME-$VERSION-$TOOLCHAIN.zip"

echo "Cleaning not needed dependencies"
GOOS=linux go mod tidy
echo "Updating dependencies"
GOOS=linux go get -u

echo "Building $APP_NAME for $TOOLCHAIN toolchain"
case $TOOLCHAIN in
	armv6)
		GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "-w -s"
		;;
	armv7)
		GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-w -s"
		;;
	aarch64)
		GOOS=linux GOARCH=arm64 go build -ldflags "-w -s"
		;;
	x86_64)
		GOOS=linux GOARCH=amd64 go build -ldflags "-w -s"
		;;
	*)
		echo "Unsupported toolchain. armv6/armv7/aarch64/x86_64"
		exit 1
		;;
esac

if [ -f "$PROJECT_FOLDER/$APP_NAME" ];then
  echo "Generating app zip"
	mkdir -p "/tmp/$DATE/bin"
	[ -d "$PROJECT_FOLDER"/files ] && cp -rf "$PROJECT_FOLDER/files" "/tmp/$DATE/"
	cp -f "$PROJECT_FOLDER/$APP_NAME" "/tmp/$DATE/bin/"
	echo "$APP_NAME" > "/tmp/$DATE/init"
	cd /tmp/"$DATE"/ || exit 1
	zip -rq "$PROJECT_FOLDER/$APP_NAME-$VERSION-$TOOLCHAIN.zip" .
	cd "$PROJECT_FOLDER" || exit 1
	rm -rf "/tmp/$DATE/"
	rm -f "$PROJECT_FOLDER/$APP_NAME"
	if [ -f "$PROJECT_FOLDER/$APP_NAME-$VERSION-$TOOLCHAIN.zip" ];then
	  echo "App zip generated correctly: $PROJECT_FOLDER/$APP_NAME-$VERSION-$TOOLCHAIN.zip"
	else
	  echo "Zip build failed"
	fi
else
	echo "Something failed. Check logs"
fi