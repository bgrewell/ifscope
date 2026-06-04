package parse

import "testing"

func TestInterrupts(t *testing.T) {
	data := []byte(` 102:   0   0   19   0   IR-PCI-MSIX-0000:17:00.0    0-edge      ice-0000:17:00.0:misc
 104:   0   0    0   0   IR-PCI-MSIX-0000:17:00.0    2-edge      ice-enp23s0np0-TxRx-0
 NMI:   1   2    3   4   Non-maskable interrupts
 LOC:   5   6    7   8   Local timer interrupts`)
	irqs := Interrupts(data)
	if len(irqs) != 2 {
		t.Fatalf("irqs = %d, want 2 (NMI/LOC skipped)", len(irqs))
	}
	if irqs[0].Number != 102 || irqs[0].Name != "ice-0000:17:00.0:misc" {
		t.Errorf("irq0 = %+v", irqs[0])
	}
	if irqs[1].Number != 104 || irqs[1].Name != "ice-enp23s0np0-TxRx-0" {
		t.Errorf("irq1 = %+v", irqs[1])
	}
}
