package main

import "fmt"

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	c := NewConfig()
	l(SEVERITY_INFO, fmt.Sprintf("LibraryPath: %s", c.LibraryPath))
	l(SEVERITY_INFO, fmt.Sprintf("CurrentPath: %s", c.CurrentPath))
	c.startHttpListener()
}
