# build 방법

bunbler.json 에 다음을 추가한다.
```json
  "environments": [
    {"arch": "amd64", "os": "linux"},
    {"arch": "amd64", "os": "windows"},
    {"arch": "amd64", "os": "darwin"}
  ]
```

터미널에서 `astilectron-bundler` 명령어를 이용해서 빌드한다. 각 OS별 폴더가 생성되는데 해당 폴더에
`static` 폴더를 복사 붙여넣기한다. 

