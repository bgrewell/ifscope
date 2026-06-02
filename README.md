# ifscope

> Unified Linux network interface inspection.

`ifscope` is a Linux CLI for inspecting host network state in one place. It
correlates interface, VLAN, route, DNS, PCIe, driver, firmware, Open vSwitch,
SR-IOV, VF, and basic connectivity data into readable tables and JSON output —
replacing the "run six commands and stitch the output together" workflow.

```bash
ifscope                 # interfaces + VLANs (default view)
ifscope all             # everything
ifscope interfaces      # interface table
ifscope sriov           # SR-IOV / VF state
ifscope test            # basic connectivity tests
ifscope --json          # machine-readable output
```

## Status

Early development. See `.planning/rel-1-plan.md` for the REL-1 scope and
milestones. The CLI is read-only by design and degrades gracefully when
optional tools (`ethtool`, `lspci`, `resolvectl`, `ovs-vsctl`) are absent.

## Build from source

Requires Go 1.26+.

```bash
make build      # builds ./bin/ifscope with version metadata
make test       # run unit tests
make check      # fmt check, vet, lint, test
```

## License

To be determined before the first public release.
