import math
import os
import csv
import sys


# parse results from go benchmark and store in new file
def storeResults(tmp,filename):
    with open(tmp) as tsv:
        with open(filename, "wb") as output:
            writer = csv.writer(output, delimiter=',')
            for line in csv.reader(tsv, delimiter="\t"):
                if line[0].startswith("Benchmark"):
                    print line
                    writer.writerow(line)



if __name__=="__main__":
    numIter = 200
    for i in range(0,numIter):
        os.system("go test -bench=BenchmarkDatasetSize -benchmem > tmp.tsv")
        storeResults("tmp.tsv", "result.csv")
        sys.stdout.write(str(i)+" of "+str(numIter)+" done\r")
    print ""

