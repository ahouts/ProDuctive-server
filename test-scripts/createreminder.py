#!/bin/python

import urllib.request
import json
import sys

if len(sys.argv) != 4:
    print("format: ./insertuser.py username password body")
    sys.exit(-1)

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2], 'Body': sys.argv[3]}
req = urllib.request.Request('https://productive.ahouts.com/reminder')
req.add_header('Content-Type', 'application/json')
req.get_method = lambda: "POST"
try:
    response = urllib.request.urlopen(req, json.dumps(dat).encode())
except urllib.request.HTTPError as e:
    error_message = e.read()
    print(error_message)
