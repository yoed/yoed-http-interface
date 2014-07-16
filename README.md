Yo'ed HTTP interface
=====================

The Yo'ed HTTP interface.
All the Yo'ed handlers should use the `Client` to communicate with the Yo'ed server over HTTP (only protocol supported now).

#Configuration

Add a `config.json` file in the same folder than the program.

##listen (string)

The `ip:port` to listen to

##serverUrl (string)

The main server URL to connect to

##handles ([]string)

The Yo handles to monitor

#Interface
 
All the Yo'ed handlers must comply to the following Go interface:

````go
type Handler interface {
	Handle(username string)
}
````

The `Handle` method is the action to execute when a Yo is received.

#Protocol

The communication protocol with the Yo'ed server is described [here](https://github.com/yoed/yoed-server/blob/master/README.md#protocol)
