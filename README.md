# ifscope

> Unified Linux network interface inspection.

`ifscope` is a Linux CLI that correlates host network state — interfaces,
VLANs, routes, DNS, PCIe NIC mapping, drivers, firmware, Open vSwitch topology,
SR-IOV/VF state, and basic connectivity checks — into readable tables and
stable JSON. It replaces the "run six commands and mentally stitch the output
together" workflow with one consistent tool.

```bash
ifscope                 # interfaces + VLANs (default view)
ifscope all             # every table
ifscope interfaces      # interface table (driver, firmware, bus, speed, port, SR-IOV)
ifscope pcie            # PCIe device mapping
ifscope routes          # routing table
ifscope dns             # resolver state
ifscope ovs             # Open vSwitch bridges/ports/VLANs
ifscope sriov           # SR-IOV PF/VF state
ifscope test            # basic connectivity tests
ifscope --json          # machine-readable output
```

## Why not just use `ip`, `ethtool`, `lspci`, …?

Those tools each show one slice of the picture. `ifscope` runs them for you and
correlates the results: a single interface row carries its link state, IPv4
addresses, driver, firmware, PCI bus, speed, port type, alt names, and SR-IOV
summary. The same data is available as stable JSON for automation.

`ifscope` is **read-only by design** — it never changes host configuration.

## Quick start

```bash
# Build (Go 1.24+)
make build
./bin/ifscope

# Or install onto your PATH
make install
ifscope --help
```

## Example output

```text
Interfaces
┌────┬────────────┬───────────────────┬───────┬───────────────┬────────┬──────────┬─────────┬──────────┬──────────────────────┬──────────┬───────┬─────┐
│ ID │ NAME       │ MAC               │ STATE │ ADDRS         │ DRIVER │ FIRMWARE │ BUS     │ SPEED    │ PORT                 │ ALTNAMES │ SRIOV │ VFS │
├────┼────────────┼───────────────────┼───────┼───────────────┼────────┼──────────┼─────────┼──────────┼──────────────────────┼──────────┼───────┼─────┤
│ 6  │ enp23s0np0 │ 40:a6:b7:ca:3a:08 │ UP    │ 192.0.2.10/24 │ ice    │ 4.30 …   │ 17:00.0 │ 100 Gb/s │ Direct Attach Copper │ …        │ 0/256 │     │
└────┴────────────┴───────────────────┴───────┴───────────────┴────────┴──────────┴─────────┴──────────┴──────────────────────┴──────────┴───────┴─────┘
```

## Commands

| Command | Description |
| --- | --- |
| `ifscope` / `ifscope show` | Interfaces + VLANs (default view) |
| `ifscope interfaces` | Interface table with driver/firmware/bus/speed/port/SR-IOV |
| `ifscope vlans` | VLAN interfaces (parent, tag, addresses) |
| `ifscope bonds` | Bonding masters with mode, active slave, and members |
| `ifscope bridges` | Linux bridges with STP state and member ports |
| `ifscope bridge-vlans` | Bridge VLAN-filtering entries per port (alias `bvlan`) |
| `ifscope tunnels` | Overlay/tunnel interfaces (VXLAN, GENEVE, GRE, …) |
| `ifscope wireguard` | WireGuard interfaces and peers (alias `wg`) |
| `ifscope pcie` | PCIe network devices (driver, kernel binding, vendor/device, NUMA, link) |
| `ifscope devlink` | devlink ports (PF/VF flavour, switchdev) |
| `ifscope routes` | Routing tables (all tables, with the table name) |
| `ifscope rules` | Routing policy rules (source-based / policy routing) |
| `ifscope neighbors` | ARP/NDP neighbor table (alias `arp`) |
| `ifscope lldp` | LLDP link-layer neighbors (chassis, port, mgmt IP) |
| `ifscope fdb` | Bridge forwarding database (MAC table) |
| `ifscope dns` | Per-link and global resolver state |
| `ifscope ovs` | Open vSwitch bridges, ports, and VLAN tags |
| `ifscope sriov` | SR-IOV PF/VF state |
| `ifscope stats` | Per-interface traffic and error counters |
| `ifscope sockets` | Listening TCP/UDP sockets (alias `ss`) |
| `ifscope netns` | Network namespaces |
| `ifscope test` | Gateway / internet / DNS ping, HTTPS GET (+ optional throughput) |
| `ifscope all` | Every inspection table |
| `ifscope version` | Build metadata |

Planned additional views (LLDP, tunnels, devlink, qdisc/shaping, and more) are
tracked in [ROADMAP.md](ROADMAP.md).

### Global flags

```text
--json            emit JSON instead of tables
--pretty          pretty-print JSON
--summary, -s     shorter summary tables
--barebones, -b   plain pipe-delimited tables
--no-color        disable color
--color           auto|always|never
--warnings        always show collection warnings
--up, -u          only interfaces that are UP
--interface NAME  filter to one interface
--driver NAME     filter by driver
--state STATE     filter by up|down|unknown
--physical / --virtual / --pf / --vf   interface-class filters
--ovs / --no-ovs  include / skip Open vSwitch data
--no-sudo         never escalate with sudo (OVS)
```

Legacy `netcheck` short flags are preserved as aliases: `-I` interfaces,
`-V` vlans, `-D` dns, `-R` routes, `-P` pcie, `-t` test, `-j` json, `-u` up,
`-s` summary. They can be combined, e.g. `ifscope -IP`.

