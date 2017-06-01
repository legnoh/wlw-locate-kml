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
-U 1 > ~/result.diff
```

- ライブラリ設置状況が変わった店舗のみを抜き出す(次回以降対応)
```sh
$ diff result-before.kml result-after.kml \
--ignore-matching-lines ".*ランキング.*" \
--ignore-matching-lines=".*<name>.*" \
--ignore-matching-lines=".*description.*" \
--ignore-matching-lines=".*coordinates.*" \
--ignore-matching-lines=".*住所.*" \
--ignore-matching-lines=".*styleUrl.*" \
--ignore-matching-lines ".*ExtendedData.*" \
--ignore-matching-lines ".*Placemark>.*" \
--ignore-matching-lines ".*Point>.*" \
--ignore-matching-lines ".*SchemaData>.*" \
--ignore-matching-lines ".*schemaUrl.*" \
--ignore-matching-lines ".*name=\"住所\".*" \
-U 4 > ~/Desktop/hoge.diff
```


## TODO
- Refactoring
- UnitTest
- Deploy to CF

## ENJOY :meat_on_bone::meat_on_bone::meat_on_bone:
