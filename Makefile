#!/bin/sh
dev:
	CGO_ENABLED=0 go build  -o /tmp/instagram .
	. ./.env && /tmp/instagram