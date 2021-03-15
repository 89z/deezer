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

- https://github.com/aeggydev/deezer-tool/issues/1
- https://github.com/deknowny/async-deethon/issues/1
