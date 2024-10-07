package banner

import (
	"fmt"
)

// prints the version message
const version = "v0.0.1"

func PrintVersion() {
	fmt.Printf("Current emailfinder version %s\n", version)
}

// Prints the Colorful banner
func PrintBanner() {
	banner := `
                           _  __ ____ _             __           
  ___   ____ ___   ____ _ (_)/ // __/(_)____   ____/ /___   _____
 / _ \ / __  __ \ / __  // // // /_ / // __ \ / __  // _ \ / ___/
/  __// / / / / // /_/ // // // __// // / / // /_/ //  __// /    
\___//_/ /_/ /_/ \__,_//_//_//_/  /_//_/ /_/ \__,_/ \___//_/ 
`
	fmt.Printf("%s\n%75s\n\n", banner, "Current emailfinder version "+version)
}
