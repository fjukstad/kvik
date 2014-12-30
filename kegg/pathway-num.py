import csv
import numpy as np

nums = []
with open('p-n.csv', 'rb') as csvfile:
    reader = csv.reader(csvfile, delimiter=',', quotechar='|', )
    next(reader) # skip first line (header)
    for row in reader:
        nums.append(int(row[1]))
    print np.mean(nums)
    print np.std(nums)
    print np.var(nums)
    print len(nums)
