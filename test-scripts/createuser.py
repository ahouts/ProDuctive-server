#!/bin/python

import requests
import json
import sys

if len(sys.argv) != 3:
    print("format: ./insertuser.py username password")
    sys.exit(-1)

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2]}
req = requests.post('https://productive.ahouts.com/user',
                    headers={'Content-Type': 'application/json'},
                    data=json.dumps(dat).encode())
print(req.text)
