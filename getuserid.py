#!/bin/python

import urllib.request
import json
import sys

if len(sys.argv) != 2:
    print("format: ./getuserid.py username")
    sys.exit(-1)

dat = { 'Email': sys.argv[1]}
req = urllib.request.Request('https://productive.ahouts.com/get_user_id')
req.add_header('Content-Type', 'application/json')
try:
    response = urllib.request.urlopen(req, json.dumps(dat).encode())
    print(response.read())
except urllib.request.HTTPError as e:
    error_message = e.read()
    print(error_message)
