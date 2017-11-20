#!/bin/python2

import requests
import json
import sys

if len(sys.argv) != 5 and len(sys.argv) != 6:
    print("format: ./createnote.py username password title body <optional projectid>")
    sys.exit(-1)

projectid = None
if len(sys.argv) == 6:
  projectid = {
                 "Int64": int(sys.argv[5]),
                 "Valid": True
              }
dat = { 'Email': sys.argv[1], 'Password': sys.argv[2], 'Title': sys.argv[3], 'Body': sys.argv[4], 'ProjectId': projectid}
print(dat)
req = requests.post('https://productive.ahouts.com/note',
                   headers={'Content-Type': 'application/json'},
                   data=json.dumps(dat).encode())
print(req.text)
