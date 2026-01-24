package file

import (
    "fmt"
    "os/exec"
    "path/filepath"
    "runtime"
)

func moveToTrash(path string) error {
    // Get absolute path
    absPath, err := filepath.Abs(path)
    if err != nil {
        return fmt.Errorf("failed to get absolute path: %w", err)
    }
    
    var cmd *exec.Cmd
    
    switch runtime.GOOS {
    case "darwin":
        // macOS: use Finder via osascript
        script := fmt.Sprintf(`tell application "Finder" to delete POSIX file "%s"`, absPath)
        cmd = exec.Command("osascript", "-e", script)
        
    case "windows":
        // Windows: use PowerShell to move to Recycle Bin
        psScript := fmt.Sprintf(
            `Add-Type -AssemblyName Microsoft.VisualBasic; `+
            `[Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile('%s', `+
            `'OnlyErrorDialogs', 'SendToRecycleBin')`,
            absPath,
        )
        cmd = exec.Command("powershell", "-NoProfile", "-Command", psScript)
        
    default:
        // Linux/Unix: use gio trash
        cmd = exec.Command("gio", "trash", absPath)
    }
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to move to trash: %w", err)
    }
    
    return nil
}
