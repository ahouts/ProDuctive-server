#!/bin/python

import urllib.request
import json
import sys

if len(sys.argv) != 3:
    print("format: ./insertuser.py username password")
    sys.exit(-1)

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2]}
req = urllib.request.Request('https://productive.ahouts.com/reminder/get')
req.add_header('Content-Type', 'application/json')
req.get_method = lambda: "PUT"
try:
    response = urllib.request.urlopen(req, json.dumps(dat).encode())
    print(response.read())
except urllib.request.HTTPError as e:
    error_message = e.read()
    print(error_message)
