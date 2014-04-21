import csv
import math
import numpy as np
from scipy import stats
from hurry.filesize import size

results = {}

with open("result.csv") as file:
    for line in csv.reader(file, delimiter=","):
        method = line[0]
        count = int(line[1].lstrip(" "))
        runtime = float(line[2].split(" ")[0])
        mem = int(line[3].split(" ")[0])

        a = line[4].lstrip(" ")
        alloc = int(a.split(" ")[0])

        res = [{"count":count, "runtime":runtime,"memory":mem,"alloc":alloc}]
        try:
            a = results[method]
            results[method] = res+a
        except KeyError:
            results[method] = res

ns = 1000000000
mb = 1048576

for i, method in  enumerate(results):
    res = results[method]

    runtimes = []
    mems = []
    count = 0

    for j in range(len(res)):
        r = res[j]
        runtimes.append(r["runtime"])
        mems.append(r["memory"])
        count = count + r["count"]

    mean = np.mean(runtimes)/ns #seconds
    std = np.std(runtimes)
    var = np.var(runtimes)
    sem = stats.sem(runtimes)/ns #seconds
    print method, "Runtime:", mean, sem, std/ns

    mean = np.mean(mems)
    std = np.std(mems)
    var = np.var(mems)
    sem = stats.sem(mems)
    print method,"Memory usage:", mean/mb, sem/mb, std/mb






