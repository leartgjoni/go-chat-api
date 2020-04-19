![build](https://github.com/leartgjoni/go-chat-api/workflows/build/badge.svg)
[![codecov](https://codecov.io/gh/leartgjoni/go-chat-api/branch/master/graph/badge.svg?token=IV6P28SQUX)](https://codecov.io/gh/leartgjoni/go-chat-api)

# go-chat-api
Chat built in Go using Websockets and Redis Pub/Sub

- The redis pub/sub system enables the api to scale horizontally.

- The CI is on github-actions. Check `./github/workflows/main.yml`

- The api is deployed into a k8s cluster in AWS. Check the <a href="https://github.com/leartgjoni/go-chat-k8s">infrastructure</a>.

- The <a href="https://github.com/leartgjoni/go-chat-app">frontend</a>. is built in React

Scripts can be found in `Makefile`, e.g. running api locally:
```
make start-local
```

#### Demo
<img src="https://raw.githubusercontent.com/leartgjoni/go-chat-api/master/demo/demo.gif" />
