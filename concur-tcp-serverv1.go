// concurrent tcp server - v1
// simple tcp chat server with concurrent go routines
// clients session ends when it send 'end'
// Author: MIchael Wolfe

package main

import (
	"bufio"
	"fmt"
	"net"
	//"os"
	"strings"
	"sync"
)


// ConnectionPool is a thread safe list of net.Conn instances
type ConnectionPool struct {
	mutex sync.RWMutex
	list  map[int]net.Conn
}


// NewConnectionPool is the factory method to create new connection pool
func NewConnectionPool() *ConnectionPool {
	fmt.Println("NewConnectionPool")
	
	pool := &ConnectionPool{
		list: make(map[int]net.Conn),
	}

	fmt.Println("before return.......")
	return pool
}

// Add collection to pool
func addv1(connection net.Conn, pool  ConnectionPool) int {
	fmt.Println("add")
	
	pool.mutex.Lock()
	
	next := len(pool.list)
        fmt.Println("len of list " , next)
	pool.list[next] = connection
	pool.mutex.Unlock()
	
	return next
}


// Get connection by id
func getv1(connectionId int,  pool ConnectionPool) net.Conn {
	fmt.Println("get .....")
	
	pool.mutex.RLock()
	connection := pool.list[connectionId]
	pool.mutex.RUnlock()
	
	return connection
}

//  Remove connection from pool
func  removev1(connectionId int, pool ConnectionPool) {
	fmt.Println("remove ....")
	
	pool.mutex.Lock()
	delete(pool.list, connectionId)
	pool.mutex.Unlock()
}

func handleConnection(nextId int,  id net.Conn, pool ConnectionPool) {
// func handleConnection(c net.Conn, connections []ConnDetails) {
	fmt.Println("func: handleConnection")

	for {
		// read from current client
		data, err := bufio.NewReader(id).ReadString('\n')	
		if err != nil {
			fmt.Println("NewReader ", err)
			return
		}
		fmt.Println("reader -> ", data)

		// when client sends 'end' close connection and remove from Pool
		temp := strings.TrimSpace(string(data))
		if temp == "end" {
			pool.mutex.RLock()
			
			fmt.Println("client ending session !")
		        id.Close()
			removev1(nextId, pool)
			
			pool.mutex.RUnlock()
			return
		}	

		// send messages to all connected clients accept current one
		for connectionid, connection := range pool.list {
			fmt.Println("id,  connection, connectionid ", id, connection, connectionid)

			if id != connection {
				fmt.Println("server,  msg: from ->  to ->  ", id, connection)
				msg := fmt.Sprintf("sending -> %s \n", data)
				writer := bufio.NewWriter(connection)
				writer.WriteString(msg)
				writer.Flush()
			}
			fmt.Println(" len pool.list -> ", len(pool.list))
		}
	}
		
	fmt.Println("end of func - server closing connection")
}

func main() {
	fmt.Println("start of concur-tcp-server .....")

	var pool  ConnectionPool
	pool.list = make(map[int]net.Conn)
	
	/*
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}
*/
//	PORT := ":" + arguments[1]
	//	l, err := net.Listen("tcp4", PORT)
	
	socket, err := net.Listen("tcp", "127.0.0.1:7777")
	
	if err != nil {
		fmt.Println(err)
		return
	}
	defer socket.Close()
	
	// prepare connection pool
	connectionPool := NewConnectionPool()
	fmt.Println("connectionPool -> ", connectionPool)

	// accept clients connections and start go routine
	for {

		connection, err := socket.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}				
		fmt.Println("remote address -> ", connection.RemoteAddr())
			
		// add conection to connection pool
		nextId := addv1(connection, pool)
		fmt.Println("next connection id -> ", nextId)

		// get next id from connection Pool
		id := getv1(nextId, pool)
		fmt.Println("next connection id -> ", id)

	        go handleConnection(nextId, id, pool)	
	}
}
