# Port Randomizer

Port Randomizer is a simple Go application designed to generate random ports and list active ports.


## Usage

### Randomize Ports

```bash
port-randomizer randomize
```

Generates a random port.

#### Flags:

- `--silent` (`-s`): Silence gamified output.

### List Active Ports

```bash
port-randomizer list-active [--tcp] [--udp] [--all]
```

Lists active ports based on the specified protocol.

#### Flags:
- `--tcp` (`-t`): List used TCP ports.
- `--udp` (`-u`): List used UDP ports.
- `--all` (`-a`): List all used ports (TCP and UDP).


## Check Well-Known Ports

```bash
port-randomizer check-well-known
```

Randomly selects a port from the well-known ports list.