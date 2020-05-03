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

In mint I'm able to use ./local_build.sh in order to make the binary command line accessible as "organize". Feel free to rename it whatever suits you, maybe I'll make that a passible option to the script. Likely need to 

    source ~./profile
    
 To get the terminal to recongize the command, probably should just add that too. Not sure how much that makes sense to do every time though
