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
(`arl`) requests. What about logging in programmatically?

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

How do I tell what formats are available for a song?

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
