## TinyORM

简易ORM框架，目前仅实现了基于sqlite3的数据库支持

```shell
go get -u github.com/CenGinHong/TinyORM
```



## CRUD

### 插入

```go
import (
	"fmt"
	"github.com/CenGinHong/TinyORM"
)

type User struct {
	Name string `tinyorm:"PRIMARY KEY"`
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func main() {
	engine, err := tinyORM.NewEngine("sqlite3", "tiny.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Close()
	session := engine.NewSession()
	err = session.Model(&User{}).CreateTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = session.Model(&User{}).Insert(user1, user2)
	if err != nil {
		fmt.Println(err)
		return
	}
}

```



### 查询

```go
import (
	"fmt"
	"github.com/CenGinHong/TinyORM"
)

type User struct {
	Name string `tinyorm:"PRIMARY KEY"`
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
	user4 = &User{"John", 32}
)

func main() {
	engine, err := tinyORM.NewEngine("sqlite3", "tiny.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	err = s.DropTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.CreateTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = s.Insert(user1, user2, user3, user4)
	if err != nil {
		fmt.Println(err)
		return
	}
	u := &User{}
	// 仅查询一条数据
	err = s.Where("Age > ?", 20).First(u)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u)
	// 条件查询
	var users []User
	if err = s.Where("Age > ?", 20).OrderBy("Age").Limit(2).Find(&users); err != nil {
		fmt.Println(err)
		return
	}
	for _, user := range users {
		fmt.Println(user)
	}
}
```



### 更新

```go
import (
	"fmt"
	"github.com/CenGinHong/TinyORM"
)

type User struct {
	Name string `tinyORM:"PRIMARY KEY"`
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func main() {
	engine, err := tinyORM.NewEngine("sqlite3", "tiny.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	err = s.DropTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.CreateTable()
	if err != nil {

		return
	}
	_, err = s.Insert(user1)
	if err != nil {
		fmt.Println(err)
		return
	}
	users := make([]User, 0)
	if err = s.Find(&users); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("after insert:")
	for _, user := range users {
		fmt.Println(user)
	}
	_, err = s.Where("Name = ?", "Tom").Update(map[string]interface{}{
		"Age": 30,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	users = make([]User, 0)
	if err = s.Find(&users); err != nil {
		fmt.Println(err)
		return
	}
    fmt.Println("after update:")
	for _, user := range users {
		fmt.Println(user)
	}
}
```



### 删除

```go
import (
	"fmt"
	"github.com/CenGinHong/TinyORM"
)

type User struct {
	Name string `tinyORM:"PRIMARY KEY"`
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func main() {
	engine, err := tinyORM.NewEngine("sqlite3", "tiny.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	err = s.DropTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.CreateTable()
	if err != nil {

		return
	}
	_, err = s.Insert(user1, user2)
	if err != nil {
		fmt.Println(err)
		return
	}
	users := make([]User, 0)
	if err = s.Find(&users); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("after insert:")
	for _, user := range users {
		fmt.Println(user)
	}
	_, err = s.Where("Name = ?", "Tom").Delete()
	if err != nil {
		fmt.Println(err)
		return
	}
	users = make([]User, 0)
	if err = s.Find(&users); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("after delete:")
	for _, user := range users {
		fmt.Println(user)
	}
}
```



### 计数

```go
import (
	"fmt"
	"github.com/CenGinHong/TinyORM"
)

type User struct {
	Name string `tinyORM:"PRIMARY KEY"`
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func main() {
	engine, err := tinyORM.NewEngine("sqlite3", "tiny.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	err = s.DropTable()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.CreateTable()
	if err != nil {

		return
	}
	_, err = s.Insert(user1, user2)
	if err != nil {
		fmt.Println(err)
		return
	}
	count, err := s.Count()
	if err != nil {
		return
	}
	fmt.Println(count);
}
```



## HOOK

以表结构作为`receiver`添加相对应的方法即可进行相应的回调.

目前支持的Hook埋点如下

```
BeforeQuery 
AfterQuery  
BeforeUpdate 
AfterUpdate 
BeforeDelete
AfterDelete 
BeforeInsert
AfterInsert
```



```go
import (
	"fmt"
	"github.com/CenGinHong/TinyORM"
	"github.com/CenGinHong/TinyORM/log"
	"github.com/CenGinHong/TinyORM/session"
)

type Account struct {
	ID       int `tinyorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *session.Session) error {
	log.Info("before inert", account)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *session.Session) error {
	log.Info("after query", account)
	account.Password = "******"
	return nil
}

func main() {
	engine, err := tinyORM.NewEngine("sqlite3", "tiny.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Close()
	s := engine.NewSession().Model(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	a := &Account{1, "123456"}
	fmt.Println("before insert")
	fmt.Println(a)
	_, _ = s.Insert(a)
	u := &Account{}
	err = s.First(u)
	fmt.Println("after insert")
	fmt.Println(u)
}
```



## 事务

使用`func (e *Engine) Transaction(f TxFunc) (result interface{}, err error)`进行事务的操作，用户将所有的操作放入回调函数`TxFunc`在，如果在该回调函数中没有返回err，则会被commit，否则将会rollback

```go
import (
	"errors"
	"fmt"
	"github.com/CenGinHong/TinyORM"
	"github.com/CenGinHong/TinyORM/session"
)

type User struct {
	Name string `tinyORM:"PRIMARY KEY"`
	Age  int
}

func main() {
	engine, err := tinyORM.NewEngine("sqlite3", "tiny.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Close()
	_ = engine.NewSession().Model(&User{}).DropTable()
	_, err = engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		if err = s.Model(&User{}).CreateTable(); err != nil {
			return nil, err
		}
		result, err = s.Insert(&User{"Tom", 18})
		return nil, err
	})
}
```

