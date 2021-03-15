# March 15 2021

This works:

~~~
deezer.ping -> deezer.getUserData -> song.getListData -> get_url
                     |                                      ^
                     |______________________________________|
~~~

Does this:

~~~
deezer.ping -> deezer.getUserData -> song.getData -> get_url
                     |                                      ^
                     |______________________________________|
~~~
