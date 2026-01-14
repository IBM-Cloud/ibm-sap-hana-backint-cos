/*
Contains all global datatypes and variables
*/
package logging

import "github.com/sirupsen/logrus"

// Datatype representing the logging text formatter for hdbbackint
type backintFormatter struct {
	logrus.TextFormatter
}

// Datatype representing the backint result messages
type BackintResultMessages []string
