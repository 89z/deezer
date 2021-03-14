from urllib import parse
import os, requests
req = requests.Session()
url = 'http://www.deezer.com/ajax/gw-light.php?'

# ping
arg = {
   'api_token': '',
   'api_version': '1.0',
   'method': 'deezer.ping'
}
res = req.get(url + parse.urlencode(arg))
session = res.json()['results']['SESSION']
print(session)

# getUserData
arg = {
   'api_token': '',
   'api_version': '1.0',
   'input': '3',
   'method': 'deezer.getUserData'
}
head = {
   'Cookie': 'sid=' + session
}
res = req.get(url + parse.urlencode(arg), headers=head)
check = res.json()['results']['checkForm']
print(check)

# pageTrack
arg = {
   'api_token': check,
   'api_version': '1.0',
   'method': 'deezer.pageTrack'
}
res = req.post(url + parse.urlencode(arg), data='{"sng_id":75498415}')
print(res.text)
