package main

import (
    "log"
    "fmt"
    "os"
    "bufio"
    "strings"
    "path/filepath"
    "io"
    // "time"

    "github.com/micmonay/keybd_event"
)

const shortCutsFile = `#
# DO NOT TOUCH
source = /tmp/current-mode.conf
bind = Ctrl Alt, c, exec, wl-copy < /dev/null
bind = Ctrl Alt, v, exec, gvm toggle
# DO NOT TOUCH
#
`

func main() {
    // Check if an argument was provided
    if len(os.Args) < 2 {
        fmt.Println("Please provide at least one argument.")
        return
    }

    // Get the first argument
    arg1 := os.Args[1]

    // Check if the second argument was provided
    var arg2 string
    if len(os.Args) >= 3 {
        arg2 = os.Args[2]
    } else {
        // If the second argument is not provided, set a default value or handle accordingly
        arg2 = "default_value" // You can set a meaningful default value here
    }

    // Use the arguments for some purpose
    fmt.Println("Input argument 1 received:", arg1)
    fmt.Println("Input argument 2 received:", arg2)

    // Example: Use the arguments to perform some operation
    performOperation(arg1, arg2)
}

// Example function that uses the arg1 variable
func performOperation(arg1, arg2 string) {
    // Create a new keyboard event instance
    kb, err := keybd_event.NewKeyBonding()
    if err != nil {
        log.Fatal(err)
    }
    // Perform some operation based on the arg1
    switch arg1 {
    case "toggle":
        toggleMotions()
    case "i":
        err := copyFile("/usr/share/global-vim-motions/insert.conf", "/tmp/current-mode.conf")
        if err != nil {
            fmt.Println("Error copying file:", err)
        } else {
            fmt.Println("File copied successfully.")
        }

    case "n":
        err := copyFile("/usr/share/global-vim-motions/normal.conf", "/tmp/current-mode.conf")
        if err != nil {
            fmt.Println("Error copying file:", err)
        } else {
            fmt.Println("File copied successfully.")
        }

        switch arg2 {
        case "down":
            fmt.Println("down")
            kb.SetKeys(keybd_event.VK_DOWN)
            if err := kb.Launching(); err != nil {
                log.Fatal(err)
            }
        case "up":
            fmt.Println("up")
            kb.SetKeys(keybd_event.VK_UP)
            if err := kb.Launching(); err != nil {
                log.Fatal(err)
            }
        case "left":
            fmt.Println("left")
            kb.SetKeys(keybd_event.VK_LEFT)
            if err := kb.Launching(); err != nil {
                log.Fatal(err)
            }
        case "right":
            fmt.Println("right")
            kb.SetKeys(keybd_event.VK_RIGHT)
            if err := kb.Launching(); err != nil {
                log.Fatal(err)
            }
        default:
            fmt.Println("Unknown arg2:", arg2)
        }
    case "v":
        err := copyFile("/usr/share/global-vim-motions/visual.conf", "/tmp/current-mode.conf")
        if err != nil {
            fmt.Println("Error copying file:", err)
        } else {
            fmt.Println("File copied successfully.")
        }

        switch arg2 {
        case "line":
            fmt.Println("home shift+end")
            // Simulate "Home" key press
            kb.SetKeys(keybd_event.VK_HOME)
            if err := kb.Launching(); err != nil {
                log.Fatal(err)
            }

            // Simulate "Shift+End" key combo
            kb.Clear()
            kb.SetKeys(keybd_event.VK_END)
            kb.HasSHIFT(true)
            if err := kb.Launching(); err != nil {
                log.Fatal(err)
            }
        case "down":
            fmt.Println("shift+down")
            kb.SetKeys(keybd_event.VK_DOWN)
            kb.HasSHIFT(true)
            if err := kb.Launching(); err != nil {
                log.Fatal(err)
            }
        case "up":
            fmt.Println("shift+up")
            kb.SetKeys(keybd_event.VK_UP)
            kb.HasSHIFT(true)
            if err := kb.Launching(); err != nil {
                log.Fatal(err)
            }
        case "left":
            fmt.Println("shift+left")
            kb.SetKeys(keybd_event.VK_LEFT)
            kb.HasSHIFT(true)
            if err := kb.Launching(); err != nil {
                log.Fatal(err)
            }
        case "right":
            fmt.Println("shift+right")
            kb.SetKeys(keybd_event.VK_RIGHT)
            kb.HasSHIFT(true)
            if err := kb.Launching(); err != nil {
                log.Fatal(err)
            }
        default:
            fmt.Println("Unknown arg2:", arg2)
        }
    default:
        fmt.Println("Unknown arg1:", arg1)
    }
    switch arg2 {
    default:
    }
}

