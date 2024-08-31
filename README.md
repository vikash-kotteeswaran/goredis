// Serve By Multiplexing
//
// 1) Subscribe server fd to multiplexer so that when polled at first, only server fd would be returned as available
// from ePolling iff server fd is free.
//
// 2) When server fd is obtained from ePolling, accept connections and subscribe the connection fds to the mulitplexer.
//
// 3) When ePolled again, if server fd is obtained again then connections are accepted again and
// the connections are subscribed again. If ePolling returns connection fds as well, process them too.
//
// 4) Simultaneously connections are accepted from subscribed server fd and already suscribed connection fds are used
// to serve the respective connection requests
