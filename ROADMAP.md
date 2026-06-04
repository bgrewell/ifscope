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

Also delivered: `queues` now includes RSS ring count, interrupt coalescing,
and RPS/XPS steering counts; `ptp` (hardware timestamping / PTP); `mdb` (bridge
multicast database); `classes` (htb/hfsc shaping classes with rate/ceil).

Remaining / optional:

| View | Source | Notes |
| --- | --- | --- |
| tc filters | `tc -json filter show dev <d>` | Complements `classes`; rules mapping traffic to classes. |
| `wifi` | `iw dev`, `iw <dev> link` | Deferred until a wireless NIC is available to test against. |

Note: `tc class` JSON needs a recent iproute2 (older versions emit text for
`class show`); ifscope degrades to a warning on those.

## Cross-cutting / future

- **Per-namespace inspection** — run any view inside a chosen netns (`ip -n` /
  `nsenter`); the `netns` view currently lists namespaces + interface counts.
- **`--watch` / interval mode** — refresh a view on an interval (rates for `stats`).
- **Prometheus exporter** — expose counters/state for scraping.

## Out of scope (for now)

- Firewall (nftables/iptables) and conntrack/NAT — large, stateful, and a
  distinct concern from read-only inspection. Revisit only if the mission
  broadens.
