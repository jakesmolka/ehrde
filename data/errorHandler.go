/* EHRDE's error functions
v1: error handler with log output
*/

package data

import (
	"log"
)

//general error handler to print error and optional more details
func errorHandler(err error, details string) {
	if err != nil {
		//log
		log.Print("ERROR: ", err)
		if details != "" {
			log.Print("details: ", details)
		}
	}
}
