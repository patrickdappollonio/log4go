package main

import l4g "github.com/ccpaging/log4go"

func main() {
	// Load the configuration (isn't this easy?)
	l4g.LoadConfiguration("config.xml")

	// And now we're ready!
	l4g.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	l4g.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
	l4g.Info("About that time, eh chaps?")

	l4g.Close()

	l4g.LoadConfiguration("config.json")

	// And now we're ready!
	l4g.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	l4g.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
	l4g.Info("About that time, eh chaps?")

	l4g.Close()
}

