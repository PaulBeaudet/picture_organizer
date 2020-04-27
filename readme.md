# Picture Organizer

The goal of this program was just solving the simple problem of moving my camera uploads folder to and a better organized picture folder

That is organized in a top level photos folder by
  * Year (eg 2020)
  * Then month and day (eg 04_27_)
    * Can later be appended with event name
  * Then pictures are renamed with the time (hh_mm_ss.jpg or hh_mm_ss.rw)

### Usage

```
go run organize.go -src "/home/username/path/with/trailing/" -dest "/same/diff/folder/"  
```

Trailing slash and exact path are important until someone fixes it to be more flexible
