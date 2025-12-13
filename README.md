# KeyStats

Listen to keyboard inputs and log them to a JSON file.

This program accesses `/dev/input` files to listen for keystrokes, so root permissions are needed to run it.

To run it:

```commandline
make build
sudo ./keystats
```

The keypresses will be accumulated and saved every 10 minutes into `keystats.json`

## Launch on Startup

1. Put the binary in `/usr/local/bin/keystats`.

2. Create a systemd entry at `/etc/systemd/system/keystats.service` with:
```toml
[Unit]
Description=Keystats

[Service]
ExecStart=/usr/local/bin/keystats
User=root
Group=root

[Install]
WantedBy=multi-user.target
```
3. Enable:
```commandline
sudo systemctl enable keystats
```
