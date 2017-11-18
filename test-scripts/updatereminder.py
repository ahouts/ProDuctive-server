#!/bin/python2

import requests
import json
import sys

if len(sys.argv) != 5:
    print("format: ./insertuser.py username password id body")
    sys.exit(-1)

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2], 'ReminderId': int(sys.argv[3]), 'Body': sys.argv[4]}
req = requests.put('https://productive.ahouts.com/reminder',
                   headers={'Content-Type': 'application/json'},
                   data=json.dumps(dat).encode())
print(req.text)
