import timeit
import csv

def local():
    with open("data/exprs.csv") as file:
        for line in csv.reader(file, delimiter=","):
            a = line


def remote():
    with open("/Users/bjorn/stallo/src/src/nowac/datastore/data-2x/exprs.csv") as file:
        for line in csv.reader(file, delimiter=","):
            a = line

if __name__=="__main__":
    a = timeit.timeit("local()", setup="from __main__ import local", number=1)
    b = timeit.timeit("remote()", setup="from __main__ import remote",number=1)

    print "local:",a
    print "remote:",b
