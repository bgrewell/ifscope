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

## Tier 3 — core delivered

Delivered: `qdisc` (root discipline per device), `offloads`, `queues`
(channels + ring sizes), `irq` (NIC interrupt CPU affinity), `multicast`
(IP group memberships).

Remaining / optional (lower priority):

| View | Source | Notes |
| --- | --- | --- |
| RSS / coalesce (extend `queues`) | `ethtool -x` (RSS indirection), `-c` (coalesce) | Adds RSS table + interrupt coalescing to the queues view. |
| RPS / XPS (extend `irq`) | `/sys/class/net/*/queues/*/{rps,xps}_cpus` | Software steering masks; usually disabled (0). |
| bridge MDB (extend `multicast`) | `bridge -json mdb` | Multicast forwarding DB; needs a snooping bridge to exercise. |
| qdisc class/filter hierarchy | `tc -json class/filter show dev <d>` | Deepens `qdisc` beyond the root discipline. |
| `wifi` | `iw dev`, `iw <dev> link` | Niche for a server tool; needs a wireless NIC. |
| PTP / hw timestamping | `ethtool -T <dev>` | Telco/finance niche. |

## Cross-cutting / future

- **Per-namespace inspection** — run any view inside a chosen netns (`ip -n` /
  `nsenter`); the `netns` view currently lists namespaces + interface counts.
- **`--watch` / interval mode** — refresh a view on an interval (rates for `stats`).
- **Prometheus exporter** — expose counters/state for scraping.

## Out of scope (for now)

- Firewall (nftables/iptables) and conntrack/NAT — large, stateful, and a
  distinct concern from read-only inspection. Revisit only if the mission
  broadens.
