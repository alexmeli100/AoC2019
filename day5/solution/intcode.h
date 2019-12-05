#include <utility>
#include <vector>
#include <tuple>

#ifndef SOLUTION_INTCODE_H
#define SOLUTION_INTCODE_H

enum class Mode { POS, IMM };

typedef std::tuple<int, Mode, Mode> instruction;


class IntCode {
public:
    IntCode() = default;
    explicit IntCode(std::vector<int> ins): instructions(std::move(ins)), pos(0), exit(false){};
    void run(int input);

private:
    void eval(int input);
    instruction next_op();
    int get_param(Mode m, int value);

    std::vector<int> instructions;
    unsigned int pos{};
    bool exit{};
    int previous{3};
};

#endif //SOLUTION_INTCODE_H
