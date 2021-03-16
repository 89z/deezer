from urllib import parse
import json, requests
session = requests.Session()
PUBLIC_API_URL = 'https://api.deezer.com/1.0/gateway.php?'

arg = {
   'method': 'mobile_userAuth',
   'api_key': '4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE',
   'output': '3',
   'input': '3',
   'sid': 'frc62578edaab2637431319896cab1944a3ba637',
}

body={
   'MAIL': 'srpen6@gmail.com',
   'consent_string': ''
   'custo_partner': '',
   'custo_version_id': '',
   'device_name': 'VirtualBox',
   'device_os': 'Android',
   'device_serial': '',
   'device_type': 'tablet',
   'google_play_services_availability': '1',
   'model': 'VirtualBox',
   'password': 'd8c44e320f3e7d2771596f56ae12e7c7',
   'platform': 'innotek GmbH_x86_64_9',
}

res = session.post(
   PUBLIC_API_URL + parse.urlencode(arg),
   headers={'User-Agent': 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox'},
   data=json.dumps(body)
)

print(json.dumps(res.json(), indent=1))
