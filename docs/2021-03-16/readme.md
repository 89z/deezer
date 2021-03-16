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
