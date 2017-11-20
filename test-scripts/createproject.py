#!/bin/python2

import requests
import json
import sys

if len(sys.argv) != 4:
    print("format: ./createproject.py username password title")
    sys.exit(-1)

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2], 'Title': sys.argv[3]}
req = requests.post('https://productive.ahouts.com/project',
                   headers={'Content-Type': 'application/json'},
                   data=json.dumps(dat).encode())
print(req.text)
