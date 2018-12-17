## sc

## 调用方法
```console
sc.exe -input=input.xlsx -output=output.xlsx
```

### 注意事项
1.  input.xlsx 格式见同目录下input.xlsx, 其目的是指明要采集的股票ID
2.  output.xlsx 如果不存在，则根据input.xlsx创建，且追加写
3.  output.xlsx 如果存在，其格式也和input.xlsx, 工具将追加写
4.  前两列时间为调用工具时的日期和时间
5.  如果调用时某股票ID并无数据返回，或者API调用失败，则使用前一行的价格填充，如果前一行不存在，则用0填充