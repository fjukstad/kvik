import numpy as np
import json
import ast
import zmq
import time
import datetime
import rpy2
from rpy2.robjects.packages import SignatureTranslatedAnonymousPackage
import rpy2.robjects as robjects


class DataEngine():

    def __init__(self,r_filename):

        # Read in r file and init data engine
        string = ""
        with open (r_filename, "r") as myfile:
            string =''.join(myfile.readlines())
        r = SignatureTranslatedAnonymousPackage(string, "r")

        self.r = r

    def add(self, a, b):
        return r.add(a,b)[0]

    def sum(self,nums):
        v = robjects.FloatVector(nums)
        rsum = robjects.r['sum']
        return rsum(v)[0]

    def std(self,a):
        v = robjects.FloatVector(a)
        rsd = robjects.r['sd']
        return rsd(v)[0]

    def var(self,a):
        v = robjects.FloatVector(a)
        rvar = robjects.r['var']
        return rvar(v)[0]

    def mean(self,a):
        v = robjects.FloatVector(a)
        rmean = robjects.r['mean']
        return rmean(v)[0]

    def command(self, string):
        print "Got the command "+string
        try:
            ret = robjects.r(string)
        except:
            return "Could not evaluate "+string
        print ret
        print len(ret)
        print str(ret)
        if len(ret) > 1:
            #return ' '.join(str(x) for x in ret)
            return str(ret)
        try:
            return ret[0]
        except IndexError:
            return "Could not evaluate "+string

def call(obj, func, attr):
    if hasattr(obj,func):
        try:
            return 0, getattr(obj,func)(*attr)
        except TypeError as err:
            print "Type error:",err
            return -1,0
    else:
        print "ERROR: Method",func,"not found"
        return -1, 0

if __name__ == "__main__":
    context = zmq.Context()
    socket = context.socket(zmq.REP)
    socket.bind("tcp://*:5555")
    #rc = socket.bind("ipc:///tmp/datastore/0")


    dataengine = DataEngine("data-engine.r")

    time = datetime.datetime.now().strftime("%H:%M:%S.%f")
    print time, "Data Engine accepting new connections..."

    while True:
        message = socket.recv()

        time = datetime.datetime.now().strftime("%H:%M:%S.%f")
        print time, "Incoming RPC"

        req = ast.literal_eval(message)

        met = req["Method"]
        args = req["Args"]

        status,result = call(dataengine,met,args)

        if result == 0:
            resp = json.dumps({"Response":result,"Status":status})
        else:
            resp = json.dumps({"Response":result,"Status":status})
        #  Send reply back to client
        socket.send(resp)

        time = datetime.datetime.now().strftime("%H:%M:%S.%f")

        print time, "Completed RPC call"
        print "sent", result



