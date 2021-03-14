from urllib import parse
import os, requests
session = requests.Session()

PRIVATE_API_URL = 'https://www.deezer.com/ajax/gw-light.php?'
args = {
   'api_version': '1.0',
   'api_token': 'null',
   'method': 'deezer.ping'
}
res = session.get(PRIVATE_API_URL + parse.urlencode(args))

PUBLIC_API_URL = 'https://api.deezer.com/1.0/gateway.php?'
args = {
   'api_key': os.environ['DZR_API'],
   'method': 'song_getData',
   'input': '3',
   'output': '3',
   'sid': res.json()['results']['SESSION']
}
res = session.post(PUBLIC_API_URL + parse.urlencode(args), data='{"sng_id":75498415}')
print(res.text)
