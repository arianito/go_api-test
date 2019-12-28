#!/bin/sh

name=`mktemp`
name2=`mktemp`
openssl genrsa -out $name 2048
openssl rsa -in $name -outform PEM -pubout -out $name2
echo '\n'
printf "export PRIVATE_KEY=\"`cat $name | base64 -w 0`\""
echo '\n'
printf "export PUBLIC_KEY=\"`cat $name2 | base64 -w 0`\""
echo '\n'