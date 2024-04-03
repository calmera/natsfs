# NATS Jetstream FUSE Driver
How cool would it be to mount a NATS JetStream object store as a filesystem? This project aims to do just that. It is a FUSE 
driver that allows you to mount a NATS JetStream object store as a filesystem.

## About
This project is more of a proof of concept than a production-ready tool. The concept to prove being the ability to mount A JetStream
Object Store directly into the linux filesystem using FUSE. I am happy to say that after an evening of bantering, swearing and
 profanities, we are able to do just that. It is written in Go and uses the https://github.com/hanwen/go-fuse 
library to make it easier to interact with the low-level stuff.  

## How to use
While it is probably possible to just reference this in your fstab in some way, I have not tried that yet. Instead You
can run the binary directly. The binary takes a few arguments to configure the connection to the NATS server.
```
natsfs -url=<url> -creds=<creds> -debug <mountpoint>
```

Should be pretty self explanatory. The `-url` argument is the URL to the NATS server, the `-creds` argument is the path
to the credentials file, and the `<mountpoint>` is the directory where you want to mount the filesystem. The `-debug`
flag will enable debug logging and is optional.

## Disclaimer
Use at your own risk. This project is not production-ready and is not recommended for use in production environments. 
It might be unstable and could potentially cause data loss, send inappropriate messages, or even summon a demon. 
Complaints will be ignored.
