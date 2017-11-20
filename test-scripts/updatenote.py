#!/bin/python2

import requests
import json
import sys

if len(sys.argv) != 7 and len(sys.argv) != 8:
    print("format: ./updatenote.py username password id title body ownerid <optional project id>")
    sys.exit(-1)

projectid = None
if len(sys.argv) == 8:
    projectid = {
        "Int64": int(sys.argv[7]),
        "Valid": True
    }

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2], 'Title': sys.argv[4], 'Body': sys.argv[5], 'OwnerId': int(sys.argv[6]), 'ProjectId': projectid}
req = requests.post('https://productive.ahouts.com/note/' + sys.argv[3],
                   headers={'Content-Type': 'application/json'},
                   data=json.dumps(dat).encode())
print(req.text)
