# **GORM** 入门指南

gorm 倾向于约定，而不是配置。默认情况下，gorm会使用`ID`作为主键，使用结构体名的蛇形复数（如：`OrderRecord `-> `order_records`）作为表名，字段名的蛇形作为列名。使用`CreateAT`、`UpdateAt`字段追踪创建、更新时间。

## 1 gorm.Model

gorm定义一个gorm.Model结构体，包括字段ID、CreatedAt、UpdateAt、DeleteAt

```go
type Model struct{
    ID            uint           `gorm:"primayKey"`
    CreatedAt     time.Time
    UpdateAt      time.Time
    DeletedAt     gorm.DeletedAt `gorm:"index"`    
}
```

### 字段标签

声明 model 时，tag 是可选的，GORM 支持以下 tag： tag 名大小写不敏感，但建议使用 `camelCase` 风格

|         标签名         | 说明                                                         |
| :--------------------: | :----------------------------------------------------------- |
|         column         | 指定 db 列名                                                 |
|          type          | 列数据类型，推荐使用兼容性好的通用类型，例如：所有数据库都支持 bool、int、uint、float、string、time、bytes 并且可以和其他标签一起使用，例如：`not null`、`size`, `autoIncrement`… 像 `varbinary(8)` 这样指定数据库数据类型也是支持的。在使用指定数据库数据类型时，它需要是完整的数据库数据类型，如：`MEDIUMINT UNSIGNED not NULL AUTO_INCREMENT` |
|          size          | 指定列大小，例如：`size:256`                                 |
|       primaryKey       | 指定列为主键                                                 |
|         unique         | 指定列为唯一                                                 |
|        default         | 指定列的默认值                                               |
|       precision        | 指定列的精度                                                 |
|         scale          | 指定列大小                                                   |
|        not null        | 指定列为 NOT NULL                                            |
|     autoIncrement      | 指定列为自动增长                                             |
| autoIncrementIncrement | 自动步长，控制连续记录之间的间隔                             |
|      **embedded**      | **嵌套字段**                                                 |
|   **embeddedPrefix**   | **嵌入字段的列名前缀**                                       |
|     autoCreateTime     | 创建时追踪当前时间，对于 `int` 字段，它会追踪秒级时间戳，您可以使用 `nano`/`milli` 来追踪纳秒、毫秒时间戳，例如：`autoCreateTime:nano` |
|     autoUpdateTime     | 创建/更新时追踪当前时间，对于 `int` 字段，它会追踪秒级时间戳，您可以使用 `nano`/`milli` 来追踪纳秒、毫秒时间戳，例如：`autoUpdateTime:milli` |
|         index          | 根据参数创建索引，多个字段使用相同的名称则创建复合索引，查看 [索引](https://gorm.io/zh_CN/docs/indexes.html) 获取详情 |
|      uniqueIndex       | 与 `index` 相同，但创建的是唯一索引                          |
|         check          | 创建检查约束，例如 `check:age > 13`，查看 [约束](https://gorm.io/zh_CN/docs/constraints.html) 获取详情 |
|         **<-**         | **设置字段写入的权限， `<-:create` 只创建、`<-:update` 只更新、`<-:false` 无写入权限、`<-` 创建和更新权限** |
|         **->**         | **设置字段读的权限，`->:false` 无读权限**                    |
|           -            | ignore this field, `-` no read/write permission, `-:migration` no migrate permission, `-:all` no read/write/migrate permission |
|        comment         | 迁移时为字段添加注释                                         |

### 嵌入结构体

对于匿名字段，GORM 会将其字段包含在父结构体中，例如：

```go
type User struct {
  gorm.Model
  Name string
}
// 等效于
type User struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
  Name string
}
```

对于正常的结构体字段，你也可以通过标签 `embedded` 将其嵌入，例如：

```go
type Author struct {
    Name  string
    Email string
}

type Blog struct {
  ID      int
  Author  Author `gorm:"embedded"`
  Upvotes int32
}
// 等效于
type Blog struct {
  ID    int64
  Name  string
  Email string
  Upvotes  int32
}
```

并且，您可以使用标签 `embeddedPrefix` 来为 db 中的字段名添加前缀，例如：

```go
type Blog struct {
  ID      int
  Author  Author `gorm:"embedded;embeddedPrefix:author_"`
  Upvotes int32
}
// 等效于
type Blog struct {
  ID          int64
    AuthorName  string
    AuthorEmail string
  Upvotes     int32
}
```

## 2 链接数据库

### mysql

```go
import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func main() {
  // 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
  dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
```



GORM 允许通过一个现有的数据库连接来初始化 `*gorm.DB`

```go
import (
  "database/sql"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)
// 现有的数据库连接
sqlDB, err := sql.Open("mysql", "mydb_dsn")
// 使用现有连接初始化gorm
gormDB, err := gorm.Open(mysql.New(mysql.Config{
  Conn: sqlDB,
}), &gorm.Config{})
```



**GORM 使用 [database/sql](https://pkg.go.dev/database/sql) 维护连接池**

```go
sqlDB, err := db.DB()

// SetMaxIdleConns 设置空闲连接池中连接的最大数量
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns 设置打开数据库连接的最大数量。
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime 设置了连接可复用的最大时间。
sqlDB.SetConnMaxLifetime(time.Hour)
```



## 3 CURD

###  1 创建

####1.1**创建记录**

```go
user := User{Name:"Jinzhu", Age:18, Birthday:time.Now()}

result := db.Create(&user)// 需要通过数据的指针来创建（需要传入引用类型之间操作原值，可以减少消耗）

user.ID					 // 返回插入数据的主键
result.Error             // 返回Error
result.RowAffected		 // 返回最大插入数目

```

> **注意** 使用`CreateBatchSize` 选项初始化 GORM 时，所有的创建& 关联 `INSERT` 都将遵循该选项

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
  CreateBatchSize: 1000,
})
```

#### 1.2 部分更新

**创建记录并更新给出的字段。**

```go
db.Select("Name", "Age", "CreatedAt").Create(&user)
```

> 等价于
>
> ```sql
> //  INSERT INTO `users` (`name`,`age`,`created_at`) VALUES ("jinzhu", 18, "2020-07-04 11:05:21.775")
> ```

**创建一个记录且一同忽略传递给略去的字段值。**

```go
db.Omit("Name", "Age", "CreatedAt").Create(&user)
```

> 等价于
>
> ```sql
> // INSERT INTO `users` (`birthday`,`updated_at`) VALUES ("2020-01-01 00:00:00.000", "2020-07-04 11:05:21.775")
> ```

#### 1.3 批量插入

将有一个slice切片传递给create方法。gorm会生成一条语句插入所有的数据，**并回填主键的值（添加完数据之后，需要获取刚刚添加的数据 id），调用钩子方法（通过一个方法来干涉另一个方法的行为，就是一个方法的返回结果或修改的变量是另一个方法执行时的if判断条件或for/while循环调用条件的）**

```go
var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
db.Create(&users)
for _, user := range users {  
	user.ID // 1,2,3
}
```

#### 1.4 创建钩子

GORM 允许用户定义的钩子有 **`BeforeSave`, `BeforeCreate`, `AfterSave`, `AfterCreate`** 创建记录时将调用这些钩子方法

**钩子方法（Hook）**是在创建、查询、更新、删除等操作之前，之后调用的函数，如果已经为模型定义了指定的方法，它会在创建、更新、查询、删除时自动被调用。如果任何回调但会错误，那么gorm会停止操作，并回滚事务

**钩子方法的函数签名应该是 `func(*gorm.DB) error`**

示例：

```go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {  
	u.UUID = uuid.New()  
    // 这种是自定义的根据自己的需求
	if !u.IsValid() {    
		err = errors.New("can't save invalid data")  
	}  
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {  
	if u.ID == 1 {    
		tx.Model(u).Update("role", "admin")  
	}  
	return
}
```

#### 1.5 根据Map创建

GORM 支持根据 `map[string]interface{}` 和 `[]map[string]interface{}{}` 创建记录

```go
db.Model(&User{}).Create(map[string]interface{}{  
	"Name": "jinzhu", "Age": 18,
})

// batch insert from `[]map[string]interface{}{}`
db.Model(&User{}).Create([]map[string]interface{}{  
    {"Name": "jinzhu_1", "Age": 18},  
    {"Name": "jinzhu_2", "Age": 20},
})
```

#### 1.6 默认值

通过标签的`default`为字段定义默认值，在插入记录到数据库的时候，**默认值会被用于填充值是0值的字段**

```go
type User struct {
  ID   int64
  Name string `gorm:"default:galeone"`
  Age  int64  `gorm:"default:18"`
}
```

**注意：** 对于声明了默认值的字段，像 `0`、`''`、`false` 等零值是不会保存到数据库。您需要使用指针类型或 Scanner/Valuer 来避免这个问题，例如：

```go
type User struct {  
    gorm.Model  
    Name string  
    Age  *int           `gorm:"default:18"`  
    Active sql.NullBool `gorm:"default:true"`
}
```

#### 1.7 Upset

判断应该是更新还是插入

- 通过判断插入的记录里是否存在主键索引或唯一索引冲突，来决定是插入还是更新。
- 当出现主键索引或唯一索引冲突时则进行update操作，否则进行insert操作

```go
import "gorm.io/gorm/clause"

// 1. 在冲突时，什么都不做
db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)

// 2. 在`id`冲突时，将列更新为默认值
db.Clauses(clause.OnConflict{
  Columns:   []clause.Column{{Name: "id"}},
  DoUpdates: clause.Assignments(map[string]interface{}{"role": "user"}),
}).Create(&users)
// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET ***; SQL Server
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE ***; MySQL

// 3. 使用SQL语句
db.Clauses(clause.OnConflict{
  Columns:   []clause.Column{{Name: "id"}},
  DoUpdates: clause.Assignments(map[string]interface{}{"count": gorm.Expr("GREATEST(count, VALUES(count))")}),
}).Create(&users)
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE `count`=GREATEST(count, VALUES(count));

// 4. 在`id`冲突时，将列更新为新值
db.Clauses(clause.OnConflict{
  Columns:   []clause.Column{{Name: "id"}},
  DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
}).Create(&users)
// MERGE INTO "users" USING *** WHEN NOT MATCHED THEN INSERT *** WHEN MATCHED THEN UPDATE SET "name"="excluded"."name"; SQL Server
// INSERT INTO "users" *** ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name", "age"="excluded"."age"; PostgreSQL
// INSERT INTO `users` *** ON DUPLICATE KEY UPDATE `name`=VALUES(name),`age=VALUES(age); MySQL

// 5. 在冲突时，更新除主键以外的所有列到新值。
db.Clauses(clause.OnConflict{
  	UpdateAll: true,
}).Create(&users)
// INSERT INTO "users" *** ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name", "age"="excluded"."age", ...;
```



### 2 查询

1. gorm提供了`First`、`Take`、`Last`方法，以便从数据库中**检索单个对象**，相当于sql中添加了`limit 1`条件

```go
// 获取第一条记录（主键升序）
db.First(&user)
// SELECT * FROM users ORDER BY id LIMIT 1;

// 获取一条记录，没有指定排序字段
db.Take(&user)
// SELECT * FROM users LIMIT 1;

// 获取最后一条记录（主键降序）
db.Last(&user)
// SELECT * FROM users ORDER BY id DESC LIMIT 1;

result := db.First(&user)
result.RowsAffected // 返回找到的记录数
result.Error        // returns error or nil

// 检查 ErrRecordNotFound 错误
errors.Is(result.Error, gorm.ErrRecordNotFound)
```

2. First与Last 会根据主键排序**（若数据表model没有定义主键，会按照其第一个字段进行排序）**，分别查询第一条和最后一条记录。**只有在目标struct是指针或者通过 db.Model()指定model时，该方法才有效。**

3. 如果主键是数字类型，可以通过内联条件进行检索。

 **内联条件：**查询条件也可以被内联到 `First` 和 `Find` 之类的方法中，其用法类似于 `Where`。

```go
// 根据主键获取记录，如果是非整型主键
db.First(&user, "id = ?", "string_primary_key")
// SELECT * FROM users WHERE id = 'string_primary_key';

// Plain SQL
db.Find(&user, "name = ?", "jinzhu")
// SELECT * FROM users WHERE name = "jinzhu";

db.Find(&users, "name <> ? AND age > ?", "jinzhu", 20)
// SELECT * FROM users WHERE name <> "jinzhu" AND age > 20;

// Struct
db.Find(&users, User{Age: 20})
// SELECT * FROM users WHERE age = 20;

// Map
db.Find(&users, map[string]interface{}{"age": 20})
// SELECT * FROM users WHERE age = 20;
```

#### 4 获取全部记录

```go
result := db.Find(&users)// SELECT * FROM users

result.RowsAffected // 返回找到的记录数，相当于 `len(users)`
result.Error        // returns error
```

#### 5 根据where条件获取

```go
// 获取第一条匹配的记录
db.Where("name = ?", "jinzhu").First(&user)
// SELECT * FROM users WHERE name = 'jinzhu' ORDER BY id LIMIT 1;

// 获取全部匹配的记录
db.Where("name <> ?", "jinzhu").Find(&users)
// SELECT * FROM users WHERE name <> 'jinzhu';

// IN
db.Where("name IN ?", []string{"jinzhu", "jinzhu 2"}).Find(&users)
// SELECT * FROM users WHERE name IN ('jinzhu','jinzhu 2');

// LIKE
db.Where("name LIKE ?", "%jin%").Find(&users)
// SELECT * FROM users WHERE name LIKE '%jin%';

// AND
db.Where("name = ? AND age >= ?", "jinzhu", "22").Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' AND age >= 22;

// Time
db.Where("updated_at > ?", lastWeek).Find(&users)
// SELECT * FROM users WHERE updated_at > '2000-01-01 00:00:00';

// BETWEEN
db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&users)
// SELECT * FROM users WHERE created_at BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';

