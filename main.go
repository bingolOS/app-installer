package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Clear screen helper
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

// ANSI colors
const (
	colorReset = "\033[0m"
	colorGold  = "\033[38;5;220m" // 256-color "gold"
	colorCyan  = "\033[36m"
)

// Print ASCII Banner
func printBanner() {
	// New ASCII logo (your choice)
	logo := `
            _____________
        .-'             '-.
      .'   /\       /\     '.
     /    /  \_____/  \      \
    /    /             \      \
   |    /   /\     /\   \     |
   |   /   /  \___/  \   \    |
   |  |   |           |   |   |
   |  |   |     Y     |   |   |
   |  |   |    / \    |   |   |
   |   \   \__/   \__/   /    |
    \    \             /     /
     \    \___________/     /
      '.                 .'
        '-.___________.-'
`

	title := `
   Paket Yöneticisi
`

	// Print in GOLD
	fmt.Print(colorGold + logo + title + colorReset)

	// (İstersen alttaki yazıyı cyan bırakabilirsin)
	// fmt.Print(colorCyan + "Paket Yöneticisi\n" + colorReset)
}

// Check if a command exists
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Install Package function
func installPackage() {
	// check for yay or pacman
	pkgManager := "pacman"
	if commandExists("yay") {
		pkgManager = "yay"
	}

	fmt.Println("Paket listesi getiriliyor... (Aramak için yazın)")

	usageCmd := fmt.Sprintf(
		`%s -Slq | fzf --preview '%s -Si {}' --layout=reverse --height=90%% --header='YÜKLEMEK için paket seçin (Enter)'`,
		pkgManager, pkgManager,
	)

	cmd := exec.Command("bash", "-c", usageCmd)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return
	}

	pkgName := strings.TrimSpace(string(output))
	if pkgName == "" {
		return
	}

	clearScreen()
	fmt.Printf("%s paketi yükleniyor...\n", pkgName)

	installCmd := exec.Command(pkgManager, "-S", pkgName)
	installCmd.Stdin = os.Stdin
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	_ = installCmd.Run()

	fmt.Println("\nMenüye dönmek için Enter'a basın...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Remove Package function
func removePackage() {
	fmt.Println("Kurulu paketler getiriliyor...")

	usageCmd := "pacman -Qq | fzf -m --preview 'pacman -Qi {}' --layout=reverse --height=90% --header='KALDIRMAK için paket(leri) seçin (TAB ile çoklu seçim, Onaylamak için Enter)'"

	cmd := exec.Command("bash", "-c", usageCmd)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return
	}

	// Split lines into arguments
	pkgs := strings.Split(result, "\n")

	clearScreen()
	fmt.Printf("Kaldırılıyor: %s\n", strings.Join(pkgs, ", "))

	args := append([]string{"-Rns"}, pkgs...)
	removeCmd := exec.Command("sudo", append([]string{"pacman"}, args...)...)
	removeCmd.Stdin = os.Stdin
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr
	_ = removeCmd.Run()

	fmt.Println("\nMenüye dönmek için Enter'a basın...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func main() {
	// Verify fzf exists
	if !commandExists("fzf") {
		fmt.Println("Hata: 'fzf' gerekli ancak bulunamadı. Lütfen yükleyin (sudo pacman -S fzf).")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		printBanner()
		fmt.Println("1. [Y]ükle - Paket Ara")
		fmt.Println("2. [K]aldır - Paket Sil (Çoklu Seçim)")
		fmt.Println("3. [C]ıkış")
		fmt.Print("\nSeçiminiz: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "1", "y", "i":
			installPackage()
		case "2", "k", "r":
			removePackage()
		case "3", "c", "q", "e":
			fmt.Println("Görüşürüz!")
			return
		default:
			continue
		}
	}
}
