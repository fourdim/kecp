#!/usr/bin/env python3

import multiprocessing
from os import execvp, remove
from pathlib import Path

MOD = Path('github.com/fourdim/kecp')

def basic_test(build_target :str, time_out :str):
    execvp('go', ('go', 'test', '-timeout', time_out, '-race', '-cover', '-covermode=atomic', build_target))

def advanced_race_test(build_target :str, time_out :str):
    execvp('go', ('go', 'test', '-timeout', time_out, '-cpu=1,9,55,99', '-race', '-count=100', '-failfast', '-cover', '-covermode=atomic', build_target))

test_cases = [
    {
        'target': 'modules/kecp-channel',
        'test_method': advanced_race_test,
        'time_out': '30s'
    },
    {
        'target': 'modules/kecp-crypto',
        'test_method': basic_test,
        'time_out': '30s'
    },
    {
        'target': 'modules/kecp-msg',
        'test_method': basic_test,
        'time_out': '30s'
    },
    {
        'target': 'modules/kecp-signal',
        'test_method': basic_test,
        'time_out': '180s'
    },
    {
        'target': 'modules/kecp-validate',
        'test_method': basic_test,
        'time_out': '30s'
    }
]

def test():
    ps = []
    for test_case in test_cases:
        path = Path.joinpath(MOD, test_case['target'])
        p = multiprocessing.Process(target=test_case['test_method'], args=(path.as_posix(), test_case['time_out']))
        p.start()
        ps.append(p)
    for p in ps:
        p.join()

if __name__ == '__main__':
    test()
