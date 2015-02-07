import numpy as np
import json
import sys
import ast
import zmq
import time
import datetime
import rpy2
from rpy2.robjects.packages import SignatureTranslatedAnonymousPackage
import rpy2.robjects as robjects


class Worker():

    def __init__(self,r_filename):

        string = ""
        with open (r_filename, "r") as myfile:
            string =''.join(myfile.readlines())
        r = SignatureTranslatedAnonymousPackage(string, "r")

        robjects.r(string)

        self.r = r

    def command(self, string):
        print "Got the command "+string
        try:
            ret = robjects.r(string)
        except:
            return "Could not evaluate "+string
        print ret
        print str(ret)
        if len(ret) > 1:
            print "Result:", str(ret)
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

    port = ""
    if len(sys.argv) < 2:
        port = "5555"
        script = "script.r"
    else:
        port = sys.argv[1]
        script = sys.argv[2]

    socket.bind("tcp://*:"+port)

    worker = Worker(script)

    time = datetime.datetime.now().strftime("%H:%M:%S.%f")
    print time, "Worker accepting new connections..."

    while True:
        message = socket.recv()

        time = datetime.datetime.now().strftime("%H:%M:%S.%f")
        print time, "Incoming RPC"

        req = ast.literal_eval(message)

        met = req["Method"]
        args = req["Args"]

        status,result = call(worker,met,args)

        if result == 0:
            resp = json.dumps({"Response":result,"Status":status})
        else:
            resp = json.dumps({"Response":result,"Status":status})
        #  Send reply back to client
        socket.send(resp)
`
        time = datetime.datetime.now().strftime("%H:%M:%S.%f")

        print time, "Completed RPC call"
        print "sent", result



