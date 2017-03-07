# Govenze

Check vendor size for repo. Checking with:

* [godep](https://github.com/tools/godep)
* [glide](https://github.com/Masterminds/glide), using with `--skip-test`
* [dep](https://github.com/golang/dep)
* [govendor](https://github.com/kardianos/govendor)

## Result examples

Result for github.com/chapsuk/frissgo

Tool | govendor | godep | dep | glide
 --- | --- | --- | --- | ---
Size | 0.93Mb | 0.94Mb | 6.93Mb | 6.99Mb
Time | 0.66s | 1.17s | 14.79s | 26.30s
Builded | true | true | true | true

Result for github.com/tools/godep

Tool | govendor | godep | glide | dep
 --- | --- | --- | --- | ---
Size | 0.08Mb | 0.08Mb | 7.99Mb | 7.99Mb
Time | 0.75s | 1.31s | 14.76s | 17.75s
Builded | true | true | true | true

Result for github.com/Masterminds/glide

Tool | godep | govendor | dep | glide
 --- | --- | --- | --- | ---
Size | 0.49Mb | 0.49Mb | 0.73Mb | 0.73Mb
Time | 1.30s | 0.92s | 6.47s | 10.37s
Builded | true | true | true | true

Result for github.com/kardianos/govendor

Tool | govendor | godep | glide | dep
 --- | --- | --- | --- | ---
Size | 3.33Mb | 3.33Mb | 12.04Mb | 12.04Mb
Time | 0.93s | 2.31s | 22.16s | 23.81s
Builded | true | true | true | true

Result for github.com/golang/dep

Tool | govendor | godep | dep | glide
 --- | --- | --- | --- | ---
Size | 0.45Mb | 0.45Mb | 0.80Mb | 0.80Mb
Time | 1.01s | 1.06s | 7.36s | 7.83s
Builded | false | false | false | false

![](https://media.giphy.com/media/Yo9Xldk1F5196/giphy.gif)

## Notices

`Godep` and `govendor` get vendors from local `GOPATH`.
`glide` and `dep` get vendors from source.

1. `godep` missing C binidings vendoring, [issue](https://github.com/tools/godep/issues/422)
1. `govendor` throw `fatal error: stack overflow` if packages duplicates in `GOPATH`
1. `glide` does not always use the latest version of the dependency.
