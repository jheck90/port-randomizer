package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"time"
	"strconv"
	"strings"

	"port-randomizer/ports"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func isWellKnownPort(port int) bool {
	for _, wellKnownPort := range ports.WellKnownPorts {
		if port == wellKnownPort {
			return true
		}
	}
	return false
}

func checkWellKnownPort(cmd *cobra.Command, args []string) {
	randomIndex := rand.Intn(len(ports.WellKnownPorts))
	randomPort := ports.WellKnownPorts[randomIndex]
	fmt.Printf("Checking well-known port: %d\n", randomPort)
}

func isPortAvailable(port int, protocol string) bool {
	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen(protocol, address)
	if err != nil {
		return false
	}
	defer ln.Close()
	return true
}

func getRandomPort() int {
	rand.Seed(time.Now().UnixNano())
	minPort := 1024
	maxPort := 65535

	for {
		port := rand.Intn(maxPort-minPort+1) + minPort
		if !isWellKnownPort(port) {
			return port
		}
	}
}

func generateRandomPort(protocol string) int {
	var port int
	for {
		port = getRandomPort()
		if !isWellKnownPort(port) && isPortAvailable(port, protocol) {
			break
		}
	}
	return port
}

func listUsedPorts(cmd *cobra.Command, args []string) {
	var protocol string

	tcpFlag, _ := cmd.Flags().GetBool("tcp")
	udpFlag, _ := cmd.Flags().GetBool("udp")
	allFlag, _ := cmd.Flags().GetBool("all")

	if tcpFlag {
		protocol = "tcp"
	} else if udpFlag {
		protocol = "udp"
	} else if allFlag {
		fmt.Println("Listing all used ports:")
		listAllUsedPorts()
		return
	} else {
		fmt.Println("Please specify a flag (-t for TCP, -u for UDP, -a for all) for listing used ports.")
		return
	}

	fmt.Printf("Listing used %s ports:\n", protocol)

	usedPorts, err := getUsedPorts(protocol)
	if err != nil {
		fmt.Printf("Error listing used ports: %v\n", err)
		return
	}

	colorizeAndPrint(usedPorts, protocol)
}

func listAllUsedPorts() {
	tcpPorts, err := getUsedPorts("tcp")
	if err != nil {
		fmt.Printf("Error listing used TCP ports: %v\n", err)
		return
	}

	udpPorts, err := getUsedPorts("udp")
	if err != nil {
		fmt.Printf("Error listing used UDP ports: %v\n", err)
		return
	}

	maxLen := len(tcpPorts)
	if len(udpPorts) > maxLen {
		maxLen = len(udpPorts)
	}

	fmt.Println("TCP Ports\t\tUDP Ports:")

	for i := 0; i < maxLen; i++ {
		tcpPort := ""
		if i < len(tcpPorts) {
			tcpPort = strconv.Itoa(tcpPorts[i])
		}

		udpPort := ""
		if i < len(udpPorts) {
			udpPort = strconv.Itoa(udpPorts[i])
		}

		fmt.Printf("%s\t\t\t%s\n", color.RedString(tcpPort), color.GreenString(udpPort))
	}
}

func getUsedPorts(protocol string) ([]int, error) {
	cmd := exec.Command("ss", "-tuln")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running ss command: %v", err)
	}

	return parseSSOutput(string(output), protocol), nil
}

func parseSSOutput(output, protocol string) []int {
	var usedPorts []int

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 6 && fields[0] == protocol {
			portStr := fields[4]
			port, err := strconv.Atoi(strings.Split(portStr, ":")[1])
			if err == nil {
				usedPorts = append(usedPorts, port)
			}
		}
	}

	return usedPorts
}

func colorizeAndPrint(ports []int, protocol string) {
	var colorFunc func(...interface{}) string

	if protocol == "tcp" {
		colorFunc = color.New(color.FgRed).SprintFunc()
	} else if protocol == "udp" {
		colorFunc = color.New(color.FgGreen).SprintFunc()
	}

	for _, port := range ports {
		fmt.Printf("%s\n", colorFunc(port))
	}
	fmt.Println()
}

func main() {

	var listCmd = &cobra.Command{
		Use:   "list-active",
		Short: "List used ports",
		Run:   listUsedPorts,
	}

	listCmd.Flags().BoolP("tcp", "t", false, "List used TCP ports")
	listCmd.Flags().BoolP("udp", "u", false, "List used UDP ports")
	listCmd.Flags().BoolP("all", "a", false, "List all used ports (TCP and UDP)")

	var checkCmd = &cobra.Command{
		Use:   "check-well-known",
		Short: "Check a well-known port",
		Run:   checkWellKnownPort,
	}

	var silent bool

	var rootCmd = &cobra.Command{Use: "port-randomizer"}
	rootCmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "Silence gamified output")

	rootCmd.AddCommand(&cobra.Command{
		Use: "randomize",
		Short: "Randomly generates a tcp port.",
		Run: func(cmd *cobra.Command, args []string) {	
			protocol := "tcp"
	
			if silent {
				randomPort := generateRandomPort(protocol)
				fmt.Println(randomPort)
			} else {
				// Gamify the output
				fmt.Print("Cycling through ports:")
	
				for i := 0; i < 10; i++ {
					randomPort := getRandomPort()
					fmt.Printf("\033[91m%d \033[0m", randomPort) // Print in red
					time.Sleep(100 * time.Millisecond)
	
					// Clear only the generated random port (assuming Unix-like systems)
					fmt.Print("\033[1K") // Clear the line from the beginning
					fmt.Print("\033[0G") // Move the cursor to the beginning
	
					cmd := exec.Command("clear") // Assuming Unix-like systems, you can use "cls" for Windows
					cmd.Stdout = os.Stdout
					cmd.Run()
					fmt.Print("Randomizing ports:")
				}
	
				randomPort := generateRandomPort(protocol)
				cmd := exec.Command("clear") // Assuming Unix-like systems, you can use "cls" for Windows
				cmd.Stdout = os.Stdout
				cmd.Run()
				fmt.Printf("\033[92m%d \033[0m\n", randomPort) // Print in green
			}
		},
		Args: cobra.NoArgs, // Accept no additional arguments
	})

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(checkCmd)


	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}