func copyFile(src string, dst string) error {
    // Open the source file
    sourceFile, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("failed to open source file: %w", err)
    }
    defer sourceFile.Close()

    // Create the destination file
    destinationFile, err := os.Create(dst)
    if err != nil {
        return fmt.Errorf("failed to create destination file: %w", err)
    }
    defer destinationFile.Close()

    // Copy the contents from the source file to the destination file
    _, err = io.Copy(destinationFile, sourceFile)
    if err != nil {
        return fmt.Errorf("failed to copy file contents: %w", err)
    }

    // Flush the writer buffer if applicable
    err = destinationFile.Sync()
    if err != nil {
        return fmt.Errorf("failed to flush to destination file: %w", err)
    }

    return nil
}

func ensureFileExists(filename string) error {
    // Check if the file exists
    _, err := os.Stat(filename)
    if os.IsNotExist(err) {
        // File does not exist, create it
        file, err := os.Create(filename)
        if err != nil {
            return fmt.Errorf("failed to create file: %w", err)
        }
        defer file.Close()

        fmt.Println("File created successfully.")
    } else if err != nil {
        return fmt.Errorf("failed to check if file exists: %w", err)
    } else {
        fmt.Println("File already exists.")
    }

    return nil
}

func appendTextToFile(filename, text string) error {
    // Open the file with the appropriate flags to create if not exists and append
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    // Write the text to the file
    _, err = file.WriteString(text)
    if err != nil {
        return fmt.Errorf("failed to write to file: %w", err)
    }

    return nil
}

func toggleMotions() {
    if err := ensureFileExists("/tmp/current-mode.conf"); err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Curent-mode file checked or created successfully.")
    }

    err := copyFile("/usr/share/global-vim-motions/normal.conf", "/tmp/current-mode.conf")
    if err != nil {
        fmt.Println("Error copying file:", err)
    } else {
        fmt.Println("File copied successfully.")
    }

    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Error getting home directory:", err)
        return
    }

    filename := filepath.Join(homeDir, ".config/castle-shell/global-vim-motions/shortcuts.conf")

    if err := ensureFileExists(filename); err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Shortcuts file checked or created successfully.")
    }

    // Open the input file for reading
    inputFile, err := os.Open(filename)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer inputFile.Close()

    // Get file info
    fileInfo, err := inputFile.Stat()
    if err != nil {
        fmt.Println("Error getting file info:", err)
    }

    // Check if the file is empty
    if fileInfo.Size() == 0 {
        fmt.Println("The file is empty.")
        if err := appendTextToFile(filename, shortCutsFile); err != nil {
            fmt.Println("Error:", err)
        } else {
            fmt.Println("Text added to the shortcuts file successfully.")
        }
        return
    } else {
        fmt.Println("The file is not empty.")
    }

    // Read the file line by line
    var lines []string
    scanner := bufio.NewScanner(inputFile)
    for scanner.Scan() {
        line := scanner.Text()
        // Check if the line contains "foo"
        if strings.Contains(line, "# source = /tmp/current-mode.conf") {
            // Replace "bar" with "foo"
            line = strings.Replace(line, "# source = /tmp/current-mode.conf", "source = /tmp/current-mode.conf", -1)

            // Check if the line contains "bar"
        } else if strings.Contains(line, "source = /tmp/current-mode.conf") {
            // Replace "foo" with "bar"
            line = strings.Replace(line, "source = /tmp/current-mode.conf", "# source = /tmp/current-mode.conf", -1)
        } 
        lines = append(lines, line)
    }

    // Check for errors during scanning
    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    // Open the file for writing
    outputFile, err := os.Create(filename)
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer outputFile.Close()

    // Write the modified lines back to the file
    writer := bufio.NewWriter(outputFile)
    for _, line := range lines {
        _, err := writer.WriteString(line + "\n")
        if err != nil {
            fmt.Println("Error writing to file:", err)
            return
        }
    }

    writer.Flush()

    fmt.Println("File updated successfully.")
}
