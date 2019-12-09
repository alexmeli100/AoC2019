#include <iostream>
#include "intcode.h"

Mode get_mode(int m) {
    return m == 1 ? Mode::IMM : Mode::POS;
}

instruction IntCode::next_op() {
    int next = instructions[pos];
    int op = next % 100;
    next /= 100;
    auto p1 = get_mode(next % 10);
    auto p2 = get_mode(next / 10);

    return instruction(op, p1, p2);
}

void IntCode::eval(int input) {
    const auto [op, p1, p2] = next_op();

    switch (op) {
        case 99:
            exit = true;
            break;
        case 1:
            instructions[instructions[pos+3]] = get_param(p1, instructions[pos+1]) + get_param(p2, instructions[pos+2]);
            pos += 4;
            break;
        case 2:
            instructions[instructions[pos+3]] = get_param(p1, instructions[pos+1]) * get_param(p2, instructions[pos+2]);
            pos += 4;
            break;
        case 3:
            instructions[instructions[pos+1]] = input;
            pos += 2;
            break;
        case 4:
            std::cout << get_param(p1, instructions[pos+1]) << std::endl;
            pos += 2;
            break;
        case 5:
            pos = get_param(p1, instructions[pos+1]) != 0 ? get_param(p2, instructions[pos+2]) : pos+3;
            break;
        case 6:
            pos = get_param(p1, instructions[pos+1]) == 0 ? get_param(p2, instructions[pos+2]) : pos+3;
            break;
        case 7:
            instructions[instructions[pos+3]] = get_param(p1, instructions[pos+1]) < get_param(p2, instructions[pos+2]) ? 1 : 0;
            pos += 4;
            break;
        case 8:
            instructions[instructions[pos+3]] = get_param(p1, instructions[pos+1]) == get_param(p2, instructions[pos+2]) ? 1 : 0;
            pos += 4;
            break;
        default:
            exit = true;
            break;
    }
}

int IntCode::get_param(Mode m, int value) {
    return m == Mode::IMM ? value : instructions[value];
}

void IntCode::run(int input) {
    while (!exit)
        eval(input);
}

