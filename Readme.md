[![Go Report Card](https://goreportcard.com/badge/github.com/tartabit/tox)](https://goreportcard.com/report/github.com/tartabit/tox)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
# tox

tox is a simple library for rapid (and potentially lossy) type conversions.  It is very helpful when trying to do quick 
conversions of interfaces when confidence in the quality of the value is high.

## Installation

Checkout the repository and run:

```bash
go get github.com/tartabit/tox
```

## Usage

```golang
import "github.com/tartabit/tox"

var myint = 5
var mystring = tox.ToString(myint)
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
