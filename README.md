# Portto Interview Assignment
* 用 golang 實作兩個 Ethereum blockchain service
* 設計 postgres / mysql DB schema
* API service
* Ethereum block indexer service
# 前言
這是我第一次使用golang，gorm、gin、go-ethereum等都是沒看過的套件，
可能會有很多地方寫得很奇怪，因為是到處拼拼湊湊出來的...
# 用 golang 實作兩個 Ethereum blockchain service
`smart_contract/` folder

其實我不太理解題目，我在網路上查golang Ethereum等關鍵字，並沒有找到用golang寫出smart contract的相關方法，因此這部分我的實作是:
* 用solidity(`smart_contract/contracts/*`)先撰寫smart contract
* 用golang(`smart_contract/main.go`)去做smart contract的部署以及互動

題目說要實作兩個，但以這種做法來說，`smart_contract/main.go`做的事情是差不多的，所以`smart_contract/main.go`只會去對`smart_contract/contracts/guess_game.sol`去做佈署及互動，另外的`smart_contract/contracts/sample.sol`就只是提供一個常見的smart contract參考，沒有寫相關的main.go。

# 設計 postgres / mysql DB schema
`portto_shop_schema.sql`

這部分我是使用mysql，並且把DB schema export出來，大致上是設計一個簡單拍賣網頁的DB schema，內容如下(**粗體**表primary key, *斜體*為foreign key):
* `user`: (賣家或買家資訊)
  * **User_id**
  * Account
  * Password
  * Name
  * Register_time
  * Bill_info
  * Email
* `cart`: (紀錄甚麼商品及其數量在哪個user的購物車裡)
  * ***User_id***
  * ***Merchandise_id***
  * Number
* `merchandise`: (紀錄賣家(user)賣的商品資訊)
  * **Merchandise_id**
  * *User_id*
  * Name
  * Price
  * Number_in_stock
  * Number_sold
  * Description
* `order`:(訂單發生時產生的資訊)
  * **Order_id**
  * *User_id*
  * Order_time
* `orderitem`: (紀錄哪些商品及其數量價格被包含在哪個訂單裡)
  * ***Order_id***
  * ***Merchandise_id***
  * Trade_price
  * Number

# API service
`APIservice/` folder

這題大致上應該與題目要求做的差不多，沒什麼大問題，可能有的問題就是沒有時間做input的例外排除

#  Ethereum block indexer service
## 根據 web3 API 透過 RPC 將區塊內的資料掃進 db
### 需以平行的方式從 block n 開始進行掃描，直到掃到最新區塊後繼續執行
`block_indexer/` folder

這題是我問題最大的一題，問題如下:
* 首先是RPC api的問題，我發現透過本地geth輸入`curl -X POST --data '{"jsonrpc":"2.0", "method":"eth_getBlockByNumber","params":[idx, true],"id":1}'`這類的api以後，得到的結果與使用ethclient的`client.BlockByNumber(context.TODO(), idx)`得到的結果不同，似乎第一種得到的結果比較完整，但是我目前能力不足，沒辦法處理第一種回傳的資料，於是選擇使用第二種。
* 再來是db與golang溝通的問題，一樣是我對golang不夠熟悉，db schema寫得很醜以外，也做不出使用正確的型別來當作其column的type，另外也因此做不出以block id做為transaction的foreign key。
* 我不確定是不是要使用`client.BlockByNumber`，但我找不到有可以batch去做這件事的function
* 我發現golang與其他語言不同的地方在於，他是使用goroutine，與我熟悉的python不同。因此本來我的實作是想要用goroutine pool(thread pool)來達成，但是嘗試了許久以後終究有幾個問題解決不了:
  * 首先是參數傳入goroutine task以後，一直會有共用記憶體的問題，導致執行結果不太對
  * 不同的參數傳入`client.BlockByNumber`並且於goroutine pool中執行後得到的結果會不同，我在想會不會是`client.BlockByNumber`本身不是thread safe的，這就代表我應該使用別的方法來與eth溝通，但基於最開始的問題，我沒辦法做到。

程式本身是能做到大致上的功能，有做到平行從block n開始掃到最後並繼續執行，但基於以上問題，是個失敗的程式。
