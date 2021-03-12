# Gateway

Currently I am using this gateway:

http://www.deezer.com/ajax/gw-light.php

This gateway requires an `api_token`. You can get this value from a HAR file,
one or more of the requests in the file will have it listed in the query string.
Another gateway is available as well:

https://api.deezer.com/1.0/gateway.php

This gateway requires an `api_key`. Here are two I found online:

~~~
4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE
ZAIVAHCEISOHWAICUQUEXAEPICENGUAFAEZAIPHAELEEVAHPHUCUFONGUAPASUAY
~~~

The question is, how were these derived? If we cannot answer this, we should
continue to use the other gateway.

https://github.com/BackInBash/DeezerAPI

The Deezer app is here:

https://microsoft.com/p/deezer-music/9nblggh6j7vv

How can I download from the Microsoft Store?
