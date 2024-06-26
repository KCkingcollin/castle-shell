package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"
)

const hyprContent string = `# color vars generated by theme-checker.service from ~/.config/castle-shell/accent-color
$prime_color = rgba(%sa6)
$alt_color = rgba(%s8c)
`

const kittyContent string = `# color vars generated by theme-checker.service from ~/.config/castle-shell/accent-color
foreground #%s
background #%s
background_opacity 0.35
`

const cssContent string = `/*color vars generated by theme-checker.service from ~/.config/castle-shell/accent-color*/
@define-color primaryColor rgba(%s, 0.5);
@define-color secondaryColor rgba(%s, 0.35);
@define-color secondaryColorDark rgba(%s, 0.45);
@define-color secondaryColorDarker rgba(%s, 0.6);
@define-color textColor rgba(%s, 1);
`

const rasiContent string = `/*color vars generated by theme-checker.service from ~/.config/castle-shell/accent-color*/
* {
    primaryColor: rgba(%s, 0.5);
    secondaryColor: rgba(%s, 0.35);
    secondaryColorDark: rgba(%s, 0.45);
    secondaryColorDarker: rgba(%s, 0.6);
    textColor: rgba(%s, 1);
}
`

func main() {
    var home = os.Getenv("HOME")
    var colorFile = path.Join(home, ".config/castle-shell/accent-color")
    var hyprConfFile = path.Join(home, ".config/castle-shell/hypr-colors.conf")
    var kittyConfFile = path.Join(home, ".config/castle-shell/kitty-colors.conf")
    var cssFile = path.Join(home, ".config/castle-shell/colors.css")
    var rasiFile = path.Join(home, ".config/castle-shell/colors.rasi")

    currentHash := genHash(colorFile)
    for {
        time.Sleep(250 * time.Millisecond)
        start := time.Now()

        newHash := genHash(colorFile)

        // loop again if the file hasn't changed yet
        if currentHash == newHash {
            continue
        } else {
            fmt.Println("Accent file has changed")
            currentHash = newHash 
        }

        // //debugging
        // fmt.Println("File hash:", currentHash)

        // get rbg values
        outPrime := getRgbFromFile(colorFile, 1)
        outAlt := getRgbFromFile(colorFile, 2)

        // convert the values
        hexPrime := rgb2hex(outPrime)
        hexAlt := rgb2hex(outAlt)
        cssPrime := num2css(outPrime)
        cssAlt := num2css(outAlt)
        hexText := "ffffff"
        cssText := "255, 255, 255"

        // send generated Hyprland config to file
        file, err := os.OpenFile(hyprConfFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
        if err != nil {
            fmt.Printf("Error opening file: %s", err)
        }
        _, err = fmt.Fprintf(file, hyprContent, hexPrime, hexAlt)
        if err != nil {
            fmt.Printf("Error writing to file: %s", err)
        }
        file.Close()

        // send generated Kitty config to file
        file, err = os.OpenFile(kittyConfFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
        if err != nil {
            fmt.Printf("Error opening file: %s", err)
        }
        _, err = fmt.Fprintf(file, kittyContent, hexText, hexAlt)
        if err != nil {
            fmt.Printf("Error writing to file: %s", err)
        }
        file.Close()
        reloadKittyConfig()

        // send generated css config to file
        file, err = os.OpenFile(cssFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
        if err != nil {
            fmt.Printf("Error opening file: %s", err)
        }
        _, err = fmt.Fprintf(file, cssContent, cssPrime, cssAlt, cssAlt, cssAlt, cssText)
        if err != nil {
            fmt.Printf("Error writing to file: %s", err)
        }
        file.Close()

        // send generated rasi config to file
        file, err = os.OpenFile(rasiFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
        if err != nil {
            fmt.Printf("Error opening file: %s", err)
        }
        _, err = fmt.Fprintf(file, rasiContent, cssPrime, cssAlt, cssAlt, cssAlt, cssText)
        if err != nil {
            fmt.Printf("Error writing to file: %s", err)
        }
        file.Close()

        // // Print updated colors for debugging
        // fmt.Printf("Primary CSS Color: rgba(%s, 0.5)\n", cssPrime)
        // fmt.Printf("Secondary CSS Color: rgba(%s, 0.35)\n", cssAlt)
        // fmt.Printf("Dark Secondary CSS Color: rgba(%s, 0.45)\n", cssAlt)
        // fmt.Printf("Darker Secondary CSS Color: rgba(%s, 0.6)\n", cssAlt)
        // fmt.Printf("Hyprland Primary Color: rgba(%sa6)\n", hexPrime)
        // fmt.Printf("Hyprland Secondary Color: rgba(%s8c)\n", hexAlt)

        cmd := exec.Command("systemctl", "restart", "--user", "waybar-hyprland.service", "swaync.service")
        _, err = cmd.Output()
        if err != nil {
            fmt.Printf("Error restarting services: %s", err)
        }

        end := time.Now()
        mainLoopTime := end.Sub(start)
        fmt.Printf("Main loop took %v to complete.\n", mainLoopTime)

        // tmep break for testing
        // break
    }
}

// create hash from file
func genHash(filePath string) (string) {
    file, err := os.Open(filePath)
    if err != nil {
        fmt.Println("Error opening file:", err)
    }
    defer file.Close()

    hash := sha256.New()

    if _, err := io.Copy(hash, file); err != nil {
        fmt.Println("Error hashing file:", err)
    }

    hashInBytes := hash.Sum(nil)
    hashInHex := fmt.Sprintf("%x", hashInBytes)

    return hashInHex
}

// create a hex value without the # from 3 uint8s
func rgb2hex(rgb [3]uint8) (string) {
    // % Format, 0 Padded, 2 Width, X Uppercase hex
    r := rgb[0]
    g := rgb[1]
    b := rgb[2]
    return fmt.Sprintf("%02X%02X%02X", r, g, b)
}

// create a css formatted string from 3 uint8s
func num2css(rgb [3]uint8) (string) {
    r := rgb[0]
    g := rgb[1]
    b := rgb[2]
    return fmt.Sprintf("%v, %v, %v", r, g, b)
}

// replaces a line in a file with another string
func getRgbFromFile(filePath  string, line int) ([3]uint8) {
    file, err := os.Open(filePath)
    if err != nil {
        fmt.Println("Error opening file:", err)
    }
    defer file.Close()

    // scanner reads a file line by line
    scanner := bufio.NewScanner(file)

    // scan through lines 
    var lineNumber int
    var targetLine string
    for scanner.Scan() {
        lineNumber++
        if lineNumber == line {
            targetLine = scanner.Text()
            break
        }
    }

    // split function that reads word by word
    scanner = bufio.NewScanner(strings.NewReader(targetLine))
    scanner.Split(bufio.ScanWords)

    // scan numbers
    var numbers [3]uint8
    index := 0
    for scanner.Scan() {
        if index >= 3 {
            break
        }
        token := scanner.Text()
        outPut := removeNonDigits(token)
        value, err := strconv.Atoi(outPut)
        if err != nil {
            fmt.Printf("Error converting token: %v\n", err)
        }
        if value < 0 || value > 255 {
            fmt.Printf("Error int out of range, can only be between 0 - 255, value: %d\n", value)
        }
        numbers[index] = uint8(value)
        index++
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Fatal error in scanner: %v\n", err)
    }

    for index < 3 {
        fmt.Printf("Not enough values in the line\n")
    }

    return numbers
}

func reloadKittyConfig() {
    cmd := exec.Command("pgrep", "-f", "kitty")

    kittyPids, err := cmd.Output()
    if err != nil {
        fmt.Println("Error running command:", err)
    }
    if kittyPids != nil {
        scanner := bufio.NewScanner(strings.NewReader(string(kittyPids)))
        for scanner.Scan() {
            pidStr := scanner.Text()
            pid, err := strconv.Atoi(pidStr)
            if err != nil {
                fmt.Printf("Invalid pid: %v\n", err)
            }
            err = syscall.Kill(pid, syscall.SIGUSR1)
            if err != nil {
                fmt.Printf("Error Failed to send sig: %v\n", err)
            }
        }
    }
}

func removeNonDigits(input string) (string) {
    var sb strings.Builder
    for _, r := range input {
        if unicode.IsDigit(r) {
            sb.WriteRune(r)
        }
    }
    return sb.String()
}
