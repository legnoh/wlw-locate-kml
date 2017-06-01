# wlw-locate-kml

- [WonderlandWars設置店一覧](https://www.google.com/maps/d/viewer?mid=1ENDxk6QqlKlyjqS4iB_1HNyD7UM)

output WLW map's KML file.

## usage
```sh
$ glide install
$ go run main.go
```

for diff...
```sh
$ diff result-before.kml result-after.kml \
--ignore-matching-lines "SimpleData*" \
--ignore-matching-lines ".*ランキング結果.*"
--ignore-matching-lines=".*styleUrl.*" \
--ignore-matching-lines=".*name.*" \
--ignore-matching-lines=".*coordinates.*" > ~/Desktop/hoge.diff
```


## TODO
- Refactoring
- UnitTest
- Deploy to CF

## ENJOY :meat_on_bone::meat_on_bone::meat_on_bone:
