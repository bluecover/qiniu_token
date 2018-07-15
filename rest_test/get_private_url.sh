#!/usr/bin/env bash
echo "Get private url from Qiniu"
http GET http://localhost:8088/object/url \
cloud==qiniu \
domain==http://p55fwwcvn.bkt.clouddn.com \
key=='test/2018/03/Fr8RAK-jHIlndVrZoZK9v3T43m5r.jpg'
