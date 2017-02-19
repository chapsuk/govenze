# Govenze

Check vendor size for repo. Checking with:

* [godep](https://github.com/tools/godep)
* [glide](https://github.com/Masterminds/glide), using with `--skip-test`
* [dep](https://github.com/golang/dep)

## Result examples

[Frissgo](https://github.com/chapsuk/frissgo)

```bash
≻ govenze -target github.com/chapsuk/frissgo

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
```

[Viper](https://github.com/spf13/viper)

```bash
≻ govenze -target github.com/spf13/viper

Vendor manager: godep
======
Size: 6.37Mb
Time: 5.81s


Vendor manager: glide
======
Size: 27.11Mb
Time: 39.99s


Vendor manager: dep
======
Size: 26.63Mb
Time: 41.48s

```

![](https://media.giphy.com/media/Yo9Xldk1F5196/giphy.gif)

## Notices

1. Godep missing C binidings vendoring, [issue](https://github.com/tools/godep/issues/422)