## JSON output

`--json` emits a single stable document; warnings are carried inside it
(`warnings[]`), never mixed into the stream. The top-level `version` field
versions the schema.

```bash
ifscope all --json | jq '.interfaces[] | {name, driver, speed, sriov}'
```

```json
{
  "version": "1.0.0",
  "host": {"hostname": "example"},
  "interfaces": [
    {
      "id": 6, "name": "ens5f0np0", "type": "physical",
      "mac": "50:7c:6f:5a:f8:c8", "state": "UP",
      "addresses": [{"family": "inet", "local": "192.0.2.10", "prefixlen": 24}],
      "driver": "ice", "firmware": "4.30", "bus": "0000:25:00.0",
      "speed": "100 Gb/s", "port": "Direct Attach Copper",
      "sriov": {"capable": true, "total_vfs": 128, "configured_vfs": 0, "vf": false}
    }
  ],
  "warnings": []
}
```

The human table shows IPv4 by default; JSON always includes IPv6 addresses too.

## SR-IOV / VF inspection

`ifscope sriov` shows physical functions (PFs), how many VFs are configured of
the maximum supported, and per-VF bus, driver binding (kernel module,
`vfio-pci`, or unbound), netdev, MAC, VLAN, spoof-check, trust, and link state.
The interface table summarizes SR-IOV as `configured/total` for PFs and
`VF of <pf>` for VFs. Use `--pf` / `--vf` to filter.

## Open vSwitch support

`ifscope ovs` lists bridges, their ports, the interfaces on each port, and
access VLAN tags / trunk lists. OVS reads are attempted unprivileged first and
retried with `sudo -n` when access is denied; pass `--no-sudo` to disable
escalation. `--ovs` adds an OVS section (and per-interface membership in JSON)
to the `show`, `interfaces`, and `all` views.

## PCIe and DPDK / detached NICs

`ifscope pcie` scans `/sys/bus/pci/devices` for Ethernet-class controllers
rather than only kernel netdevs, so it also surfaces NICs that have been
**detached from the kernel** — e.g. bound to `vfio-pci`/`uio_pci_generic` for
DPDK, or left unbound. The `BIND` column classifies each device:

- `kernel` — a normal kernel netdev is present
- `dpdk` — bound to a userspace/passthrough driver, no netdev
- `unbound` — no driver bound
- `detached` — a driver is bound but exposes no netdev

## Routing tables and policy rules

`ifscope routes` collects from **all** routing tables (`ip route show table all`),
not just `main`, and shows each route's table — so `local`, `default`, and any
custom tables used by policy routing are visible.

`ifscope rules` shows the routing policy rules (`ip rule`): priority, match
`from`/`to`, input/output interface, fwmark, and the table each rule selects.
Together these represent **source-based / policy routing**, which the main table
alone doesn't reveal.

## Connectivity tests

`ifscope test` runs: gateway ping (resolved via the route to the ping target),
internet ping (no DNS), DNS-name ping, and an HTTPS GET. The throughput
download is **opt-in** (`--throughput`) to avoid surprise large transfers.
Targets are configurable:

```bash
ifscope test --ping-target 8.8.8.8 --dns-target example.com \
  --web-target https://example.com/ --count 2 --timeout 3s
ifscope test --throughput
```

`ifscope test` exits non-zero (10) if any test fails.

## Dependencies

| Tool | Required? | Used for |
| --- | --- | --- |
| `ip` (iproute2) | required | interfaces, VLANs, routes, VF attributes |
| `ethtool` | optional | driver, firmware, speed, port |
| `lspci` (pciutils) | optional | PCIe device descriptions, vendor/device IDs |
| `resolvectl` (systemd) | optional | DNS resolver state |
| `ovs-vsctl` | optional | Open vSwitch topology |
| `ping` (iputils) | optional | connectivity ping tests |
| sysfs (`/sys`) | built-in | PCIe attributes, SR-IOV PF/VF state |

Missing optional tools degrade gracefully: affected fields are left empty and a
warning is recorded (shown with `--warnings`, or in JSON `warnings[]`).

## Permissions

`ifscope` runs as a normal user. Some data needs elevated privileges:

- **OVS**: `ovs-vsctl` typically requires root; `ifscope` auto-retries with
  `sudo -n` (disable with `--no-sudo`).
- **Some ethtool/SR-IOV fields** may be unavailable without privileges; they
  are reported as empty rather than failing the command.

## Exit codes

| Code | Meaning |
| --- | --- |
| 0 | success |
| 1 | general error |
| 2 | invalid arguments |
| 10 | one or more connectivity tests failed |
| 20 | a required dependency was missing |
| 30 | permission denied for a requested collector |

Inspection commands return 0 even when optional data is missing.

## Known limitations

- Linux only.
- Bonds are summarized (mode, active slave, members); deeper team/bridge
  inspection, network namespaces, and VRF are out of scope for this release.
- Connectivity tests are reachability checks, not a network benchmark.

## Build from source

```bash
make build      # ./bin/ifscope with version metadata
make test       # unit tests
make check      # gofmt check, vet, lint, test
```

## License

Licensed under the [Apache License, Version 2.0](LICENSE). Attribution must be
preserved per the [NOTICE](NOTICE) file when redistributing or creating
derivative works.
