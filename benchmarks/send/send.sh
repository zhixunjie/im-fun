for ((i=0;i<100000;i++)); do
curl -XPOST  "http://127.0.0.1:8080/im/send/to/users" \
-d '{
        "tcp_session_ids": [
            {
                "user_id": 1001,
                "user_key": "x4u5mmq6gh2md5dl"
            }
        ],
        "sub_id": 0,
        "message": "你在干什么呢？"
    }'
done