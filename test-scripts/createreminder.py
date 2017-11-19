#!/bin/python2

import requests
import json
import sys

if len(sys.argv) != 4:
    print("format: ./createreminder.py username password body")
    sys.exit(-1)

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2], 'Body': sys.argv[3]}
req = requests.post('https://productive.ahouts.com/reminder',
                   headers={'Content-Type': 'application/json'},
                   data=json.dumps(dat).encode())
print(req.text)
