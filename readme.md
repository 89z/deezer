# Deezer

[![Go Reference]](//pkg.go.dev/github.com/89z/deezer)

[Go Reference]:https://pkg.go.dev/static/img/badge.svg

> Listen, um, I don’t really know what’s going on here, but I’m leaving. I
> don’t know where exactly, but I’m gonna start over.
>
> Come with me. Please.
>
> [Paint it Black (2016)](//wikipedia.org/wiki/Paint_It_Black_%282016_film%29)

Download audio from Deezer

## Album

http://api.deezer.com/album/7476543

## Thanks

- https://github.com/moon004/Go-deezer-downloader
- https://pkg.go.dev/github.com/godeezer/lib/deezer
- https://pkg.go.dev/github.com/joshbarrass/deezerdl/pkg/deezer

## Cookies

1. http://www.deezer.com/login
2. Web Developer
3. Network
4. Log in
5. Network Settings
6. Save All As HAR

This will give you an `arl` cookie that is valid for six months. The program will
then make a `GET` request with the `arl` cookie, which will return a `sid`
session cookie.
