#for ((i=0;i<100000;i++)); do
for ((i=0;i<10;i++)); do
curl -XPOST  "http://127.0.0.1:8080/im/send/to/users/by/ids" \
-d '{
        "uni_ids": [
           "1"
        ],
        "sub_id": 0,
        "message": "你在干什么呢？"
    }'
echo
done