```

#### 6 Struct & Map 条件

```go
// Struct:需要传入一个已经初始化了的模型（由于其他字段为默认值是0值，将不被查询数据）
db.Where(&User{Name: "jinzhu", Age: 20}).First(&user)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20 ORDER BY id LIMIT 1;

// Map:直接传入一个map，不像struct一样，使用map遇到0值时候，仍然会被查询包含所有 key-value 的查询条件
db.Where(map[string]interface{}{"name": "jinzhu", "age": 20}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 20;

// 主键切片条件
db.Where([]int64{20, 21, 22}).Find(&users)
// SELECT * FROM users WHERE id IN (20, 21, 22);
```

**注：**当使用 struct 进行查询时，可以通过向 `Where()` 传入 struct 来指定查询条件的字段、值、表名，例如：

```go
db.Where(&User{Name: "jinzhu"}, "name", "Age").Find(&users)
// SELECT * FROM users WHERE name = "jinzhu" AND age = 0;

db.Where(&User{Name: "jinzhu"}, "Age").Find(&users)
// SELECT * FROM users WHERE age = 0;
```

#### 7 根据not条件获取

构建 NOT 条件，用法与 `Where` 类似

```go
db.Not("name = ?", "jinzhu").First(&user)
// SELECT * FROM users WHERE NOT name = "jinzhu" ORDER BY id LIMIT 1;

// Not In
db.Not(map[string]interface{}{"name": []string{"jinzhu", "jinzhu 2"}}).Find(&users)
// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");

// Struct
db.Not(User{Name: "jinzhu", Age: 18}).First(&user)
// SELECT * FROM users WHERE name <> "jinzhu" AND age <> 18 ORDER BY id LIMIT 1;

// 不在主键切片中的记录
db.Not([]int64{1,2,3}).First(&user)
// SELECT * FROM users WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;
```

#### 8 or条件的使用

```go
db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';

// Struct
db.Where("name = 'jinzhu'").Or(User{Name: "jinzhu 2", Age: 18}).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);

