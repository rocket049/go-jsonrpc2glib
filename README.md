# go-jsonrpc2glib
Help developer create server compatibly `jsonrpc-glib-1.0` .Use `github.com/powerman/rpc-codec/jsonrpc2` inside.
帮助程序员开发兼容`jsonrpc-glib-1.0`的`jsonrpc2`服务器。

### Install
`go get github.com/rocket049/go-jsonrpc2glib`

### Example:
```
//server.go
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"strings"

	"github.com/rocket049/go-jsonrpc2glib"
)

type Arith int

type ParamsT struct {
	Arg string
}

func (t *Arith) Hello(args []string, reply *string) error {
	*reply = strings.Join(args, "\n")
	return nil
}

func main() {
	arith := new(Arith)
	rpc.Register(arith)
	l, e := net.Listen("tcp", "127.0.0.1:6666")
	defer l.Close()
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go jsonrpc2glib.ServeGlib(conn, nil)
	}

	waitSig()
}

func waitSig() {
	var c chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println("\nSignal:", s)
}

//client.vala
//valac --pkg jsonrpc-glib-1.0 jsonclient.vala
using Jsonrpc;

owned SocketConnection rpcConnect(string host,uint16 port){
	Resolver resolver = Resolver.get_default ();
	List<InetAddress> addresses = resolver.lookup_by_name (host, null);
	InetAddress address = addresses.nth_data (0);
	SocketClient client = new SocketClient ();
	SocketConnection conn = client.connect(new InetSocketAddress (address, port));
	return conn;
}

void rpcClient(SocketConnection conn){
	var c = new Jsonrpc.Client(conn);
	string[] v = {"Hello friend.","无限恐怖","abcde"};
	var params = new Variant.strv(v);
	Variant res;
	try{
		for (int i=0;i<5;i++){
			var ok = c.call("Arith.Hello",params,null,out res);
			if(ok){
				stdout.printf("%s\n",res.get_string());
			}else{
				stdout.printf("error\n");
			}
		}
	}catch (Error e) {
		stdout.printf ("Error: %s\n", e.message);
	}
	conn.close();
}

void main(){
	var conn = rpcConnect("localhost",6666);
	rpcClient(conn);
	stdout.printf("end\n");
}
```
*end*
