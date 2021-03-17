# March 16 2021

With this server:

~~~
https://media.deezer.com/v1/get_url
~~~

If you send unauthenticated (`deezer.ping`) request for `MP3_320`, you get this
response:

~~~
License token has no sufficient rights on requested media.
~~~

Same for authenicated (`arl`) with "Deezer Free". So "Deezer Premium" would be
required with this server. So in order to get the URL, we need `MD5_ORIGIN`.
Unauthenticated (`deezer.ping`) requests do not return this, only authenicated
(`arl`) requests. What about logging in programmatically? Here are the requests
that are made:

~~~
{
  method: 'get',
  url: 'https://api.deezer.com/1.0/gateway.php?method=mobile_auth&api_key=4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE&output=3&uniq_id=b1bffb3423b0bd6bba1a9019f595a5e6',
  headers: {
    'User-Agent': 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox'
  }
}
{
  method: 'get',
  url: 'https://api.deezer.com/1.0/gateway.php?method=api_checkToken&api_key=4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE&output=3&auth_token=84bce3f5f024528422030fe15a073cf639c9343207f173d9d1017f9646a7814536f74daf1798cd3c6ccc581b9ea7920325227e6ab4f4f8d71ce41de590e7d796',
  headers: {
    'User-Agent': 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox'
  }
}
{
  method: 'post',
  url: 'https://api.deezer.com/1.0/gateway.php?method=mobile_userAuth&api_key=4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE&output=3&input=3&sid=frc7f4658ee4475182f717fcfa8995dafbb6650c',
  headers: {
    'User-Agent': 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox'
  },
  data: {
    mail: 'srpen6@gmail.com',
    password: 'd91e4efdcf834c5c914071e007ae99b1',
    device_serial: '',
    platform: 'innotek GmbH_x86_64_9',
    custo_version_id: '',
    custo_partner: '',
    model: 'VirtualBox',
    device_name: 'VirtualBox',
    device_os: 'Android',
    device_type: 'tablet',
    google_play_services_availability: '1',
    consent_string: ''
  }
}
{
  method: 'post',
  url: 'https://api.deezer.com/1.0/gateway.php?method=mobile.pageAlbum&api_key=4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE&output=3&input=3&sid=frc7f4658ee4475182f717fcfa8995dafbb6650c',
  headers: {
    'User-Agent': 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox'
  },
  data: {
    alb_id: '7476543',
    user_id: 4239375622,
    lang: 'en',
    header: true,
    tab: 0
  }
}
~~~

This is interesting:

~~~
REQUEST {
  method: 'get',
  url: 'https://api.deezer.com/1.0/gateway.php?method=mobile_auth&api_key=4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE&output=3&uniq_id=64dfbd9ab81f90fd3ed5f39a8909f442',
  headers: {
    'User-Agent': 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox'
  }
}
RESPONSE {
  error: [],
  results: {
    CONFIG: {
      HOST_STREAM_MOB: 'http://e-cdn-proxy-{0}.deezer.com/mobile/1/',
      HOST_CUSTOS_MOB: 'http://e-cdn-content.dzcdn.net/mobile/custos/',
      HOST_API_MOB: 'https://api.deezer.com/1.0/gateway.php',
      HOST_FILES: 'http://e-cdn-files.dzcdn.net',
      HOST_CONTENT: 'http://e-cdn-content.dzcdn.net',
~~~

Here are the formats:

int    | string
-------|---------
0      | MP3_MISC
1      | MP3_128
2 4 16 | SBC_256
3      | MP3_320
5      | MP3_256
6      | AAC_64
7      | MP3_192
8      | AAC_96
9      | FLAC
10     | MP3_64
11     | MP3_32
12     | SBC
13     | MP4_RA1
14     | MP4_RA2
15     | MP4_RA3

How do I tell what formats are available for a song? Like this:

~~~
REQUEST {
  method: 'post',
  url: 'https://api.deezer.com/1.0/gateway.php?method=song_getData&api_key=4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE&output=3&input=3&sid=frde28a171f0faaea4ebd1f5a704bbb2a9b20f1a',
  headers: {
    'User-Agent': 'Deezer/6.1.22.49 (Android; 9; Tablet; us) innotek GmbH VirtualBox'
  },
  data: { SNG_ID: '75498415' }
}
RESPONSE
{
  results: {
    FILESIZE_AAC_64: '0',
    FILESIZE_MP3_64: '0',
    FILESIZE_MP3_128: '3909589',
    FILESIZE_MP3_256: '0',
    FILESIZE_MP3_320: '0',
    FILESIZE_MP4_RA1: '0',
    FILESIZE_MP4_RA2: '0',
    FILESIZE_MP4_RA3: '0',
    FILESIZE_FLAC: '0',
  }
}
~~~

Are they all actually available though?

int    | string   | gw-light.php | get_url
-------|----------|--------------|--------
1      | MP3_128  | good         | good
3      | MP3_320  | good         | good
9      | FLAC     | good         | bad
0      | MP3_MISC |              | ugly
10     | MP3_64   |              | ugly
12     | SBC      |              | ugly
13     | MP4_RA1  |              | bad
14     | MP4_RA2  |              | bad
15     | MP4_RA3  |              | bad
2 4 16 | SBC_256  |              | bad
6      | AAC_64   |              | bad
8      | AAC_96   |              | bad
5      | MP3_256  |              | bad
7      | MP3_192  |              | bad
11     | MP3_32   |              | bad

`get_url` results are with:

~~~
OFFER_NAME: 'Deezer Premium'
RENEWAL: '1'
~~~

Lets see if anyone has a way to get this:

~~~
MP3_256
~~~

Issues:

- https://github.com/TheDeezLoader/pyDeezloader/issues/1
- https://github.com/acgonzales/pydeezer/issues/17
- https://github.com/cumulus13/deezersdk/issues/1
- https://github.com/frannyfx/freezer/issues/4
