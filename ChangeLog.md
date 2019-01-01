#新的变化

2019-01-01 支持`jsonrpc-glib-1.0`的`notify`消息，调用方法 `Notify(method string, params interface{})`.示例：

```
conn,_ :=  l.Accept()
myconn := jsonrpc2glib.NewMyConn(conn)
go rpc.ServeCodec(jsonrpc2.NewServerCodec(myconn, nil))
myconn.Notify("SomeMethod",params)
```
