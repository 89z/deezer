# SID

Instead of using `arl` cookie, you can use `sid` cookie. This allows you to avoid
making a `deezer.getUserData` request. The question is, how long is the cookie
good for? Here are results for time since SID generation:

duration | time                          | result
---------|-------------------------------|-------
0        | 2021-03-11T10:16:22.938-06:00 | good
4h0m0s   | 2021-03-11 14:16:22.938 -0600 | good
8h0m0s   | 2021-03-11 18:16:22.938 -0600 | bad

Here are results for time since last request:

duration | time                          | result
---------|-------------------------------|-------
0        | 2021-03-11T17:34:01.362-06:00 | good
2m       | 5:37 PM                       | good
16m      | 5:54 PM                       | good
32m      | 6:41 PM                       | good
1h       | 7:41 PM                       | good
2h       |
4h       |

https://github.com/kmille/deezer-downloader/issues/49
