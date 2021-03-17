# March 16 2021

Of the formats listed around the web, which are actually available?

int    | string   | gw-light.php | get_url
-------|----------|--------------|--------
1      | MP3_128  | good         | good
3      | MP3_320  | good         | good
9      | FLAC     | good         | bad
13     | MP4_RA1  |              | bad
14     | MP4_RA2  |              | bad
15     | MP4_RA3  |              | bad
2 4 16 | SBC_256  |              | bad
6      | AAC_64   |              | bad
8      | AAC_96   |              | bad
5      | MP3_256  |              | bad
7      | MP3_192  |              | bad
11     | MP3_32   |              | bad
0      | MP3_MISC |              | ugly
10     | MP3_64   |              | ugly
12     | SBC      |              | ugly

`get_url` results are with:

~~~
https://media.deezer.com/v1/get_url
OFFER_NAME: 'Deezer Premium'
RENEWAL: '1'
~~~

- https://github.com/TheDeezLoader/pyDeezloader/issues/1
- https://github.com/acgonzales/pydeezer/issues/17
- https://github.com/cumulus13/deezersdk/issues/1
- https://github.com/frannyfx/freezer/issues/4
