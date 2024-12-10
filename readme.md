# server-core 💫
The module allows you to run several separate modules.
For example, it can be a gin server and some kind of service job (it can be anything). 
`server-core` provides the ability to communicate between different modules
via the event bus, and also allows you to build a dependency (requirements) tree 
(for example, a separate database module that will be used in different modules)

## Import
This module can be imported into a go project using the command

```bash 
go get github.com/AlekseiKromski/server-core
```
## Example of usage in your project

```go
package main

import (
	"github.com/AlekseiKromski/server-core/core"
)

func main() {
	c := core.NewCore()
	c.Init([]core.Module{
		// Your modules
	})
}

```
