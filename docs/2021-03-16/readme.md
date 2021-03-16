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
required with this server.
