#!/bin/python

import requests
import json
import sys

if len(sys.argv) != 2:
    print("format: ./getuserid.py username")
    sys.exit(-1)

dat = { 'Email': sys.argv[1]}
req = requests.put('https://productive.ahouts.com/user/getid',
                   headers={'Content-Type': 'application/json'},
                   data=json.dumps(dat).encode())
print(req.text)
