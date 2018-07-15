
#!/usr/bin/env bash
http POST http://localhost:8080/callback/oss-put-object \
bucket='moremom-obj' \
object='joehart.jpg' \
etag='78F2F5E8F6B7FE9F793F27F0FE291F61' \
size='17689' \
imageInfo.format='jpg' \
imageInfo.width='457' \
imageInfo.height='343' \
appName='moremom' \
appUserID='123456' \
appBusiness='avatar' \
appToken='random_123456'
