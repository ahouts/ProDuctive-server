#!/bin/python2

import requests
import json
import sys

if len(sys.argv) != 4:
    print("format: ./getreminder.py username password id")
    sys.exit(-1)

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2]}
req = requests.put('https://productive.ahouts.com/reminder/' + sys.argv[3], 
                   headers={'Content-Type': 'application/json'}, 
                   data=json.dumps(dat).encode())
print(req.text)
