# Table data

数据库一张表的数据存储在一个目录中，每行都是一个JSON文件，可以手工维护。这个包将可以便利的读取这种特定结构的目录。

这个目录结构说明如下

```
/path/to/table                      # 目录
    __meta__.json                   # 这个表的元信息 (这个文件是可选的)
    {row}.json                      # 一行一个json文件
    {row}/index.json                # 也可以将一行装入到一个目录中
    {row}.{column}.{ext}            # 行中的某一个列可以储存在文件中，最后上传到ObjectStore中(参见sdobjectstore)
    {row}/{column}.{ext}            # 这些文件也可以装入一个目录中
    {row}.{column}.{sub}.{ext}      # 许多文件可以放在一列中，打包成json的链接形式
    {row}/{column}.{sub}.{ext}      # 这些文件也可以放在目录中
```

注意：这里`{row}`不是行的ID，只是在文件系统中保存行数据的目录名，真正的id写在row的数据中。

例如下面是具体的例子

```
local_data1/Product
└── 开关
    ├── image_url.white.png         # image_url字段中的white文件
    ├── image_url.red.png           # image_url字段中的red文件
    ├── index.json                  # 保存有这行数据的json
    └── logo_url.png                # logo_url字段中对应的文件

```