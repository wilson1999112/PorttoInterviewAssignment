# Portto Interview Assignment
* 用 golang 實作兩個 Ethereum blockchain service
* 設計 postgres / mysql DB schema
* API service
* Ethereum block indexer service
# 前言
大概學了幾天的golang以後，稍微寫得有點感覺了
# API service
`APIservice/` folder

這題大致上應該與題目要求做的差不多，沒什麼大問題，可能有的問題就是沒有時間做input的例外排除

#  Ethereum block indexer service
## 根據 web3 API 透過 RPC 將區塊內的資料掃進 db
### 需以平行的方式從 block n 開始進行掃描，直到掃到最新區塊後繼續執行
`block_indexer/` folder
經過從基本開始學習golang以後，對goroutine有了解一些些，基本上有按照題目要求的執行並經過測試
但DB schema的部分還是有點生疏，因此做的DB是簡化版本，有些value沒有塞進去

程式基本上會有以下四個階段:
1. 先把穩定的區塊掃進db(`stable=true`)
2. 開始掃不穩定的區塊並追蹤其block hash(`stable=false`)
3. 當最新區塊變多了以後，將最前面的不穩定區塊標記為`stable=true`，並更新資料庫
   *  此時若有發現其最新的block hash與原本掃進db的不一樣，則刪除並重新插入最新資料
4. 重複第三步驟

