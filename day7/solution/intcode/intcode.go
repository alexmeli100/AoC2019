package intcode

import "log"

type IntCode struct {
	tape []int
	pc   int
	In   chan int
	Out  chan int
	base int
	Last int
	halt bool
}

type mode int
type opFunc func(i *IntCode, mode []mode)

var opFuncs = map[int]opFunc{
	1:  add,
	2:  mul,
	3:  read,
	4:  write,
	5:  jumpTrue,
	6:  jumpFalse,
	7:  lessThan,
	8:  equalTo,
	9:  changeBase,
	99: halt,
}

const (
	pos    mode = 0
	imm    mode = 1
	rel    mode = 2
	Memory int  = 8000
)

func NewVm(tape []int) *IntCode {
	t := make([]int, Memory)
	copy(t, tape)

	return &IntCode{
		tape: t,
		pc:   0,
		base: 0,
		In:   make(chan int, 2),
		Out:  make(chan int, 2),
		halt: false,
	}
}

func (vm *IntCode) getValue(m mode, value int) int {
	if m == imm {
		return value
	} else if m == rel {
		return vm.tape[vm.base+value]
	}

	return vm.tape[value]
}

func (vm *IntCode) nextOp() (int, []mode) {
	next := vm.tape[vm.pc]
	op := next % 100
	next /= 100
	p1 := next % 10
	p2 := next / 10

	return op, []mode{mode(p1), mode(p2)}
}

func (vm *IntCode) eval() {
	op, modes := vm.nextOp()

	if f, ok := opFuncs[op]; ok {
		f(vm, modes)
	} else {
		opError(op)
		vm.halt = true
	}
}

func (vm *IntCode) Run() {
	for !vm.halt {
		vm.eval()
	}
	close(vm.Out)
}

// OPCODES

func add(vm *IntCode, m []mode) {
	vm.tape[vm.tape[vm.pc+3]] = vm.getValue(m[0], vm.tape[vm.pc+1]) + vm.getValue(m[1], vm.tape[vm.pc+2])
	vm.pc += 4
}

func mul(vm *IntCode, m []mode) {
	vm.tape[vm.tape[vm.pc+3]] = vm.getValue(m[0], vm.tape[vm.pc+1]) * vm.getValue(m[1], vm.tape[vm.pc+2])
	vm.pc += 4
}

func read(vm *IntCode, m []mode) {
	//vm.tape[vm.getValue(m[0], vm.pc+1)] = <-vm.In
	//vm.pc += 2
	input := <-vm.In
	switch m[0] {
	case pos:
		vm.tape[vm.tape[vm.pc+1]] = input
	case rel:
		log.Println(vm.tape[vm.pc+1])
		log.Println(vm.base)
		vm.tape[vm.base+vm.tape[vm.pc+1]] = input
	default:
		log.Fatal("Invalid param mode ", m[0])
	}
	vm.pc += 2
}

func write(vm *IntCode, m []mode) {
	vm.Last = vm.getValue(m[0], vm.tape[vm.pc+1])
	vm.Out <- vm.Last
	vm.pc += 2
}

func jumpTrue(vm *IntCode, m []mode) {
	if vm.getValue(m[0], vm.tape[vm.pc+1]) != 0 {
		vm.pc = vm.getValue(m[1], vm.tape[vm.pc+2])
	} else {
		vm.pc += 3
	}
}

func jumpFalse(vm *IntCode, m []mode) {
	if vm.getValue(m[0], vm.tape[vm.pc+1]) == 0 {
		vm.pc = vm.getValue(m[1], vm.tape[vm.pc+2])
	} else {
		vm.pc += 3
	}
}

func lessThan(vm *IntCode, m []mode) {
	p1 := vm.getValue(m[0], vm.tape[vm.pc+1])
	p2 := vm.getValue(m[1], vm.tape[vm.pc+2])

	if p1 < p2 {
		vm.tape[vm.tape[vm.pc+3]] = 1
	} else {
		vm.tape[vm.tape[vm.pc+3]] = 0
	}

	vm.pc += 4
}

func equalTo(vm *IntCode, m []mode) {
	p1 := vm.getValue(m[0], vm.tape[vm.pc+1])
	p2 := vm.getValue(m[1], vm.tape[vm.pc+2])

	if p1 == p2 {
		vm.tape[vm.tape[vm.pc+3]] = 1
	} else {
		vm.tape[vm.tape[vm.pc+3]] = 0
	}

	vm.pc += 4
}

func changeBase(vm *IntCode, m []mode) {
	vm.base += vm.getValue(m[0], vm.tape[vm.pc+1])
	vm.pc += 2
}

func halt(vm *IntCode, m []mode) {
	vm.halt = true
}

func opError(op int) {
	log.Println("Unknown op ", op)
}
