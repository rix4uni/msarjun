package banner

import (
	"fmt"
)

// prints the version message
const version = "v0.0.2"

func PrintVersion() {
	fmt.Printf("Current msarjun version %s\n", version)
}

// Prints the Colorful banner
func PrintBanner() {
	banner := `
                                  _             
   ____ ___   _____ ____ _ _____ (_)__  __ ____ 
  / __  __ \ / ___// __  // ___// // / / // __ \
 / / / / / /(__  )/ /_/ // /   / // /_/ // / / /
/_/ /_/ /_//____/ \__,_//_/ __/ / \__,_//_/ /_/ 
                           /___/
`
	fmt.Printf("%s\n%60s\n\n", banner, "Current msarjun version "+version)
}
