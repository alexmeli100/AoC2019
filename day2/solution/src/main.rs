use std::fs;

#[derive(Debug, Clone)]
struct IntCode {
    ints: Vec<usize>,
    pos: usize
}

#[derive(Debug)]
enum OpCode {
    Add,
    Mul,
    Halt
}

impl From<usize> for OpCode {
    fn from(x: usize) -> Self {
        match x {
            1  => OpCode::Add,
            2  => OpCode::Mul,
            99 => OpCode::Halt,
            _  => panic!("Invalid op code")
        }
    }
}

impl IntCode {
    fn next_op(&mut self) -> OpCode {
        OpCode::from(self.ints[self.pos])
    }

    fn run(&mut self) {
        loop {
            let op = self.next_op();

            match op {
                OpCode::Add => {
                    let out = self.ints[self.pos+3];
                    self.ints[out] = self.ints[self.ints[self.pos+1]] + self.ints[self.ints[self.pos+2]];
                },
                OpCode::Mul => {
                    let out = self.ints[self.pos+3];
                    self.ints[out] = self.ints[self.ints[self.pos+1]] * self.ints[self.ints[self.pos+2]];
                }
                OpCode::Halt => break
            }

            self.pos += 4;
        }
    }
}

fn part1(int_code: &IntCode, noun: usize, verb: usize) -> usize {
    let mut next = int_code.clone();
    next.ints[1] = noun;
    next.ints[2] = verb;

    next.run();
    next.ints[0]
}

fn part2(int_code: &IntCode) -> usize {
    for noun in 0..99 {
        for verb in 0..99 {
            if part1(int_code, noun, verb) == 19690720 {
                return noun * 100 + verb
            }
        }
    }

    0
}

fn main() {
    let input = read_input("input.txt").expect("Error opening file");
    let int_code = IntCode{ ints: input, pos: 0};

    println!("Part 2: {}", part2(&int_code));
}

fn read_input(path: &str) -> Result<Vec<usize>, std::io::Error> {
    let input = fs::read_to_string(path)?;

    let res = input.split(',')
        .map(|x| x.parse::<usize>().unwrap())
        .collect();

    Ok(res)
}