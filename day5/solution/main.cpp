#include <iostream>
#include <fstream>
#include <streambuf>
#include "intcode.h"

std::vector<int> parseFile() {
    std::ifstream t("../input.txt");
    std::vector<int> program;

    while(t) {
        int i;
        char c;
        t >> i >> c;

        program.push_back(i);
    }

    return program;
}

int main() {
    auto p = parseFile();
    int input = 5;

    IntCode tape(p);
    tape.run(input);

    return 0;
}

