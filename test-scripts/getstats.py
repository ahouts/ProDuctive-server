#!/bin/python2

import requests
import json
import sys

if len(sys.argv) != 1:
    print("format: ./getuser.py")
    sys.exit(-1)

req = requests.get('https://productive.ahouts.com/stats', 
                   headers={'Content-Type': 'application/json'})
print(req.text)