// Map
db.Where("name = 'jinzhu'").Or(map[string]interface{}{"name": "jinzhu 2", "age": 18}).Find(&users)
// SELECT * FROM users WHERE name = 'jinzhu' OR (name = 'jinzhu 2' AND age = 18);
```

#### 9 使用select检索指定字段

```go
db.Select("name", "age").Find(&users)
// SELECT name, age FROM users;

db.Select([]string{"name", "age"}).Find(&users)
// SELECT name, age FROM users;

db.Table("users").Select("COALESCE(age,?)", 42).Rows()
// SELECT COALESCE(age,'42') FROM users;
```

#### 10 Order

**ORDER BY 语句用于根据指定的列对结果集进行排序**

1. asc 升序，可以省略，是数据库默认的排序方式

2. desc 降序，跟升序相反。

```go
db.Order("age desc, name").Find(&users)
// 首要按age的降序排列，其次按name进行排列
// SELECT * FROM users ORDER BY age desc, name;

// 多个 order
// 首要按age的降序排列，其次按name进行排列
db.Order("age desc").Order("name").Find(&users)
// SELECT * FROM users ORDER BY age desc, name;

db.Clauses(clause.OrderBy{
  Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{[]int{1, 2, 3}}, WithoutParentheses: true},
}).Find(&User{})
// ORDER BY FIELD 指自定义排序：按照字段id的1、2、3进行排序的
// field(value,str1,str2,str3,str4)，value与str1、str2、str3、str4比较，返回1、2、3、4，如遇到null或者不在列表中的数据则返回0.
// SELECT * FROM users ORDER BY FIELD(id,1,2,3)
```



#### 11 Limit&Offset

**`Limit` 指定获取记录的最大数量**

**`Offset` 指定在开始返回记录之前要跳过的记录数量**

```go
db.Limit(3).Find(&users)
// SELECT * FROM users LIMIT 3;

