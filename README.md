<p align="center">
<img alt="CloverDB Logo" src=".github/logo.png#gh-light-mode-only" width="300px">
<img alt="CloverDB Logo" src=".github/logo-white.png#gh-dark-mode-only" width="300px">
</p>
<h2 align="center">Lightweight document-oriented NoSQL Database</h2>

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  
[![Go Reference](https://pkg.go.dev/badge/badge/github.com/xfp-881643/clover.svg)](https://pkg.go.dev/github.com/xfp-881643/clover)
[![Go Report Card](https://goreportcard.com/badge/github.com/xfp-881643/clover)](https://goreportcard.com/report/github.com/xfp-881643/clover)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Coverage Status](https://coveralls.io/repos/github/ostafen/clover/badge.svg)](https://coveralls.io/github/ostafen/clover)
[![Join the chat at https://gitter.im/cloverDB/community](https://badges.gitter.im/cloverDB/community.svg)](https://gitter.im/cloverDB/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

> 🇬🇧 English | [🇨🇳 简体中文](README-CN.md) | [🇪🇸 Spanish](README-ES.md)

**CloverDB** is a lightweight NoSQL database designed for being simple and easily maintainable, thanks to its small code base. It has been inspired by [tinyDB](https://github.com/msiemens/tinydb).

## Features

- Document oriented
- Written in pure Golang
- Simple and intuitive api
- Easily maintainable

## Why CloverDB?

**CloverDB** has been written for being easily maintenable. As such, it trades performance with simplicity, and is not intented to be an alternative to more performant databases such as **MongoDB** or **MySQL**.
However, there are projects where running a separate database server may result overkilled, and, for simple queries, network delay may be the major performance bottleneck.
For there scenario, **CloverDB** may be a more suitable alternative.

## Database Layout

**CloverDB** abstracts the way collections are stored on disk through the **StorageEngine** interface. The default implementation is based on the [Badger](https://github.com/dgraph-io/badger) database key-value store. However, you could easily write your own storage engine implementation.

## Installation
Make sure you have a working Go environment (Go 1.13 or higher is required). 
```shell
  GO111MODULE=on go get github.com/xfp-881643/clover
```

## Databases and Collections

CloverDB stores data records as JSON documents, which are grouped together in collections. A database is made up of one or more collections.

### Database

To store documents inside collections, you have to open a Clover database using the `Open()` function.

```go
import (
	"log"
	c "github.com/xfp-881643/clover"
)

...

db, _ := c.Open("clover-db")

// or, if you don't need persistency
db, _ := c.Open("", c.InMemoryMode(true))

defer db.Close() // remember to close the db when you have done
```

### Collections

CloverDB stores documents inside collections. Collections are the **schemaless** equivalent of tables in relational databases. A collection is created by calling the `CreateCollection()` function on a database instance. New documents can be inserted using the `Insert()` or `InsertOne()` methods. Each document is uniquely identified by a **Version 4 UUID** stored in the **_id** special field and generated during insertion.

```go
db, _ := c.Open("clover-db")
db.CreateCollection("myCollection") // create a new collection named "myCollection"

// insert a new document inside the collection
doc := c.NewDocument()
doc.Set("hello", "clover!")

// InsertOne returns the id of the inserted document
docId, _ := db.InsertOne("myCollection", doc)
fmt.Println(docId)
```

### Importing and Exporting Collections

CloverDB is capable of easily importing and exporting collections to JSON format regardless of the storage engine used.

```go
// dump the content of the "todos" collection in a "todos.json" file
db.ExportCollection("todos", "todos.json")

...

// recover the todos collection from the exported json file
db.DropCollection("todos")
db.ImportCollection("todos", "todos.json")

docs, _ := db.Query("todos").FindAll()
for _, doc := range docs {
  log.Println(doc)
}
```

## Queries

CloverDB is equipped with a fluent and elegant API to query your data. A query is represented by the **Query** object, which allows to retrieve documents matching a given **criterion**. A query can be created by passing a valid collection name to the `Query()` method.

### Select All Documents in a Collection

The `FindAll()` method is used to retrieve all documents satisfying a given query.

```go
docs, _ := db.Query("myCollection").FindAll()

todo := &struct {
    Completed bool   `clover:"completed"`
    Title     string `clover:"title"`
    UserId    int    `clover:"userId"`
}{}

for _, doc := range docs {
    doc.Unmarshal(todo)
    log.Println(todo)
}
```

### Filter Documents with Criteria

In order to filter the documents returned by `FindAll()`, you have to specify a query Criteria using the `Where()` method. A Criteria object simply represents a predicate on a document, evaluating to **true** only if the document satisfies all the query conditions. 


The following example shows how to build a simple Criteria, matching all the documents having the **completed** field equal to true.

```go
db.Query("todos").Where(c.Field("completed").Eq(true)).FindAll()

// or equivalently
db.Query("todos").Where(c.Field("completed").IsTrue()).FindAll()
```

In order to build very complex queries, we chain multiple Criteria objects by using the `And()` and `Or()` methods, each returning a new Criteria obtained by appling the corresponding logical operator.

```go
// find all completed todos belonging to users with id 5 and 8
db.Query("todos").Where(c.Field("completed").Eq(true).And(c.Field("userId").In(5, 8))).FindAll()
```

Naturally, you can also create Criteria involving multiple fields. CloverDB provides you with two equivalent ways to accomplish this:

```go
db.Query("myCollection").Where(c.Field("myField1").Gt(c.Field("myField2")))
// or, if you prefer
db.Query("myCollection").Where(c.Field("myField1").Gt("$myField2"))
```

### Sorting Documents

To sort documents in CloverDB, you need to use `Sort()`. It is a variadic function which accepts a sequence of SortOption, each allowing to specify a field and a sorting direction.
A sorting direction can be one of 1 or -1, respectively corresponding to ascending and descending order. If no SortOption is provided, `Sort()` uses the **_id** field by default.

```go
// Find any todo belonging to the most recent inserted user
db.Query("todos").Sort(c.SortOption{"userId", -1}).FindFirst()
```

### Skip/Limit Documents

Sometimes, it can be useful to discard some documents from the output, or simply set a limit on the maximum number of results returned by a query. For this purpose, CloverDB provides the `Skip()` and `Limit()` functions, both accepting an interger $n$ as parameter.

```go
// discard the first 10 documents from the output,
// also limiting the maximum number of query results to 100
db.Query("todos").Skip(10).Limit(100).FindAll()
```

### Update/Delete Documents

The `Update()` method is used to modify specific fields of documents in a collection. The `Delete()` method is used to delete documents. Both methods belong to the Query object, so that it is easy to update and delete documents matching a particular query.

```go
// mark all todos belonging to user with id 1 as completed
updates := make(map[string]interface{})
updates["completed"] = true

db.Query("todos").Where(c.Field("userId").Eq(1)).Update(updates)

// delete all todos belonging to users with id 5 and 8
db.Query("todos").Where(c.Field("userId").In(5,8)).Delete()
```

To update or delete a single document using a specific document id, use `UpdateById()` or `DeleteById()`, respectively:

```go
docId := "1dbce353-d3c6-43b3-b5a8-80d8d876389b"
// update the document with the specified id
db.Query("todos").UpdateById(docId, map[string]interface{}{"completed": true})
// or delete it
db.Query("todos").DeleteById(docId)
```

## Data Types

Internally, CloverDB supports the following primitive data types: **int64**, **uint64**, **float64**, **string**, **bool** and **time.Time**. When possible, values having different types are silently converted to one of the internal types: signed integer values get converted to int64, while unsigned ones to uint64. Float32 values are extended to float64.

For example, consider the following snippet, which sets an uint8 value on a given document field:

```go
doc := c.NewDocument()
doc.Set("myField", uint8(10)) // "myField" is automatically promoted to uint64

fmt.Println(doc.Get("myField").(uint64))
```

Pointer values are dereferenced until either **nil** or a **non-pointer** value is found:

``` go
var x int = 10
var ptr *int = &x
var ptr1 **int = &ptr

doc.Set("ptr", ptr)
doc.Set("ptr1", ptr1)

fmt.Println(doc.Get("ptr").(int64) == 10)
fmt.Println(doc.Get("ptr1").(int64) == 10)

ptr = nil

doc.Set("ptr1", ptr1)
// ptr1 is not nil, but it points to the nil "ptr" pointer, so the field is set to nil
fmt.Println(doc.Get("ptr1") == nil)
```

Invalid types leaves the document untouched:

```go
doc := c.NewDocument()
doc.Set("myField", make(chan struct{}))

log.Println(doc.Has("myField")) // will output false
```

## Contributing

**CloverDB** is actively developed. Any contribution, in the form of a suggestion, bug report or pull request, is well accepted :blush:

Major contributions and suggestions have been gratefully received from (in alphabetical order):

- [ASWLaunchs](https://github.com/ASWLaunchs)
- [jsgm](https://github.com/jsgm)
- [segfault99](https://github.com/segfault99)