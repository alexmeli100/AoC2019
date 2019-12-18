use std::iter;
use std::fs;
use std::collections::HashMap;

const BASE: &[i64] = &[0, 1, 0, -1];

fn get_result(seq: Vec<i64>) -> Vec<i64>  {
    let mut cache: HashMap<usize, Vec<i64>> = HashMap::new();

    iter::successors(Some(seq), |s| next_seq(s, &mut cache))
        .skip(100)
        .take(1)
        .next().unwrap()
}

fn next_seq(seq: &Vec<i64>, cache: &mut HashMap<usize, Vec<i64>>) -> Option<Vec<i64>> {
    let next = (0..seq.len())
        .map(|i| {
            let pat = cache.entry(i+1).or_insert(get_patten(BASE, i+1));
            let res: i64 = seq.iter()
                .zip(pat.iter().cycle().skip(1))
                .map(|(x, y)| x*y).sum();

            (res % 10).abs()
        }).collect();

    Some(next)
}

fn get_patten(pat: &[i64], n: usize) -> Vec<i64> {
    pat.iter()
        .flat_map(|x| iter::repeat(*x).take(n))
        .collect()
}

fn get_string(n: &[i64]) -> String {
    n.iter().map(|i| i.to_string()).collect::<String>()
}

fn part2(seq: Vec<i64>) -> String {
    let index = seq.iter().take(7).fold(0, |acc, val| val + acc * 10) as usize;
    let s = seq.iter().cycle().take(seq.len() * 10000).map(|x| *x).collect::<Vec<i64>>();
    let res = get_result(s);

    get_string(&res[index..index+8])
}

fn part2_fast(seq: Vec<i64>) -> String {
    let index = seq.iter().take(7).fold(0, |acc, val| val + acc * 10) as usize;
    let mut s = seq.iter().cycle().take(seq.len() * 10000).skip(index).map(|x| *x).collect::<Vec<i64>>();

    for _ in 0..100 {
        for i in (0..s.len()-1).rev() {
            s[i] = (s[i] + s[i+1]) % 10;
        }
    }


    get_string(&s[0..8])
}

fn part1(seq: Vec<i64>) -> String {
    let res = get_result(seq);

    get_string(&res[0..8])
}

fn main() {

    let input = read_input().unwrap();
    println!("{:?}", part2_fast(input));
}

fn read_input() -> Result<Vec<i64>, std::io::Error> {
    let input = fs::read_to_string("input.txt")?;

    let res = input.trim().chars().map(|c| c.to_digit(10).unwrap() as i64).collect();

    Ok(res)
}
