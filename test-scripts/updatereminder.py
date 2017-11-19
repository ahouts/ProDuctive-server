#!/bin/python2

import requests
import json
import sys

if len(sys.argv) != 5:
    print("format: ./updatereminder.py username password id body")
    sys.exit(-1)

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2], 'Body': sys.argv[4]}
req = requests.post('https://productive.ahouts.com/reminder/' + sys.argv[3],
                   headers={'Content-Type': 'application/json'},
                   data=json.dumps(dat).encode())
print(req.text)
