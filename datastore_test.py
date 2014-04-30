import math
import os
import csv
import sys


# parse results from go benchmark and store in new file
def storeResults(tmp,filename):
    with open(tmp) as tsv:
        with open(filename, "a") as output:
            writer = csv.writer(output, delimiter=',')
            for line in csv.reader(tsv, delimiter="\t"):
                if line[0].startswith("Benchmark"):
                    writer.writerow(line)



if __name__=="__main__":
    numIter = 200
    sizes = ["","-2x","-5x","-10x"]
    for i in range(0,numIter):
        for size in sizes:
            tmp = "tmp"+size+".tmp"
            outfile = "result"+size+".csv"

            if size == "":
                outfile = "result-1x.csv"

            os.system("go test -bench=. -benchmem -cpuprofile=\""+size+"\" -timeout=2000m > "+tmp)
            storeResults(tmp, outfile)
            #sys.stdout.write(str(i)+" of "+str(numIter)+" done\r")
        print i,"of",numIter,"done"
    print ""

