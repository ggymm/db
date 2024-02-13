```shell
go install modernc.org/goyacc@latest
```

```shell
goyacc -o yy_parser.go -v yacc.output sql.y
```


### 语法规则

{%
嵌入代码
%}
文法定义
%%
文法规则
%%
嵌入代码


%union	    用来定义一个类型并映射golang的一个数据类型（可以是一个自定义类型）
%struct	    同%union 建议使用%union
%token	    定义非终结符 是一个union中定义的类型空间 可无类型空间
%type	    定义终结符
%start	    定义从哪个终结符开始解析 默认规则段中的第一个终结符
%left	    定义规则结合性质 左优先
%right	    定义规则结合性质 右优先
%nonasso	定义规则结合性质 不结合
%perc term	定义优先级与 term 一致
