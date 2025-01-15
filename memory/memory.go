package memory

import "fmt"

const (
	ROMStart          uint16 = 0x0000
	ROMEnd            uint16 = 0x7FFF
	VRAMStart         uint16 = 0x8000
	VRAMEnd           uint16 = 0x9FFF
	ExternalRAMStart  uint16 = 0xA000
	ExternalRAMEnd    uint16 = 0xBFFF
	InternalRAM0Start uint16 = 0xC000
	InternalRAM0End   uint16 = 0xCFFF
	InternalRAM1Start uint16 = 0xD000
	InternalRAM1End   uint16 = 0xDFFF
	OAMStart          uint16 = 0xFE00
	OAMEnd            uint16 = 0xFE9F
	IOPortsStart      uint16 = 0xFF00
	IOPortsEnd        uint16 = 0xFF7F
	HRAMStart         uint16 = 0xFF80
	HRAMEnd           uint16 = 0xFFFF
)

// Memory structure
type Memory struct {
	rom  []byte       // ROM Data
	vram [0x2000]byte // Video RAM
	ram  [0x2000]byte // Internal RAM (0 + 1)
	oam  [0xA0]byte   // OAM
	io   [0x80]byte   // I/O Ports
	hram [0x80]byte   // High RAM
}

// NewMemory initializes the Memory structure
func NewMemory(rom []byte) *Memory {
	m := &Memory{
		rom: rom,
	}
	return m
}

// Read retrieves the value at a given address
func (m *Memory) Read(addr uint16) byte {
	switch {
	case addr >= ROMStart && addr <= ROMEnd:
		// Read from ROM
		return m.rom[addr]
	case addr >= VRAMStart && addr <= VRAMEnd:
		// Read from Video RAM
		return m.vram[addr-0x8000]
	case addr >= ExternalRAMStart && addr <= ExternalRAMEnd:
		// Read from External RAM (if implemented)
		return m.ram[addr-0xA000]
	case addr >= InternalRAM0Start && addr <= InternalRAM0End:
		// Read from Internal RAM 0
		return m.ram[addr-0xC000]
	case addr >= InternalRAM1Start && addr <= InternalRAM1End:
		// Read from Internal RAM 1
		return m.ram[addr-0xD000]
	case addr >= OAMStart && addr <= OAMEnd:
		// Read from OAM
		return m.oam[addr-0xFE00]
	case addr >= IOPortsStart && addr <= IOPortsEnd:
		// Read from I/O Ports
		return m.io[addr-0xFF00]
	case addr >= HRAMStart && addr <= HRAMEnd:
		// Read from High RAM
		return m.hram[addr-0xFF80]
	default:
		// Handle invalid memory access
		fmt.Printf("Invalid memory read at address: %04X\n", addr)
		return 0xFF // Return a default value for invalid access
	}
}

// Write sets the value at a given address
func (m *Memory) Write(addr uint16, value byte) {
	switch {
	case addr >= ROMStart && addr <= ROMEnd:
		// ROM should be read-only in most cases, do nothing or handle it
		fmt.Printf("Invalid write to ROM at address: %04X\n", addr)
	case addr >= VRAMStart && addr <= VRAMEnd:
		// Write to Video RAM
		m.vram[addr-0x8000] = value
	case addr >= ExternalRAMStart && addr <= ExternalRAMEnd:
		// Write to External RAM (if implemented)
		m.ram[addr-0xA000] = value
	case addr >= InternalRAM0Start && addr <= InternalRAM0End:
		// Write to Internal RAM 0
		m.ram[addr-0xC000] = value
	case addr >= InternalRAM1Start && addr <= InternalRAM1End:
		// Write to Internal RAM 1
		m.ram[addr-0xD000] = value
	case addr >= OAMStart && addr <= OAMEnd:
		// Write to OAM
		m.oam[addr-0xFE00] = value
	case addr >= IOPortsStart && addr <= IOPortsEnd:
		// Write to I/O Ports
		m.io[addr-0xFF00] = value
	case addr >= HRAMStart && addr <= HRAMEnd:
		// Write to High RAM
		m.hram[addr-0xFF80] = value
	default:
		// Handle invalid memory access
		fmt.Printf("Invalid memory write at address: %04X\n", addr)
	}
}
