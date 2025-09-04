# GoPOST
GoPOST is a proof-of-concept DDoS tool designed to exploit unprotected HTTP POST endpoints to overwhelm the target server with data.

The tool has built-in proxy support, but requires self-compiling with your proxies in the ``proxies.txt``  file. ``This tool will not run without proxies``, as it would get ratelimited extremely quickly by the target server. I have packaged some example internal proxies with the project.

GoPOST is multithreaded, fast, and efficient, with various CLI options designed to facilitate versatility and ease-of-use. Please consult the application's ``help`` output for further details.