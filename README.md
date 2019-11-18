# portcheck
A tool to scan through a port range to see if any traffic gets dropped invisibly


You've turned firewall off, checked netstat and nothing is already bound, but your packets still get dropped?


Check to see if something else might be getting in the way with portcheck.

Some deep dark networking layer of your cloud provider might be eating all the packets you ever send on some high port. This tool helps to make debugging that problem easy. It enumerates though the ports you specify to check if traffic might be getting dropped.

## Usage:

Set up the server:

```
portcheck-server -p 2000 -u > /tmp/serverlogs.out
```


Then, connect via the port range specified by your client:

```
portcheck -p 30000-50000  -u  [myhost.mynet]:2000 > /tmp/clientlogs.out
```

The client sends a quick "Hi I'm connecting from [IP]:[port]" message, which the server logs, and the client logs any port that is unable to connect.


```
portcheck(master*)$ ./bin/portcheck -p 30000-50000 -u localhost:2000
FAILED: local UDP port 34611 failed to access 127.0.0.1:2000
FAILED: local UDP port 34644 failed to access 127.0.0.1:2000
FAILED: local UDP port 41313 failed to access 127.0.0.1:2000
FAILED: local UDP port 45855 failed to access 127.0.0.1:2000
FAILED: local UDP port 47109 failed to access 127.0.0.1:2000
FAILED: local UDP port 48854 failed to access 127.0.0.1:2000
```

Note that ports that fail to bind (in the above example, I had a process already bound to all of the listed ports) still fail - you might compare this list against netstat output:

```
portcheck(master*)$ netstat -anup  | grep ESTABLISHED
(Not all processes could be identified, non-owned process info
 will not be shown, you would have to be root to see it all.)
udp6       0      0 2606:a000:1125:43:48854 2404:6800:4008:c00::443 ESTABLISHED 27187/chrome --type 
udp6       0      0 2606:a000:1125:43:36705 2404:6800:4004:80b::443 ESTABLISHED 27187/chrome --type 
udp6       0      0 2606:a000:1125:43:39457 2404:6800:4004:806::443 ESTABLISHED 27187/chrome --type 
udp6       0      0 2606:a000:1125:43:41313 2404:6800:4003:c04::443 ESTABLISHED 27187/chrome --type 
udp6       0      0 2606:a000:1125:43:42574 2404:6800:4004:806::443 ESTABLISHED 27187/chrome --type 
udp6       0      0 2606:a000:1125:43:60791 2404:6800:4004:818::443 ESTABLISHED 27187/chrome --type 
udp6       0      0 2606:a000:1125:43:45643 2404:6800:4003:c04::443 ESTABLISHED 27187/chrome --type 
```

As a debugging tool, it's critial that you use this on a reliable network, and without any other network activities running which might block port access in an expected way.



The client can test all of udp, tcp, and multiple hosts at a time, if you set up multiple servers:

```
portcheck-server -p 2000 -u > /tmp/server2000udp.out&
portcheck-server -p 2001 -t > /tmp/server2001tcp.out&

portcheck -p 30000-50000 -t -u localhost:2000 localhost:2001 > /tmp/client.out
```

Note that, in the above example, all UDP connections to port 2001 would fail, and all TPC connections to port 2000 would fail.
