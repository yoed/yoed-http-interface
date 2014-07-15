yoed-client-interface
=====================

The Yo'ed clients interface and Base client.
All the Yo'ed clients should use the `BaseYoedClient` to communicate with the Yo'ed server and read the basic configuration

#Configuration

Add a `config.json` file in the same folder than the program.

##listen (string)

The `ip:port` to listen to

##serverUrl (string)

The main server URL to connect to

##handles ([]string)

The Yo handles to monitor

#Interface
 
All the Yo'ed client must comply to the following Go interface:

````go
type YoedClient interface {
	Handle(username string)
	GetConfig() *BaseYoedClientConfig
}
````

The `Handle`method is the action to do when a Yo is received.
The `GetConfig` method returns the configuration, containing the fields described above.

#Protocol

The communication protocol with the Yo'ed server is described [here](https://github.com/yoed/yoed-server)
