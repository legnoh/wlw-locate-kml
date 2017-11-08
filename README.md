# wlw-locate-kml

- [WonderlandWars設置店一覧](https://www.google.com/maps/d/viewer?mid=1ENDxk6QqlKlyjqS4iB_1HNyD7UM)

output WLW map's KML file.

## usage
```sh
$ glide install
$ go run main.go
```

### for diff...
- 新規出店・退店店舗のみを抜き出す
```sh
$ diff result-before.kml result-after.kml \
--ignore-matching-lines=".*name.*" \
--ignore-matching-lines=".*description.*" \
--ignore-matching-lines=".*coordinates.*" \
--ignore-matching-lines ".*ランキング.*" \
--ignore-matching-lines ".*styleUrl.*" \
-U 0 > frequent.diff
```

## TODO
- Refactoring
- UnitTest
- Deploy to CF

## ENJOY :meat_on_bone::meat_on_bone::meat_on_bone:
