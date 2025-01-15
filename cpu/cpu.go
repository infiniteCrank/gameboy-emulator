package cpu

import (
	"fmt"
)

// Define the CPU structure with registers and flags
type CPU struct {
	A, F   byte   // Accumulator and Flags
	B, C   byte   // Register B and C
	D, E   byte   // Register D and E
	H, L   byte   // Register H and L
	SP     uint16 // Stack Pointer
	PC     uint16 // Program Counter
	Cycles int    // Cycle counter
}

// Flags
const (
	FlagZ = 0x80 // Zero flag
	FlagN = 0x40 // Negative flag
	FlagH = 0x20 // Half-carry flag
	FlagC = 0x10 // Carry flag
)

// Initialize CPU
func NewCPU() *CPU {
	return &CPU{
		A:      0,
		F:      0,
		B:      0,
		C:      0,
		D:      0,
		E:      0,
		H:      0,
		L:      0,
		SP:     0xFFFE, // Initial Stack Pointer
		PC:     0x0100, // Starting address for Game Boy
		Cycles: 0,
	}
}

// Memory interface for reading and writing
type Memory interface {
	Read(addr uint16) byte
	Write(addr uint16, value byte)
}

// Execute method for fetching and executing instructions
func (cpu *CPU) Execute(memory Memory) {
	opcode := memory.Read(cpu.PC) // Fetch the opcode
	cpu.PC++

	switch opcode {
	case 0x00: // NOP
		cpu.Cycles += 4

	// LD instructions
	case 0x01: // LD BC, d16
		cpu.B = memory.Read(cpu.PC)
		cpu.C = memory.Read(cpu.PC + 1)
		cpu.PC += 2
		cpu.Cycles += 12
	case 0x02: // LD (BC), A
		addr := (uint16(cpu.B) << 8) | uint16(cpu.C)
		memory.Write(addr, cpu.A)
		cpu.Cycles += 8
	case 0x0A: // LD A, (BC)
		cpu.A = memory.Read((uint16(cpu.B) << 8) | uint16(cpu.C))
		cpu.Cycles += 8
	case 0x3E: // LD A, d8
		cpu.A = memory.Read(cpu.PC)
		cpu.PC++
		cpu.Cycles += 8
	case 0x32: // LD (HL-), A
		addr := (uint16(cpu.H) << 8) | uint16(cpu.L)
		memory.Write(addr, cpu.A)
		cpu.L--
		cpu.Cycles += 8
	case 0x36: // LD (HL), d8
		addr := (uint16(cpu.H) << 8) | uint16(cpu.L)
		memory.Write(addr, memory.Read(cpu.PC))
		cpu.PC++
		cpu.Cycles += 12

	// ADD instructions
	case 0xC6: // ADD A, d8
		value := memory.Read(cpu.PC)
		cpu.PC++
		cpu.Add(value)
		cpu.Cycles += 8
	case 0x87: // ADD A, A
		cpu.Add(cpu.A)
		cpu.Cycles += 4
	case 0x80: // ADD A, B
		cpu.Add(cpu.B)
		cpu.Cycles += 4
	case 0x82: // ADD A, D
		cpu.Add(cpu.D)
		cpu.Cycles += 4

	// SUB instructions
	case 0xD6: // SUB d8
		value := memory.Read(cpu.PC)
		cpu.PC++
		cpu.Sub(value)
		cpu.Cycles += 8

	// INC instructions
	case 0x04: // INC B
		cpu.B++
		cpu.SetZeroFlagIfNeeded(cpu.B)
		cpu.Cycles += 4
	case 0x0C: // INC C
		cpu.C++
		cpu.SetZeroFlagIfNeeded(cpu.C)
		cpu.Cycles += 4
	case 0x14: // INC D
		cpu.D++
		cpu.SetZeroFlagIfNeeded(cpu.D)
		cpu.Cycles += 4
	case 0x1C: // INC E
		cpu.E++
		cpu.SetZeroFlagIfNeeded(cpu.E)
		cpu.Cycles += 4
	case 0x24: // INC H
		cpu.H++
		cpu.SetZeroFlagIfNeeded(cpu.H)
		cpu.Cycles += 4
	case 0x2C: // INC L
		cpu.L++
		cpu.SetZeroFlagIfNeeded(cpu.L)
		cpu.Cycles += 4
	case 0x3C: // INC A
		cpu.A++
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4

	// DEC instructions
	case 0x05: // DEC B
		cpu.B--
		cpu.SetZeroFlagIfNeeded(cpu.B)
		cpu.Cycles += 4
	case 0x0D: // DEC C
		cpu.C--
		cpu.SetZeroFlagIfNeeded(cpu.C)
		cpu.Cycles += 4
	case 0x15: // DEC D
		cpu.D--
		cpu.SetZeroFlagIfNeeded(cpu.D)
		cpu.Cycles += 4
	case 0x1D: // DEC E
		cpu.E--
		cpu.SetZeroFlagIfNeeded(cpu.E)
		cpu.Cycles += 4
	case 0x25: // DEC H
		cpu.H--
		cpu.SetZeroFlagIfNeeded(cpu.H)
		cpu.Cycles += 4
	case 0x2D: // DEC L
		cpu.L--
		cpu.SetZeroFlagIfNeeded(cpu.L)
		cpu.Cycles += 4
	case 0x3D: // DEC A
		cpu.A--
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4

	// XOR instructions
	case 0xA8: // XOR B
		cpu.A ^= cpu.B
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xA9: // XOR C
		cpu.A ^= cpu.C
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xAA: // XOR D
		cpu.A ^= cpu.D
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xAB: // XOR E
		cpu.A ^= cpu.E
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xAC: // XOR H
		cpu.A ^= cpu.H
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xAD: // XOR L
		cpu.A ^= cpu.L
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xAE: // XOR (HL)
		cpu.A ^= memory.Read((uint16(cpu.H) << 8) | uint16(cpu.L))
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 8

	// OR instructions
	case 0xB0: // OR B
		cpu.A |= cpu.B
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xB1: // OR C
		cpu.A |= cpu.C
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xB2: // OR D
		cpu.A |= cpu.D
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xB3: // OR E
		cpu.A |= cpu.E
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xB4: // OR H
		cpu.A |= cpu.H
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xB5: // OR L
		cpu.A |= cpu.L
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xB6: // OR (HL)
		cpu.A |= memory.Read((uint16(cpu.H) << 8) | uint16(cpu.L))
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 8

	// Control Flow Instructions
	case 0xC3: // JP a16
		addr := uint16(memory.Read(cpu.PC)) | (uint16(memory.Read(cpu.PC+1)) << 8)
		cpu.PC = addr
		cpu.Cycles += 16

	case 0x18: // JR r8
		offset := int8(memory.Read(cpu.PC))
		cpu.PC += uint16(offset) + 1
		cpu.Cycles += 12

	case 0xCD: // CALL a16
		addr := uint16(memory.Read(cpu.PC)) | (uint16(memory.Read(cpu.PC+1)) << 8)
		cpu.Push(cpu.PC, memory) // Push current PC to stack
		cpu.PC = addr
		cpu.Cycles += 24

	case 0xC9: // RET
		cpu.PC = cpu.Pop(memory) // Pop from stack to PC
		cpu.Cycles += 16

	default:
		fmt.Printf("Unknown opcode: %02X at PC: %04X\n", opcode, cpu.PC-1)
	}
}

