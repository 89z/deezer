# March 13 2021

So it seems this API may not work anymore:

https://api.deezer.com/1.0/gateway.php

I want to see if anyone is using it successfully.

~~~
deezer gateway api_key
~~~

- https://github.com/BackInBash/DeezerAPI/issues/4
- https://github.com/IvanMMM/bp-tags-grabber/issues/1
- https://github.com/Kaporos/Deezpy/issues/1
- https://github.com/L0v4iy/deezer-api-client/issues/1
- https://github.com/svbnet/diezel/issues/6
- https://github.com/yne/dzr/issues/8

To get audio URL, we need `license_token` value. To get `license_token` value,
we need `session` value. To get `session` value, we make request to:

~~~
http://www.deezer.com/ajax/gw-light.php?method=deezer.ping&api_version=1.0&api_token
~~~

Can I use this `session` value with the old process?
