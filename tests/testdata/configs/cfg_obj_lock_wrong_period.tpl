[cloud_storage]
auth_mode = apikey
auth_keypath = {{APIKEYFILE}}
bucket = bucket
region = us-east
endpoint_url = https://s3.direct.us-east.cloud-object-storage.appdomain.cloud

[objects]
object_tags = key1=val1,key2=val2
object_lock_retention_mode = cmp
object_lock_retention_period = 1,2

[backint]
max_concurrency = 8
multipart_chunksize = 512MB