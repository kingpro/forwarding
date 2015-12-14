# forwarding
A very liteweight tool to forward data over tcp,  written in Go.

这有点像一个turn服务, 多个客户端连接到这个服务上，它负责转发任意两端(多端)之间的tcp层报文, 所以这必然是一个神器。具体能干嘛? 你懂的，嘿嘿~


#TODO:
```
1 握手报文加密
2 超时关闭
3 自动生成chid
```


#Client A want to broadcast data to client B, the process is
```
client A [TCP stream encrypted by AES256]----> forwarding-server ----> [TCP stream encrypted by AES256] client B
```


# Step 1 
client send the following json to forwarding-server: 

```
{
    "req":"hs1",  // hs1 a.k.a first step of handshake process
    "chid":"ch0" // optional. Sever will generate if not set.
}
```


e.g.
client A send the following json to forwarding-server: 

```
{"req":"hs1","chid":"ch0"}
```

client B send the following json to forwarding-server:

```
 {"req":"hs1","chid":"ch1"}

```


# Step 2 
forwarding-server send back the following json to the client

```
{
    "rsp":"hs2",   
    "chid":"ch1" 
}

```


# Step 3 
client send the following json to forwarding-server: 

```
client send a json
{  
    "req":"hs3",   
    "recvers":["ch2","ch3","ch4"],   // optional, not set means broadcast to all others.
    "timeout":3600, // optinal, seconds
}

```


e.g.

client A: 

```
{"req":"hs3","recvers":["ch1"]}

```
client B:

```
 {"req":"hs3","recvers":["ch0"]}

```


# Step 4
forwarding-server send back the following json to the client

```
{   
    "rsp":"hs4",   
    "result":"OK" // "OK" means initial sucessfully, "FAIL" means initial failed.
    "desc":"" // optional.
}
```


# Step 5
client send data to forwarding-server who will foward the data to all receivers directly.


# Step 6
A close the connection, then forwarding-server clear this session infomations.

