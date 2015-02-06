import zmq
import json
from random import randint


def randomArray(length=10, min=0, max=100):
    res = []
    for i in range(length):
        res.append(randint(min,max))
    return res


if __name__ == "__main__":
    context = zmq.Context()
    socket = context.socket(zmq.REQ)

    while True:
        socket.connect("tcp://localhost:5555")

        arr = randomArray()
        a = str(arr)
        b = a.replace(",","\,")
        #socket.send("{\"Method\":\"std\", \"Args\":"+b+"}")
        socket.send("{\"Method\":\"add\", \"Args\":[1,2]}")

        resp = socket.recv()
        print resp
        socket.close()
