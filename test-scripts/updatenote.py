#!/bin/python2

import requests
import json
import sys

if len(sys.argv) != 7:
    print("format: ./updatenote.py username password id title body ownerid")
    sys.exit(-1)

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2], 'Title': sys.argv[4], 'Body': sys.argv[5], 'OwnerId': int(sys.argv[6]), 'ProjectId': None}
req = requests.post('https://productive.ahouts.com/note/' + sys.argv[3],
                   headers={'Content-Type': 'application/json'},
                   data=json.dumps(dat).encode())
print(req.text)
