package cpu

import (
	"clockworkgnome/memory"
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
	IM     byte   // Interrupt Master Flag
	Timer  int    // Timer for emulation
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
		IM:     0, // Interrupt Master Flag
		Timer:  0, // Initialize Timer
	}
}

// Execute method for fetching and executing instructions
func (cpu *CPU) Execute(memory memory.Memory) {
	opcode := memory.Read(cpu.PC) // Fetch the opcode
	cpu.PC++

	switch opcode {
	// Jump Instructions
	case 0xC3: // JP a16
		addr := uint16(memory.Read(cpu.PC)) | (uint16(memory.Read(cpu.PC+1)) << 8)
		cpu.PC = addr
		cpu.Cycles += 16

	case 0xC2: // JP NZ, a16
		addr := uint16(memory.Read(cpu.PC)) | (uint16(memory.Read(cpu.PC+1)) << 8)
		if cpu.F&FlagZ == 0 { // Jump if Zero flag is clear
			cpu.PC = addr
			cpu.Cycles += 16
		} else {
			cpu.PC += 2
			cpu.Cycles += 12
		}

	case 0xDA: // JP Z, a16
		addr := uint16(memory.Read(cpu.PC)) | (uint16(memory.Read(cpu.PC+1)) << 8)
		if cpu.F&FlagZ != 0 { // Jump if Zero flag is set
			cpu.PC = addr
			cpu.Cycles += 16
		} else {
			cpu.PC += 2
			cpu.Cycles += 12
		}

	// JR Instructions
	case 0x18: // JR r8
		offset := int8(memory.Read(cpu.PC))
		cpu.PC += uint16(offset) + 1
		cpu.Cycles += 12

	case 0x20: // JR NZ, r8
		offset := int8(memory.Read(cpu.PC))
		if cpu.F&FlagZ == 0 { // Jump if Zero flag is clear
			cpu.PC += uint16(offset)
		}
		cpu.PC++
		cpu.Cycles += 12

	case 0x28: // JR Z, r8
		offset := int8(memory.Read(cpu.PC))
		if cpu.F&FlagZ != 0 { // Jump if Zero flag is set
			cpu.PC += uint16(offset)
		}
		cpu.PC++
		cpu.Cycles += 12

	// CALL Instructions
	case 0xCD: // CALL a16
		addr := uint16(memory.Read(cpu.PC)) | (uint16(memory.Read(cpu.PC+1)) << 8)
		cpu.Push(cpu.PC, memory) // Push current PC to stack
		cpu.PC = addr
		cpu.Cycles += 24

	case 0xC4: // CALL NZ, a16
		addr := uint16(memory.Read(cpu.PC)) | (uint16(memory.Read(cpu.PC+1)) << 8)
		if cpu.F&FlagZ == 0 { // Call if Zero flag is clear
			cpu.Push(cpu.PC, memory)
			cpu.PC = addr
			cpu.Cycles += 24
		} else {
			cpu.PC += 2
			cpu.Cycles += 12
		}

	case 0xCC: // CALL Z, a16
		addr := uint16(memory.Read(cpu.PC)) | (uint16(memory.Read(cpu.PC+1)) << 8)
		if cpu.F&FlagZ != 0 { // Call if Zero flag is set
			cpu.Push(cpu.PC, memory)
			cpu.PC = addr
			cpu.Cycles += 24
		} else {
			cpu.PC += 2
			cpu.Cycles += 12
		}
	case 0x3E: // LD A, d8
		cpu.A = memory.Read(cpu.PC) // Load immediate value into register A
		cpu.PC++
		cpu.Cycles += 8 // 8 cycles for LD A, d8

	case 0xC6: // ADD A, d8
		d8 := memory.Read(cpu.PC) // Get the immediate value
		cpu.Add(d8)               // Add to A
		cpu.PC++
		cpu.Cycles += 8 // 8 cycles for ADD A, d8

	// RET Instructions
	case 0xC9: // RET
		cpu.PC = cpu.Pop(memory) // Pop from stack to PC
		cpu.Cycles += 16

	case 0xD9: // RETI
		cpu.PC = cpu.Pop(memory) // Pop from stack to PC
		cpu.Cycles += 16
		// Handle additional logic required for Return from Interrupt here if needed

	case 0xC0: // RET NZ
		if cpu.F&FlagZ == 0 { // Return if Zero flag is clear
			cpu.PC = cpu.Pop(memory)
			cpu.Cycles += 16
		} else {
			cpu.Cycles += 8 // If not returning, just consume cycles
		}

	case 0xC8: // RET Z
		if cpu.F&FlagZ != 0 { // Return if Zero flag is set
			cpu.PC = cpu.Pop(memory)
			cpu.Cycles += 16
		} else {
			cpu.Cycles += 8 // If not returning, just consume cycles
		}

	case 0xD0: // RET NC
		if cpu.F&FlagC == 0 { // Return if Carry flag is clear
			cpu.PC = cpu.Pop(memory)
			cpu.Cycles += 16
		} else {
			cpu.Cycles += 8 // If not returning, just consume cycles
		}

	case 0xD8: // RET C
		if cpu.F&FlagC != 0 { // Return if Carry flag is set
			cpu.PC = cpu.Pop(memory)
			cpu.Cycles += 16
		} else {
			cpu.Cycles += 8 // If not returning, just consume cycles
		}

	// Logical AND Instructions
	case 0xA4: // AND B
		cpu.A &= cpu.B
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xA5: // AND C
		cpu.A &= cpu.C
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 4
	case 0xA6: // AND (HL)
		cpu.A &= memory.Read((uint16(cpu.H) << 8) | uint16(cpu.L))
		cpu.ClearCarryFlag()
		cpu.SetZeroFlagIfNeeded(cpu.A)
		cpu.Cycles += 8

	// Logical OR Instructions
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

	// BIT instructions (bit manipulation)
	case 0xCB: // Example prefix for BIT operation
		switch memory.Read(cpu.PC) {
		case 0x40: // BIT 0, B
			cpu.SetZeroFlagIfNeeded(cpu.B & 0x01)
			cpu.Cycles += 8
			cpu.PC++
		case 0x41: // BIT 0, C
			cpu.SetZeroFlagIfNeeded(cpu.C & 0x01)
			cpu.Cycles += 8
			cpu.PC++
		case 0x42: // BIT 0, D
			cpu.SetZeroFlagIfNeeded(cpu.D & 0x01)
			cpu.Cycles += 8
			cpu.PC++
		// Add more BIT cases for each register...

		default:
			fmt.Printf("Unhandled BIT operation\n")
		}

	// Placeholder for timer handling (time-based operations)
	// Timer management can be expanded later

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
	result := uint16(cpu.A) + uint16(value)
	if result > 0xFF {
		cpu.SetCarryFlag() // Set carry flag if there's an overflow
	} else {
		cpu.ClearCarryFlag()
	}
	cpu.A = byte(result) // Store the lower 8 bits
	cpu.SetZeroFlagIfNeeded(cpu.A)
}

// SUB operation
func (cpu *CPU) Sub(value byte) {
	result := uint16(cpu.A) - uint16(value)
	if result == 0 {
		cpu.SetZeroFlag()
	} else {
		cpu.ClearZeroFlag()
	}
	if result > 0xFF {
		cpu.SetCarryFlag() // Set carry flag if there's a borrow
	} else {
		cpu.ClearCarryFlag()
	}
	cpu.A = byte(result) // Store the lower 8 bits
}

// Stack operations
func (cpu *CPU) Push(value uint16, memory memory.Memory) {
	cpu.SP -= 2
	memory.Write(cpu.SP, byte(value&0xFF))
	memory.Write(cpu.SP+1, byte(value>>8))
}

func (cpu *CPU) Pop(memory memory.Memory) uint16 {
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
	ROMData := []byte{0x01, 0x34, 0x12, 0x02, 0x80, 0x3E, 0x0A, 0xC6, 0x02, 0xC9} // Sample ROM data
	mem := memory.NewMemory(ROMData)                                              // Initialize memory with ROM data
	cpu := NewCPU()

	// Initialize Accumulator A
	cpu.A = 5 // Set Accumulator A to 5

	// Execute instructions
	for cpu.PC < uint16(len(ROMData)) {
		cpu.Execute(mem) // Execute instructions in memory
	}

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
