# Генерация кода по SQL запросам

при помощи XO/XO

go get -u golang.org/x/tools/cmd/goimports

go get -u github.com/xo/xo@v0.0.0-20210416025017-9a3ddc1e1407

go install github.com/xo/xo@v0.0.0-20210416025017-9a3ddc1e1407
mv $GOPATH/bin/xo $GOPATH/bin/xo-v0.0.0-9a3ddc1e1407

go get -u ./...
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o xo -v main.go

# Правила

## Генератор xo

- в каталоге scripts следующая иерархия папок
  + scripts
    + xo
      + <DB_NAME>
        + sql
          + query
            + many (SQL запросы выбирающие массив)
            + one  (SQL запросы выбирающие одно значение)
            + uid  (SQL запросы UPDATE/INSERT/DELETE)
          + templates (шаблоны генерации кода)
          
- Генерация запускается из каталога scripts/xo/<DB_NAME> 
```shell script
make gen
```

## Имена файлов

имя sql файла

<Возвращаемый тип>-<Имя функции>.sql

<Возвращаемый тип> также влияет на имя генерируемого файла с кодом

на выходе сформируется файл с именем в нижнем регистре
<Возвращаемый тип>.xo.go
если SQL запросов с одинаковым типом несколько то первый запрос создает файл с кодом,
остальные добавляют код в конец файла

```
/sql/query/many/ - возвращается тип []*<Возвращаемый тип>
/sql/query/one/ - возвращается тип *<Возвращаемый тип>
/sql/query/uid/ - возвращается тип error 
```

## Параметры запроса

Для формирования параметров функции в GO в sql запросе
```
$1, $2, $3 ... 
```
заменяются на конструкцию 
```
%%<Name> <Go-Type>%%
```
например
```
%%vatRate float64%%
``` 
в результате функция будет иметь входной параметр
```
func GetSomeExample(.., vatRate float64 ...)
```


## Особенности XO

### Возврат простых типов int, string, etc

пока я не добился генерации кода для функций, возвращающих простые типы int, string.
но xo генерирует структуры например с одним полем - например ID::int, 

Пока выхожу из положения так:

добавляю к типу приставку Simple, а в запросе возвращаю AS value
в результате в GO генерируется такой код
```go

type SimpleBool struct {
    Value bool
}

type SimpleInt struct {
    Value int
}

```

### Двойное и более, включение одинакового параметра запроса

Если в SQL нужно 2 и более раз включить какое-то одно значение
например
```SQL
SELECT

WHERE
    first_name ILIKE $1
OR  last_name ILIKE $1
```

xo сгененрирует функцию с N одинаковыми входящими параметрами, что вызовет ошибку компиляции GO.

для решения этой ситуации нужно модифицировать SQL следующим образом

```SQL
-- 1. Объявить секцию WITH
WITH
    sqlvars
        (name, age) AS
--      (VALUES ($1::varchar, $2::int))        
        (VALUES (%%name string%%::varchar, %%age int%%::int))

SELECT

FROM
    a_table as a

-- 2. Приджоинить sqlvars объявленные в заголовке WITH
    CROSS JOIN sqlvars AS sv

WHERE
-- 3. обращаться к sv.*
    a.first_name ILIKE sv.name
OR  a.last_name ILIKE sv.name

AND a.age = sv.age
```

### Массивы в качестве параметров запроса

В SQL указываем тип interface{} для входящего параметра

```sql
WITH
    sqlvars
        (taskTypes) AS
        (VALUES (%%taskTypes interface{}%%::int[]))
```

при вызове функции параметр оборачиваем в pq.Array
```go

taskTypes := []int{1,2,3}

data, err := SalesPlanByYearByMonth(db, ..., pq.Array(taskTypes))

```

для pgx оборачивать в pq.Array не нужно
можно передать массив напрямую

### Custom Update/Insert/Delete

https://github.com/xo/xo/issues/169
к сожалению xo не умеет генерировать код на кастомные Update/Insert/Delete SQL.

для выхода из этой ситуации я написал программу xouid
программа получает на вход все те же параметры что и генератор xo
по шаблону из sql файлов генерирует код
sql запрос валидируется при помощи конструкции EXPLAIN <SQL>.
