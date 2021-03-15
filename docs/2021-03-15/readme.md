# March 15 2021

This works:

~~~
deezer.ping -> deezer.getUserData -> song.getListData -> get_url
                     |                                      ^
                     |______________________________________|
~~~

This does not work:

~~~
deezer.ping -> deezer.getUserData -> song.getData -> get_url
                     |                                      ^
                     |______________________________________|
~~~

So what is `song.getData` used for?

https://github.com/deknowny/async-deethon/issues/1
