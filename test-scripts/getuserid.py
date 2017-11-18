#!/bin/python

import requests
import json
import sys

if len(sys.argv) != 2:
    print("format: ./getuserid.py username")
    sys.exit(-1)

dat = { 'Email': sys.argv[1]}
req = requests.post('https://productive.ahouts.com/get_user_id',
                   headers={'Content-Type': 'application/json'},
                   data=json.dumps(dat).encode())
print(req.text)
