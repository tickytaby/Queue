Flash Sale Queue Backend Service

1. Project Goal:
To design and implement an efficient and resilient backend service capable of handling high levels of concurrent buyer load, typical of a flash sale, or limited inventory.

2. Skills and Tech to Learn:
- Goroutines: understanding an implementation of concurrent execution units.
- Channels: to ensure safe communication and synchronization between goroutines.
- WebSockets: implementing full-duplex, persistent communication between server and client. Provide real time updates.

3. Planned Architectural Approach:
- HTTP Server: handling of initial client connections and serve static content.
- Websocket endpoint: upgrade the HTTP connection to support websockets for real time communication with clients.
- Buffered channel to act as a queue (order ingestion).
- Worker pool: set of goroutines that will read the order requests and push to the buffered channel.
- Queue Monitor: a single goroutine to handle the orders in the buffered channel, to ensure FIFO.