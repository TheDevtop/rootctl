# Rootctl
Minimal chroot launcher for OpenBSD

Usage: `rootctl [name]`

### /etc/rootctl.conf
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
