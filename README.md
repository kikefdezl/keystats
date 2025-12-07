# KeyStats

Listen to keyboard inputs and log them to a JSON file.

This program accesses `/dev/input` files to listen for keystrokes, so root permissions are needed to run it.

To run it:

```commandline
make build
sudo ./keystats
```

The keypresses will be accumulated and saved every 10 minutes into `keystats.json`
