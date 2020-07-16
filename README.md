# SSH Tail
This is a CLI app that will setup SSH connections to multiple hosts specified in the given spec file using a key of your choice, tail the named file, and aggregate the output to the calling terminal's STDOUT.

**Note:** This utility uses the `tail` executable on the remote host to facilitate its base functionality. This limitation is mostly because I haven't figured out any other way yet. PRs welcome!

## Installation
It's easy to get started, just run this command!
```bash
go get github.com/drognisep/sshtail
```

## Examples
An example file can be output to "test.yml" by running
```bash
sshtail spec init --with-comments test.yml
```

Here's the output.
```yaml
# Hosts and files to tail
hosts:
  host1:
    hostname: remote-host-1
    # Excluding the username here will default it to the current user name
    file: /var/log/syslog
    # Default SSH port
    port: 22
  host2:
    hostname: remote-host-2
    username: me
    file: /var/log/syslog
    port: 22
# This section is optional for portability
keys:
  host1:
    # Defaults to this value
    path: ~/.ssh/id_rsa
  host2:
    # If all of these values are the same, then 'sshtail usekey' may be more convenient.
    path: ~/.ssh/id_rsa
```

## Hosts
This section is used to specify the host machines to connect to. `hostname` and `file` are required, but `port` may be excluded if the default SSH port of 22 is desired.

The values of "host1" and "host2" can be anything you wish, and are primarily used to match a specified host with a given key path, and to tag the output to your terminal like so:
```
[ host1 ] A line posted to /var/log/syslog on remote-host-1...
[ host1 ] And another one...
```

## Keys
This section is entirely optional, but an entry here overrides both user home configuration and the default value, as long as the key tag (like "host1") matches up with a host tag.

If there are more `keys` entries than `hosts` entries, a warning is printed to the terminal.

## Common Commands
This will create a spec file useful for understanding the format, exactly like what is shown above.
```bash
sshtail spec init --with-comments <spec file name>
```

To make it a bit more useful, this will exclude the `keys` section for portability and won't print comments.
```bash
sshtail spec init --exclude-keys <spec file name>
```

The default path for the SSH key used is `~/.ssh/id_rsa`. If you don't want to put a `keys` section in your spec file this command can be used to override the default by placing the given path in your config file (`~/.sshkeys.yaml`).

**Note:** The given file is not currently validated as a real key. I plan on fixing this soon.
```bash
sshtail usekey /new/default/key/here
```

Finally, to execute a spec use this command. If the configured key is encrypted then the user will be asked to enter their pass phrase each time it is referenced. This is for security purposes because I don't want to cache the pass phrase in memory.
```bash
sshtail spec run <spec file name>
```