// SetZeroFlagIfNeeded sets the zero flag if the value is zero
func (cpu *CPU) SetZeroFlagIfNeeded(value byte) {
	if value == 0 {
		cpu.SetZeroFlag()
	} else {
		cpu.ClearZeroFlag()
	}
}

// Helper functions to manage flags
func (cpu *CPU) SetZeroFlag() {
	cpu.F |= FlagZ
}

func (cpu *CPU) ClearZeroFlag() {
	cpu.F &^= FlagZ
}

func (cpu *CPU) SetCarryFlag() {
	cpu.F |= FlagC
}

func (cpu *CPU) ClearCarryFlag() {
	cpu.F &^= FlagC
}

// ADD operation
func (cpu *CPU) Add(value byte) {
	result := cpu.A + value
	if result == 0 {
		cpu.SetZeroFlag()
	} else {
		cpu.ClearZeroFlag()
	}
	if result < cpu.A {
		cpu.SetCarryFlag() // Carry flag if there is an overflow
	} else {
		cpu.ClearCarryFlag()
	}
	cpu.A = result
}

// SUB operation
func (cpu *CPU) Sub(value byte) {
	result := cpu.A - value
	if result == 0 {
		cpu.SetZeroFlag()
	} else {
		cpu.ClearZeroFlag()
	}
	if result > cpu.A {
		cpu.SetCarryFlag() // Set carry flag if there's a borrow
	} else {
		cpu.ClearCarryFlag()
	}
	cpu.A = result
}

