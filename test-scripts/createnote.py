#!/bin/python2

import requests
import json
import sys

if len(sys.argv) != 5:
    print("format: ./createnote.py username password title body")
    sys.exit(-1)

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2], 'Title': sys.argv[3], 'Body': sys.argv[4], 'ProjectId': None}
req = requests.post('https://productive.ahouts.com/note',
                   headers={'Content-Type': 'application/json'},
                   data=json.dumps(dat).encode())
print(req.text)
