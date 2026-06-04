# ifscope roadmap

Planned views beyond the current set. Effort is a rough guide: **S** ≈ a
collector + parser + table (a few hours), **M** ≈ S plus correlation or an
optional dependency, **L** ≈ multiple sources or per-namespace/entering work.

All new views follow the existing pattern (pure `parse`, `collect` behind
`run.Runner`/`sysfs.FS`, `render` table + JSON) and degrade gracefully when an
optional tool is absent (warning, not failure).

## Available now

`show` (interfaces + VLANs + bonds + bridges), `interfaces`, `vlans`, `bonds`,
`bridges`, `pcie` (incl. DPDK/detached NICs), `devlink`, `routes` (all tables),
`rules` (policy routing), `neighbors` (ARP/NDP), `fdb`, `dns`, `ovs`, `sriov`,
`stats`, `sockets`, `netns`, `test`, `all`.

## Tier 2 — complete

All Tier 2 views are delivered: `devlink`, `fdb`, `sockets`, `tunnels`, `lldp`,
`wireguard`, `bridge-vlans`.

Notes from live validation (in LXD):
- The `tunnels` view parses `ip -d link` **text**, not `-json` — iproute2
  emits malformed JSON for vxlan (a stray `fan-map` token). Verified with live
  vxlan/gre/geneve interfaces.
- The `lldp` view uses `lldpcli -f json0` (array-wrapped, unambiguous) rather
  than `-f json` (object/array-ambiguous). Verified with two lldpd nodes on a
  bridge with LLDP forwarding enabled.

## Tier 3 — planned (advanced / niche)

Delivered so far: `qdisc` (root discipline per device), `offloads`.

| View | Source | Key fields | Effort | Notes |
| --- | --- | --- | --- | --- |
| `queues` — channels / RSS / rings | `ethtool -l` (channels), `-x` (RSS), `-g` (ring), `-c` (coalesce) | combined/rx/tx queues, RSS indirection, ring sizes, coalesce usecs | M | Per-NIC; several ethtool calls. |
| IRQ / RPS / XPS affinity | `/proc/interrupts`, `/proc/irq/*/smp_affinity_list`, `/sys/class/net/*/queues/*/{rps,xps}_cpus` | per-queue IRQ→CPU, RPS/XPS CPU masks | L | Correlate IRQs to NIC queues to CPUs/NUMA. |
| multicast / MDB | `ip -json maddr`, `bridge -json mdb` | group, dev, users; bridge mdb entries | S | |
| `wifi` | `iw dev`, `iw <dev> link` | SSID, freq, signal, bitrate, mode | S | Needs `iw`; low priority for a server tool. |
| PTP / hw timestamping | `ethtool -T <dev>` | PHC index, tx/rx timestamp capabilities | S | Telco/finance niche. |
| qdisc class/filter hierarchy | `tc -json class/filter show dev <d>` | class rate/ceil, filter rules | M | Deepens the `qdisc` view beyond the root discipline. |

## Cross-cutting / future

- **Per-namespace inspection** — run any view inside a chosen netns (`ip -n` /
  `nsenter`); the `netns` view currently lists namespaces + interface counts.
- **`--watch` / interval mode** — refresh a view on an interval (rates for `stats`).
- **Prometheus exporter** — expose counters/state for scraping.

## Out of scope (for now)

- Firewall (nftables/iptables) and conntrack/NAT — large, stateful, and a
  distinct concern from read-only inspection. Revisit only if the mission
  broadens.
