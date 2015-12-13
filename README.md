# forwarding
A very liteweight tool to forward data over tcp,  written in Go.

这有点像一个turn服务，可以转发任意两端(多端)之间的tcp层报文，所以这必然是一个神器。具体能干嘛? 你懂的，嘿嘿~



#Client A want to broadcast data to client B, the process is
```
client A [TCP stream encrypted by AES256]----> forwarding-server ----> [TCP stream encrypted by AES256] client B
```


# Step 1 
A send a json

{
    "sync_id":id0
    "req":"handshake_0"
    "channel_id":"ch0" // optional
}


# Step 2 
forwarding-server send back a json to A
{
    "sync_id":id1,
    "channel_id":"ch1"
    "rsp":"handshake_1"
}


# Step 3 
A send a json
{
    "sync_id":id2
    "req":"handshake_2"
    "recv_channel_ids":["ch2","ch3","ch4"],
    "max_timeout_sec":3600, // optinal
}


# Step 4
forwarding-server send back a json to A
{
    "sync_id":id3
    "req":"handshake_3"
    "handshake_result":"OK"
}


# Step 5
A send data to forwarding-server who fowarding the data to B (all recv_channel_ids receivers) directly.


# Step 6
A close the connection, then forwarding-server clear this session infomations.



