## stock_crawler

调用新浪API查询输入表格中指定股票ID的最新报价，并追加到输出表格的工具

## install
```console
go get -u github.com/tianlin/stock_crawler
dep ensure
go install github.com/tianlin/stock_crawler/sc
```

### usage
```console
sc -input=input.xlsx -output=output.xlsx
```