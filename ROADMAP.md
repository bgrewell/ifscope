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

## Tier 2 — planned (remaining)

| View | Source | Key fields | Effort | Notes |
| --- | --- | --- | --- | --- |
| `lldp` — link-layer neighbors | `lldpcli show neighbors -f json` (lldpd) | local port, remote chassis/sysname, port id/desc, mgmt IP, VLAN | M | Optional dep (lldpd); huge for datacenter topology. Degrade if absent. |
| WireGuard peers | `wg show all dump` | peer key, endpoint, allowed-ips, handshake | S | Extends the `tunnels` view (it already lists wireguard devices); peers need `wg` (+root). |
| bridge VLANs | `bridge -json vlan show` | port, vlan ids, pvid, flags | S | Could fold into the `bridges` view or `fdb`. |

Delivered from the original Tier 2 list: `devlink`, `fdb`, `sockets`, `tunnels`.

Note: the `tunnels` view parses `ip -d link` **text**, not `-json` — iproute2
emits malformed JSON for vxlan (a stray `fan-map` token). Verified in an LXD
container with live vxlan/gre/geneve interfaces.

## Tier 3 — planned (advanced / niche)

| View | Source | Key fields | Effort | Notes |
| --- | --- | --- | --- | --- |
| `qdisc` — traffic shaping / QoS | `tc -json qdisc/class/filter show dev <d>` | qdisc kind (fq_codel/htb/mq/…), handle, parent, class rate/ceil, filters | M | The shaping/QoS view. One call per device; summarize hierarchy. |
| `queues` — channels / RSS / rings | `ethtool -l` (channels), `-x` (RSS), `-g` (ring), `-c` (coalesce) | combined/rx/tx queues, RSS indirection, ring sizes, coalesce usecs | M | Per-NIC; several ethtool calls. |
| `offloads` — NIC features | `ethtool -k <dev>` | GRO/GSO/TSO/LRO/checksum/rx-vlan… on/off/fixed | S | Simple key/value table. |
| IRQ / RPS / XPS affinity | `/proc/interrupts`, `/proc/irq/*/smp_affinity_list`, `/sys/class/net/*/queues/*/{rps,xps}_cpus` | per-queue IRQ→CPU, RPS/XPS CPU masks | L | Correlate IRQs to NIC queues to CPUs/NUMA. |
| multicast / MDB | `ip -json maddr`, `bridge -json mdb` | group, dev, users; bridge mdb entries | S | |
| `wifi` | `iw dev`, `iw <dev> link` | SSID, freq, signal, bitrate, mode | S | Needs `iw`; low priority for a server tool. |
| PTP / hw timestamping | `ethtool -T <dev>` | PHC index, tx/rx timestamp capabilities | S | Telco/finance niche. |

## Cross-cutting / future

- **Per-namespace inspection** — run any view inside a chosen netns (`ip -n` /
  `nsenter`); the `netns` view currently lists namespaces + interface counts.
- **`--watch` / interval mode** — refresh a view on an interval (rates for `stats`).
- **Prometheus exporter** — expose counters/state for scraping.

## Out of scope (for now)

- Firewall (nftables/iptables) and conntrack/NAT — large, stateful, and a
  distinct concern from read-only inspection. Revisit only if the mission
  broadens.
