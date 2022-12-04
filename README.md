# Rootctl

Minimal chroot launcher for OpenBSD

(This repository is **archived**, because it reached steady-state a long time ago)

### Installation:

1. `git clone https://github.com/TheDevtop/rootctl.git`
2. `cd rootctl`
3. `go build`
4. (root) `cp ./rootctl /usr/local/bin/rootctl`
5. (root) `rootctl [name]`

### Configuration: /etc/rootctl.conf
```json
{
  "demo": {
    "Path": "/usr/local/demo-root",
    "Cmd": "/bin/sh",
    "Args": ["-i"],
    "Env": []
  }
}
```
