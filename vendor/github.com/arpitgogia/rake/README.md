Golang based implementation of Rake
---
RAKE is short for Rapid Automatic Keyword Extraction. The original research paper "Automatic keyword extraction from individual documents by Stuart Rose, Dave Engel, Nick Cramer and Wendy Cowley" can be found [here](https://www.researchgate.net/profile/Stuart_Rose/publication/227988510_Automatic_Keyword_Extraction_from_Individual_Documents/links/55071c570cf27e990e04c8bb.pdfs)

### Installation and Usage
* Install by executing `go get github.com/arpitgogia/rake`
* Use as shown below:
```
package main
import (
    "fmt"
    "github.com/arpitgogia/rake"
)

func main() {
    rake.WithText("Avengers: Infinity War")
    rake.WithFile("~/test.txt")
}
```

### Web API

Make a GET request on `https://frozen-lowlands-96920.herokuapp.com/rake?text=<text>`

### To Do

- [X] Basic implementation
- [X] Clean up and organize code
- [X] Implement package-like abstraction
- [X] Convert to a REST API
