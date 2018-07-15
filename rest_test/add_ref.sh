#!/usr/bin/env bash
http POST http://localhost:8088/v1/stash/addref \
userID:=123 \
objectID:=21 \
tag=picviewer \
object:='{"cloud":"qiniu","bucket":"images","key":"est/2018/03/Fto5o-5ea0sNMlW_75VgGJCv2AcJ.mp4","etag":"Fqe2rcRvLIaPCWy3FCKw2Qw9kL7o","mimeType":"video/mp4","size":2190427,,"tag":"picviewer"}'
