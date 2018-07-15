#!/usr/bin/env bash
echo "Get access token from Qiniu"
http GET http://localhost:8088/v1/oss/secrets?cloud=qiniu&bucket==images&options=''