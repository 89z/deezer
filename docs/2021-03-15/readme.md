# March 15 2021

## Fail

~~~
deezer.ping -> deezer.getUserData -> song.getData -> get_url
                     |                                  ^
                     |__________________________________|
~~~

## Pass cookie

~~~
deezer.getUserData -> deezer.pageTrack
~~~

~~~
deezer.getUserData -> song.getData
~~~

## Pass no cookie

~~~
deezer.ping -> deezer.getUserData -> song.getListData -> get_url
                     |                                      ^
                     |______________________________________|
~~~

~~~
deezer.ping -> deezer.getUserData -> deezer.pageTrack -> get_url
                     |                                      ^
                     |______________________________________|
~~~

Now what about this:

~~~
http://api-v3.deezer.com/1.0/gateway.php
http://api.deezer.com/1.0/gateway.php
deezer gateway api_key
~~~

https://github.com/svbnet/diezel

Here are some methods used:

~~~
shouldn't need to call this normally
api_checkToken

shouldn't need to call this normally
mobile_auth

Signs in with an email and password
mobile_userAuth

Restores a signed-in session
mobile_userAutoLog

Gets a song by ID
song_getData
~~~
