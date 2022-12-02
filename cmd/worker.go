
package main

import (
	"github.com/skyhookml/skyhookml/skyhook"
	gouuid "github.com/google/uuid"

	_ "github.com/skyhookml/skyhookml/ops"

	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: ./worker [external IP] [port]")
		fmt.Println("example: ./worker localhost 8081")
		return
	}
	myIP := os.Args[1]
	myPort := skyhook.ParseInt(os.Args[2])

	mode := "docker"
	if len(os.Args) >= 4 {
		mode = os.Args[3]
	}

	workingDir, err := os.Getwd()