// 通过 -1 消除 Limit 条件
db.Limit(10).Find(&users1).Limit(-1).Find(&users2)
// SELECT * FROM users LIMIT 10; (users1)
// SELECT * FROM users; (users2)

db.Offset(3).Find(&users)
// SELECT * FROM users OFFSET 3;

db.Limit(10).Offset(5).Find(&users)
// SELECT * FROM users OFFSET 5 LIMIT 10;

// 通过 -1 消除 Offset 条件
db.Offset(10).Find(&users1).Offset(-1).Find(&users2)
// SELECT * FROM users OFFSET 10; (users1)
// SELECT * FROM users; (users2)

// select * from users LIMIT 3 OFFSET 1 表示跳过1条数据,从第2条数据开始取，取3条数据，也就是取2,3,4三条数据
```

#### 12.Group By & Having

**GROUP BY**语法可以根据**给定数据列的每个成员**对查询结果进行**分组统计**，最终得到一个**分组汇总表**

**HAVING语句**通常与**GROUP BY语句**联合使用，用来**过滤**由GROUP BY语句返回的**记录集**，弥补了WHERE关键字不能与聚合函数联合使用的不足。

```go
type result struct {
  Date  time.Time
  Total int
}

db.Model(&User{}).Select("name, sum(age) as total").Where("name LIKE ?", "group%").Group("name").First(&result)
// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "group%" GROUP BY `name` LIMIT 1


db.Model(&User{}).Select("name, sum(age) as total").Group("name").Having("name = ?", "group").Find(&result)
// SELECT name, sum(age) as total FROM `users` GROUP BY `name` HAVING name = "group"

rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Rows()
defer rows.Close()
for rows.Next() {
  ...
}

rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()
defer rows.Close()
for rows.Next() {
  ...
}

type Result struct {
  Date  time.Time
  Total int64
}
db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)
```

### 3 高级查询





### 4 更新

#### 1 保存所有字段

save 会保存所有的字段，即使字段是零值

```go
db.First(&user)
user.Name = "jinzhu 2"
user.Age = 100

db.save(&user)
// UPDATE users SET name='jinzhu 2', age=100, birthday='2016-01-01', updated_at = '2013-11-17 21:34:10' WHERE id=111;
```

#### 2 更新单个列

当使用update更新单个列的时候，需要指定条件，否则会返回错误。

当使用了Model方法，且该对象主键有值，该值会被用于构建条件

```go
// 条件更新
db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;

// User 的 ID 是 `111`
db.Model(&user).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

// 根据条件和 model 的值进行更新
db.Model(&user).Where("active = ?", true).Update("name", "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;

```

















































































































