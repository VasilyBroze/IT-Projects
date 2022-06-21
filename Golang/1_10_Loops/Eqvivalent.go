package main
import "fmt"
//import "strconv"
func main() {
    var a, b, otv string
    //var as, bs string
    fmt.Scan(&a, &b)
    //564 8954
    for i:=0; i<len(a); i++ {
        for n:=0; n<len(b); n++ {
            if a[i] == b[n] {
                if otv == "" {
                    otv += string(a[i])
                } else {
                    otv += " " + string(a[i])
                }
                break
            }
        }
    }
   fmt.Println(otv)

}




