from urllib import parse
import json, requests
session = requests.Session()
PUBLIC_API_URL = 'https://api.deezer.com/1.0/gateway.php?'

arg = {
   'api_key': '4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE',
   'method': 'mobile_userAuth',
   'output': '3',
   'sid': 'fr737eb0c0365971df92b1d18554382a04ab95b1',
}

body={
   'MAIL': 'srpen6@gmail.com',
   'password': '78e715eb052ff77f9b92f3eb59a1360e',
}

res = session.post(
   PUBLIC_API_URL + parse.urlencode(arg), data=json.dumps(body)
)

print(json.dumps(res.json(), indent=1))
