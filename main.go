package main

import (
	cpuPkg "clockworkgnome/cpu"    // Adjust this import to match your project structure
	memPkg "clockworkgnome/memory" // Adjust this import to match your project structure
	"fmt"
)

func main() {
	fmt.Println("Starting Game Boy Emulator...")

	// Initialize the memory and CPU with sample ROM data
	ROMData := []byte{
		0x3E, // LD A, d8
		0x10, // Load 16 into A (A = 16)
		0xC6, // ADD A, d8
		0x02, // ADD A, 2 (A should become 18)
		0xC9, // RET (return)
	}

	mem := memPkg.NewMemory(ROMData) // Initialize memory with ROM data
	cpu := cpuPkg.NewCPU()           // Create a new CPU instance
	cpu.PC = 0x0000                  // Set initial program counter

	// Set initial CPU values
	cpu.A = 5 // Set Accumulator A to 5

	// Main emulation loop
	for {
		cpu.Execute(mem) // Execute the next instruction

		// Print CPU Registers and Flags after execution
		fmt.Printf("A: %d (0x%02X)\n", cpu.A, cpu.A)
		fmt.Printf("B: %d (0x%02X)\n", cpu.B, cpu.B)
		fmt.Printf("C: %d (0x%02X)\n", cpu.C, cpu.C)
		fmt.Printf("D: %d (0x%02X)\n", cpu.D, cpu.D)
		fmt.Printf("E: %d (0x%02X)\n", cpu.E, cpu.E)
		fmt.Printf("H: %d (0x%02X)\n", cpu.H, cpu.H)
		fmt.Printf("L: %d (0x%02X)\n", cpu.L, cpu.L)
		fmt.Printf("SP: %04X\n", cpu.SP)
		fmt.Printf("PC: %04X\n", cpu.PC)
		fmt.Printf("F: %02X (Flags: Z: %v, N: %v, H: %v, C: %v)\n", cpu.F,
			cpu.F&cpuPkg.FlagZ != 0, // Convert bool to int (0 or 1)
			cpu.F&cpuPkg.FlagN != 0,
			cpu.F&cpuPkg.FlagH != 0,
			cpu.F&cpuPkg.FlagC != 0,
		)

		// Simple exit condition
		if cpu.PC >= uint16(len(ROMData)) { // Check if PC exceeds ROM data size
			fmt.Println("Ending emulation loop.")
			break // Exit the emulation loop
		}
	}
}
