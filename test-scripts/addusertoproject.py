#!/bin/python2

import requests
import json
import sys

if len(sys.argv) != 5:
    print("format: ./addusertoproject.py username password projectid newuserid")
    sys.exit(-1)

dat = { 'Email': sys.argv[1], 'Password': sys.argv[2], 'NewUserId': int(sys.argv[4])}
req = requests.post('https://productive.ahouts.com/project/' + sys.argv[3] + '/add_user',
                   headers={'Content-Type': 'application/json'},
                   data=json.dumps(dat).encode())
print(req.text)
