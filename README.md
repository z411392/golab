延續在其他語言的寫作習慣，使用 golang 製作的 boilerplate。

### 專案架構

#### main.go
主要的進入點。

執行後會用 corbra 解析命令列參數，再接著執行不同的 subcommand。

#### cmd
- 註冊 subcommands。

> 目前只有 `serve`、`auth` 兩個 command。之後如果要作 workers，不同的 queue 會有各自對應的 workers，每一類 workers 會有一個啟動用的 subcommand。

#### boot
- 註冊 http route / handlers
- 註冊 dependencies

> 之後若使用 amqp，也會註冊 queue / consumers。

#### resources
使用到的靜態檔案。

#### tests
撰寫測試的地方。分成：
- e2e
- integration
- unit

#### adapters
- 與外部系統介接的類別。

> 可以用 ports 作出抽象，註冊在 container，測試時改提供 mock 類別。

#### constants
是為了避免 magic number / strings。

目前只先放了 gin context 的 keys。

#### utils
將一些常用的 snippets 作成 functions 方便使用。

#### modules
按使用情境區分出不同的模組。


### APIs
![image](https://github.com/z411392/golab/blob/0ae336c9f6253a9c9d1c37ba5a6c642e579f1a7f/figure.jpg)

- GET /liveness_check
- GET /readiness_check
- GET /auth/me