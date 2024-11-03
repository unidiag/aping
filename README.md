# aping

## Fast scanner for your /24 network in 1 second.

### Compile:
```
go build -ldflags "-linkmode external -extldflags '-static'" -o aping
```

### Usage:
```
root@nuc:/aping# ./aping
10.8.0.1
10.8.0.6
10.8.0.26
Took: 1.112s
root@nuc:/aping# ./aping 192.168.1.0
192.168.1.1
192.168.1.20
192.168.1.25
192.168.1.150
192.168.1.152
192.168.1.153
192.168.1.155
192.168.1.156
192.168.1.159
Took: 1.105s
root@nuc:/aping# 
```