// Stack operations
func (cpu *CPU) Push(value uint16, memory Memory) {
	cpu.SP -= 2
	memory.Write(cpu.SP, byte(value&0xFF))
	memory.Write(cpu.SP+1, byte(value>>8))
}

func (cpu *CPU) Pop(memory Memory) uint16 {
	value := uint16(memory.Read(cpu.SP)) | (uint16(memory.Read(cpu.SP+1)) << 8)
	cpu.SP += 2
	return value
}

// SimpleMemory implementation
type SimpleMemory struct {
	data [65536]byte // 64 KB of memory
}

func (m *SimpleMemory) Read(addr uint16) byte {
	return m.data[addr]
}

func (m *SimpleMemory) Write(addr uint16, value byte) {
	m.data[addr] = value
}

// Convert boolean to int (0 or 1)
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Main function to demonstrate CPU execution
func main() {
	mem := &SimpleMemory{}
	cpu := NewCPU()

	// Load sample instructions into memory
	mem.Write(0x0100, 0x01) // LDI BC, d16
	mem.Write(0x0101, 0x34) // Low byte
	mem.Write(0x0102, 0x12) // High byte
	mem.Write(0x0103, 0x02) // LDI (BC), A
	mem.Write(0x0104, 0x80) // ADD A, A
	mem.Write(0x0105, 0x3E) // LD A, d8
	mem.Write(0x0106, 0x0A) // Load 10 into A
	mem.Write(0x0107, 0xC6) // ADD A, d8
	mem.Write(0x0108, 0x02) // ADD A, 2
	mem.Write(0x0109, 0xC9) // RET

	// Set initial values
	cpu.A = 5        // Set Accumulator A to 5
	cpu.Execute(mem) // Execute LDI BC, d16
	cpu.Execute(mem) // Execute LDI (BC), A
	cpu.Execute(mem) // Execute ADD A, A
	cpu.Execute(mem) // Execute LD A, d8
	cpu.Execute(mem) // Execute ADD A, d8

	// Print CPU Registers and Flags
	fmt.Printf("A: %d (0x%02X)\n", cpu.A, cpu.A)
	fmt.Printf("B: %d (0x%02X)\n", cpu.B, cpu.B)
	fmt.Printf("C: %d (0x%02X)\n", cpu.C, cpu.C)
	fmt.Printf("D: %d (0x%02X)\n", cpu.D, cpu.D)
	fmt.Printf("E: %d (0x%02X)\n", cpu.E, cpu.E)
	fmt.Printf("H: %d (0x%02X)\n", cpu.H, cpu.H)
	fmt.Printf("L: %d (0x%02X)\n", cpu.L, cpu.L)
	fmt.Printf("SP: %04X\n", cpu.SP)
	fmt.Printf("PC: %04X\n", cpu.PC)
	fmt.Printf("F: %02X (Flags: Z: %d, N: %d, H: %d, C: %d)\n", cpu.F,
		btoi(cpu.F&FlagZ != 0), // Convert bool to int (0 or 1)
		btoi(cpu.F&FlagN != 0),
		btoi(cpu.F&FlagH != 0),
		btoi(cpu.F&FlagC != 0),
	)
}
