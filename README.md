**August 2021**  |  **v1.0**


Bofry/host-fasthttp 快速入門指南
==================================

----------------------------------------------------------------
#### 目錄
  - [簡介](#簡介)
  - [環境要求](#環境要求)
  - [快速入門](#快速入門)
  - [相關資源](#相關資源)
----------------------------------------------------------------

## 簡介

⠿ **Bofry/host-fasthttp** 是使用 [*valyala/fasthttp*](https://github.com/valyala/fasthttp) 作為主體，專門提供 golang 開發者快速建置 *WebAPI* 應用服務為目的而製作的框架。主要特色如下：

1. <span class="underline">**非全功能的框架**</span>，主要目的是快速開發 *WebAPI* 服務，且目前並沒有計畫增加 *WebServer* 功能的相關設計，但仍可以參考 *valyala/fasthttp* 的方式來開發 *WebServer* 的應用程式。
2. 框架遵循 <span class="underline">[**RMR (*Reource-Method-Representation*)**](https://www.peej.co.uk/articles/rmr-architecture.html) **語意**</span>來作為應用程式開發配置的原則。
3. 應用程式<span class="underline">**使用自訂 HTTP method**</span> 來組織應用程式的請求呼叫，比如：`PEEK`、`FETCH`、`COMMIT`…等，而非傳統的 `GET`、`POST`、`DELETE`…等。
4. *Bofry/host-fasthttp* 所支援的內建功能:
   1. 自動配置 *Route Table*。
   2. 通用的應用程式設定（*config*）處理；支援解析環境變數、*JSON*、*YAML* 以及命令列參數。部份類型支援設定值驗證。
   3. 通用的 *HTTP querystring*、*body* 等輸入的解析與驗證。
   4. 額外提供處理請求非標準 *HTTP method* 的相容性標頭。
   5. 提供額外的 *Error Response Handler*，處理失敗的 *HTTP* 請求回應。
   6. 提供額外的 *Logging* 支援。
   7. 內建安全停止（*Graceful Shutdown*）機制。
   8. 框架使用 *DI* 機制來減低耦合與簡化操作複雜度。

[🔝回目錄](#%e7%9b%ae%e9%8c%84)


----------------
## 環境要求

1. GO 1.14 以上

    下載及安裝指引：https://golang.org/doc/install

2. Git

3. 將 GOPATH 的 bin 目錄加入你的 system PATH 變數中

    - Linux 系統於 `~/.bashrc` 中加入
      ```bash
      export PATH=<your_gopath_bin_dir>:$PATH
      ```

    - Mac 系統於 `~/.bash_profile` 中加入
      ```bash
      export PATH=<your_gopath_bin_dir>:$PATH
      ```

    - Windows 系統於系統環境變數中
      1. 按 `Windows Key` + `R` 開啟執行對話框
      2. 輸入 `control sysdm.cpl,,3`
      3. 按 **環境變數(N)** 按鈕 或 `Alt` + `N`
      4. 在 **系統變數(S)** 方塊中，修改 **Path** 變數，加入 GOPATH 的 bin 目錄

    > **⠿ 取得 GOPATH 的 bin 完整路徑，依下列操作步驟**：
    >
    > 使用 `go env GOPATH` 查看 GOPATH 後，在後面加上 `/bin` 即是 GOPATH 的 bin 目錄，例：
    >   ```bash
    >   $ go env GOPATH
    >   /Users/go
    >   ```
    > 因此推得 GOPATH bin 目錄的完整路徑為 _/User/go/bin_

4. 安裝所需的 go tool

    - Go 1.15 含以前使用下面指令：
    ```bash
    go get -v github.com/Bofry/go-tools/gorun
    go get -v github.com/Bofry/go-tools/host-fasthttp
    go get -v github.com/Bofry/go-tools/gen-host-fasthttp-resource
    go get -v github.com/joho/godotenv/cmd/godotenv
    ```
    - Go 1.16 含以後使用下面指令：
    ```bash
    go install github.com/Bofry/go-tools/gorun
    go install github.com/Bofry/go-tools/host-fasthttp
    go install github.com/Bofry/go-tools/gen-host-fasthttp-resource
    go install github.com/joho/godotenv/cmd/godotenv
    ```

    - Go 1.16 含以後使用下面指令：
    ```bash
    go install github.com/Bofry/go-tools/gorun
    go install github.com/Bofry/go-tools/host-fasthttp
    go install github.com/Bofry/go-tools/gen-host-fasthttp-resource
    go install github.com/joho/godotenv/cmd/godotenv
    ```
    > 🐾 see *Deprecation of 'go get' for installing executables* <sup>[🔗](https://go.dev/doc/go-get-install-deprecation)</sup>

[🔝回目錄](#%e7%9b%ae%e9%8c%84)


----------------
## 快速入門

⠿ 本節提供快速建構 **myapp** 應用程式的基本步驟。


- **步驟一: 決定你的專案資料夾**

  ⠿ 建立名稱為 **myapp** 的專案。
    ```bash
    $ mkdir myapp

    $ cd myapp
    ```

- **步驟二: 初始化你的專案**
  1. 使用 `go mod init` 初始化 golang 專案。
      ```bash
      $ go mod init apiservice
      ```
  2. 使用 `host-fasthttp init` 進行專案的基本配置。
      ```bash
      $ host-fasthttp init
      ```
      > 💬 `host-fasthttp init` 會產生 .env、app.go、config.yaml、.conf/、internal/appContext.go……等專案檔，並進行 http server 最基本的配置。

- **步驟三: 加入預設的 http request handler**
  1. 使用編輯器打開 **app.go** 檔。
  2. 找到第11行，
      ```go
      type ResourceManager struct {}
      ```
      更改為：
      ```go
      type ResourceManager struct {
        *DefaultResource `url:"/"`
      }
      ```
      > 💬 若你使用 IDE，比如 vscode 則可能會標示語法錯誤，請先忽略。
  3. 執行 `go generate` 來配置 http request handler 檔案。
      ```bash
      $ go generate
      ```
      > 💬 這個動作為幫你產生 **defaultResouce.go** 的檔案，並會放在專案內的 `/resource` 目錄中。
      >
      > 💬 若使用 IDE 且支援 go generate 指令的話，可以在程式碼上方直接按下 `run go generate` 的按鈕，也會有相同的效果。

- **步驟四: 啟動專案**
  ```bash
  $ gorun
  ```
  > 💬 `gorun` 是一個結合 `godotenv` 與 `go run` 命令的工具。
  >
  > 💬 使用 `Ctrl` + `C` 來關閉服務。

- **步驟五: 檢查**
  ```bash
  $ curl -XPING -sv http://127.0.0.1:10074/
  ```
  將會得到以下的輸出：
  ```bash
  PONG
  *   Trying 127.0.0.1:10074...
  * Connected to 127.0.0.1 (127.0.0.1) port 10074 (#0)
  > PING / HTTP/1.1
  > Host: 127.0.0.1:10074
  > User-Agent: curl/7.71.1
  > Accept: */*
  >
  * Mark bundle as not supporting multiuse
  < HTTP/1.1 200 OK
  < Server: WebAPI
  < Date: Mon, 23 Aug 2021 07:42:43 GMT
  < Content-Type: text/plain
  < Content-Length: 4
  < Connection: close
  <
  { [4 bytes data]
  * Closing connection 0
  ```

----------------
## 相關資源

  1. 有關 `go mod` 的名稱規範詳見 `go mod init` 的官方說明[連結](https://golang.org/ref/mod#go-mod-init)。


[🔝回目錄](#%e7%9b%ae%e9%8c%84)
