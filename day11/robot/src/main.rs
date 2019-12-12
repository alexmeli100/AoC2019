use std::collections::{HashMap};
use num_complex::Complex;
use std::net::TcpStream;
use std::io::{Write, BufReader, BufRead};
use std::error::Error;
use std::fmt::{Display, Formatter};
use core::fmt;

type Panel = HashMap<Complex<i32>, Color>;

enum Color {
    Black,
    White
}

struct Robot {
    dir: Complex<i32>,
    pos: Complex<i32>,
    panels: Panel
}

impl Robot {
    fn new() -> Robot {
        let mut panels = HashMap::new();
        panels.insert(Complex::new(0, 0), Color::White);

        Robot {
            dir: Complex::new(0, 1),
            pos: Complex::new(0, 0),
            panels
        }
    }

    fn draw(&mut self, stream: &mut TcpStream) -> Result<(), Box<dyn Error>> {
        let mut reader = BufReader::new(stream.try_clone()?);

        loop {
            let mut input = 0;

            if let Some(color) = self.panels.get(&self.pos) {
                match color {
                    Color::Black => input = 0,
                    Color::White => input = 1
                }
            }

            stream.write(format!("{}\n", input).as_bytes())?;

            let mut out = String::new();
            reader.read_line(&mut out)?;

            let command = out.trim();
            if command == "DONE" {
                break;
            }

            let ins = command.split(' ').map(|n| n.parse::<usize>().unwrap()).collect::<Vec<usize>>();
            self.update(ins);
        }

        Ok(())
    }

    fn update(&mut self, ins: Vec<usize>) {
        self.paint(ins[0]);
        self.walk(ins[1]);
    }

    fn paint(&mut self, c: usize) {
        match c {
            0 => self.panels.insert(self.pos, Color::Black),
            1 => self.panels.insert(self.pos, Color::White),
            _ => panic!("Unknown color")
        };
    }

    fn walk(&mut self, c: usize) {
        match c {
            0 => self.dir *= Complex::new(0, 1),
            1 => self.dir *= Complex::new(0, -1),
            _ => panic!("Unknow direction")
        };

        self.pos += self.dir;
    }

    fn bounds_x(&self) -> (i32, i32) {
        self.panels.keys().fold((i32::min_value(), i32::max_value()), |(mx_x, mn_x), x| {
            (std::cmp::max(mx_x, x.re), std::cmp::min(mn_x, x.re))
        })
    }

    fn bounds_y(&self) -> (i32, i32) {
        self.panels.keys().fold((i32::min_value(), i32::max_value()), |(mx_y, mn_y), x| {
            (std::cmp::max(mx_y, x.im), std::cmp::min(mn_y, x.im))
        })
    }
}

impl Display for Robot {
    fn fmt(&self, f: &mut Formatter) -> fmt::Result {
        let (mx_x, mn_x) = self.bounds_x();
        let (mx_y, mn_y) = self.bounds_y();

        for y in mn_y..mx_y +1 {
            for x in mn_x..mx_x +1 {
                let color = self.panels.get(&Complex::new(x, y)).unwrap_or(&Color::Black);

                match color {
                    Color::White => "▋".fmt(f)?,
                    Color::Black => " ".fmt(f)?
                }
            }
            writeln!(f, "")?;
        }

        Ok(())
    }
}

fn main() {
    let mut stream = TcpStream::connect("localhost:8080").expect("Error creating stream");
    let mut r = Robot::new();

    match r.draw(&mut stream) {
        Ok(()) => (),
        Err(e) => panic!(format!("{:?}", e))
    };

    let painted = r.panels.len();

    println!("{}", painted);
    print!("{}", r);
}
