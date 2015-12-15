# forwarding
A very liteweight tool to forward data over tcp,  written in Go.

本项目有点像一个turn服务, 多个客户端连接到这个服务上，它负责转发任意两端(多端)之间的tcp层报文, 所以这必然是一个神器。

所有客户端连上服务握手成功以后会拿到通道标识即 chid， 握手过程中设置好接收者的chid列表，后续发给服务的所有数据都会被转发到设置的接收者那里。

握手报文可以设置加密，负载数据是否加密服务不关心，它只是原样转发不作解析。


具体能干嘛? 你懂的，嘿嘿 ^=^



#TODO:
```
1 握手报文加密
2 超时关闭
3 自动生成chid
4 chan的缓冲等完善
5 架构完善
6 控制指令完成
```


#Client A want to send data to client B, the process is
```
client A [TCP stream encrypted by AES256]----> forwarding-server ----> [TCP stream encrypted by AES256] client B
```



#The concrete steps of all the steps is:


# Step 1 , this is first step of handshake process
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


# Step 2  , this is second step of handshake process
forwarding-server send back the following json to the client

```
{
    "rsp":"hs2",   
    "chid":"ch1" 
}

```


# Step 3  , this is third step of handshake process
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


# Step 4 , this is forth step of handshake process
forwarding-server send back the following json to the client

```
{   
    "rsp":"hs4",   
    "result":"OK" // "OK" means initial sucessfully, "FAIL" means initial failed.
    "desc":"" // optional.
}
```


# Step 5 , send data after handshake sucessfully
client send data to forwarding-server who will foward the data to all receivers directly.


# Step 6 , close session when everthing is done
client close the connection, then forwarding-server clear this session infomations.

