import csv
import math
import numpy as np
from scipy import stats
from hurry.filesize import size
import matplotlib.pyplot as plt


blu = "#005CAA"
ora = "#CF5112"
gre = "#009060"
pin = "#C25E99"
yel = "#E19300"

def plot(d,i,name="figure", color=gre):

    ax = plt.subplot(2,3,i)
    #ax = plt.figure()
    ax.spines["right"].set_visible(False)
    ax.spines["top"].set_visible(False)
    ax.tick_params(axis='both', direction='out')
    ax.get_xaxis().tick_bottom()   # remove unneeded ticks
    ax.get_yaxis().tick_left()

    binwidth = 3
    plt.hist(d, bins=len(d)/3,color=color,histtype='stepfilled')
    plt.axis([min(d)-int(min(d)/2.2), max(d)+1, 0, 100])
    plt.title(name)
    plt.xlabel('Runtime(s)')
    plt.ylabel('Frequency')
    #a = pl.frange(0.5,int(max(d))+1,1)
    #plt.xticks(a)
    #plt.xticks(range(0,int(max(d)+1)))
    plt.yticks(range(0,40,12))



if __name__ == "__main__":
    sizes = ["-1x","-2x","-5x", "-10x", "-20x"]

    results = {}
    for s in sizes:
        filename = "result"+s+".csv"
        print "-----------"+s+"------------"
        with open(filename) as file:
            for line in csv.reader(file, delimiter=","):
                method = line[0]
                count = int(line[1].lstrip(" "))
                a = line[2].lstrip(" ")
                a = a.split(" ")[0]
                runtime = float(a)


                a = line[3].lstrip(" ")
                a = a.split(" ")[0]
                mem = int(a)

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

        d1 = []
        d2 = []

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
            gmean = stats.gmean(runtimes)
            print method, "Runtime (mean,std):", mean, std/ns
            #plot(runtimes,i,method)
            d1.append(mean)

            mean = np.mean(mems)
            std = np.std(mems)
            var = np.var(mems)
            sem = stats.sem(mems)
            print method,"Memory usage (mean, std):", mean/mb, std/mb

            d2.append(mean)

    #plt.subplots_adjust(hspace=.5,wspace=.75)
    #plt.show()


    #d1.sort()
    #d2.sort()
    #print d1
    #print d2



