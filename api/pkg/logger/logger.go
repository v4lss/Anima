// Package logger provides a thin structured logger wrapping the standard log package.
// Swap this for zerolog/zap in v2 without touching any other package.
package logger

import (
	"log"
	"os"
)

var std = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

func Info(msg string, args ...any)  { std.Printf("[INFO]  "+msg, args...) }
func Warn(msg string, args ...any)  { std.Printf("[WARN]  "+msg, args...) }
func Error(msg string, args ...any) { std.Printf("[ERROR] "+msg, args...) }
func Fatal(msg string, args ...any) { std.Fatalf("[FATAL] "+msg, args...) }
