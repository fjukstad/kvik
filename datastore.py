import numpy as np
import json
import ast
import zmq
import time
import datetime


class RPC():

    def add(self, a, b):
        return a+b

    def sum(self,nums):
        res = 0
        for num in nums:
            res += num

        return res

    def std(self,a):
        return np.std(a)

    def var(self,a):
        return np.var(a)

    def mean(self,a):
        return np.mean(a)


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

    rpc = RPC()

    while True:
        #  Wait for next request from client
        message = socket.recv()

        time = datetime.datetime.now().strftime("%H:%M:%S.%f")

        # print time, "Incoming request"

        req = ast.literal_eval(message)

        met = req["Method"]
        args = req["Args"]

        #result = getattr(globals, ()[met](args)
        status,result = call(rpc,met,args)

        if result == 0:
            resp = json.dumps({"Response":result,"Status":status})
        else:
            resp = json.dumps({"Response":result,"Status":status})
        #  Send reply back to client
        socket.send(resp)

        time = datetime.datetime.now().strftime("%H:%M:%S.%f")

        print time, "Completed RPC call"


