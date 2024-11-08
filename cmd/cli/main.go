package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aryadiahmad4689/infra-lib/src/folder"
	"github.com/aryadiahmad4689/infra-lib/src/helm"
)

// fungsi untuk memeriksa apakah Helm sudah terinstal
func isHelmInstalled() bool {
	_, err := exec.LookPath("helm")
	return err == nil
}

// fungsi untuk meminta input dari pengguna
func promptUser(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", prompt)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y"
}

func main() {
	// Check if there are enough arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: infra create <base_folder_name>")
		return
	}

	// Ambil argumen base path
	basePath := os.Args[2]

	// Inisialisasi Folder struct dengan basePath
	folder := folder.NewFolder(basePath)

	// Buat struktur folder
	err := folder.Create()
	if err != nil {
		fmt.Printf("Failed to create directory structure: %v\n", err)
		return
	}

	fmt.Println("All directories created successfully.")

	helmClient := helm.NewHelm(basePath)
	if err := helmClient.AddAndPullCharts("src/list-tools.txt"); err != nil {
		fmt.Printf("Failed to add and pull Helm charts: %v\n", err)
		return
	}

	fmt.Println("Infrastructure setup completed successfully.")

	// Tanyakan kepada pengguna apakah mereka ingin memuat semua tools
	if promptUser("Do you want to load all tools") {
		// Periksa apakah Helm sudah diinstal
		if !isHelmInstalled() {
			fmt.Println("Helm is not installed. Please install Helm and rerun the command: infra tools")
			return
		}

		// Inisialisasi Helm struct dan tarik chart
		helmClient := helm.NewHelm(basePath)
		if err := helmClient.AddAndPullCharts("src/list-tools.txt"); err != nil {
			fmt.Printf("Failed to add and pull Helm charts: %v\n", err)
			return
		}

		fmt.Println("Infrastructure setup completed successfully.")
	} else {
		fmt.Println("Skipping tool installation.")
	}
}
