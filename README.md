# forwarding
A very liteweight tool to forward data over tcp,  written in Go.



#Client A want to broadcast data to client B, the process is
```
client A [TCP stream encrypted by AES256]----> forwarding-server ----> [TCP stream encrypted by AES256] client B
```


# 1 A send a json

{
    "sync_id":id0
    "req":"handshake_0"
}

# 2 forwarding-server send back a json to A
{
    "sync_id":id1,
    "channel_id":"ch0"
    "rsp":"handshake_1"
}

# 3 A send a json
{
    "sync_id":id2
    "req":"handshake_2"
    "recv_channel_ids":["ch1","ch2","ch3"],
    "max_timeout_sec":3600, // optinal
}

# 4 forwarding-server send back a json to A
{
    "sync_id":id3
    "req":"handshake_3"
    "handshake_result":"OK"
}

# 5 A send data to forwarding-server who fowarding the data to B (all recv_channel_ids receivers) directly.

# 6 A close the connection, then forwarding-server clear this session infomations.



