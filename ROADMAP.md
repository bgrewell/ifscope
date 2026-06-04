# ifscope roadmap

Planned views beyond the current set. Effort is a rough guide: **S** ≈ a
collector + parser + table (a few hours), **M** ≈ S plus correlation or an
optional dependency, **L** ≈ multiple sources or per-namespace/entering work.

All new views follow the existing pattern (pure `parse`, `collect` behind
`run.Runner`/`sysfs.FS`, `render` table + JSON) and degrade gracefully when an
optional tool is absent (warning, not failure).

## Available now

`show` (interfaces + VLANs + bonds + bridges), `interfaces`, `vlans`, `bonds`,
`bridges`, `pcie` (incl. DPDK/detached NICs), `routes` (all tables), `rules`
(policy routing), `neighbors` (ARP/NDP), `dns`, `ovs`, `sriov`, `stats`,
`netns`, `test`, `all`.

## Tier 2 — planned

| View | Source | Key fields | Effort | Notes |
| --- | --- | --- | --- | --- |
| `lldp` — link-layer neighbors | `lldpcli show neighbors -f json` (lldpd) | local port, remote chassis/sysname, port id/desc, mgmt IP, VLAN | M | Optional dep (lldpd); huge for datacenter topology. Degrade if absent. |
| `tunnels` — overlay/tunnel endpoints | `ip -d -json link` (info_data) + `wg show <dev> dump` | type (vxlan/gre/geneve/wireguard), local, remote, VNI/key, port, ttl | M | WireGuard peers need `wg` (+root). Builds on the type classification already done. |
| `devlink` — device/eswitch | `devlink -j dev`, `devlink -j port` | eswitch mode (legacy/switchdev), port flavour, pf/vf number | M | Complements SR-IOV/DPDK; switchdev is key for offload setups. |
| `fdb` — bridge forwarding DB | `bridge -json fdb show` | mac, dev/port, vlan, flags, state | S | Pairs with the bridge view. |
| bridge VLANs | `bridge -json vlan show` | port, vlan ids, pvid, flags | S | Could fold into the `bridges` view or `fdb`. |
| `sockets` — listening ports | `ss -tulpnH` | proto, local addr:port, state, process (pid/name) | M | Process needs root. Slightly adjacent to interface inspection — include behind its own command. |

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
