# Govenze

Check vendor size for repo. Checking with:

* [godep](https://github.com/tools/godep)
* [glide](https://github.com/Masterminds/glide)
* [dep](https://github.com/golang/dep)

## Result (2017-02-19)

```bash
≻ /govenze -target github.com/chapsuk/frissgo
2017/02/19 19:41:16 Detected GOPATH: /Users/mak/go
2017/02/19 19:41:16 Create tmp dir: /Users/mak/.govenze

Vendor manager: godep
======
Size: 0.76Mb
Time: 1.60s


Vendor manager: glide
======
Size: 7.48Mb
Time: 16.51s


Vendor manager: dep
======
Size: 6.65Mb
Time: 30.68s

2017/02/19 19:41:47 Tmp dir /Users/mak/.govenze deleted
```

## Notices

1. Godep missing C binidings vendoring, [issue](https://github.com/tools/godep/issues/422)
