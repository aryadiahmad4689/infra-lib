package helm

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Helm struct untuk mengelola operasi Helm
type Helm struct {
	BasePath string
}

// NewHelm adalah constructor untuk Helm struct
func NewHelm(basePath string) *Helm {
	return &Helm{BasePath: basePath}
}

// AddAndPullCharts menambahkan repository dari file dan menarik chart-nya ke folder tools
func (h *Helm) AddAndPullCharts(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Error opening file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "helm repo add") {
			if err := h.processHelmRepoCommand(line); err != nil {
				fmt.Printf("Error processing Helm repo command: %v\n", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading file: %v", err)
	}

	// Cleanup the list-tools directory after extraction
	if err := h.cleanupListToolsFolder(); err != nil {
		return fmt.Errorf("Error cleaning up list-tools directory: %v", err)
	}

	return nil
}

// processHelmRepoCommand memproses perintah helm repo add dan helm pull
func (h *Helm) processHelmRepoCommand(line string) error {
	// Pecah perintah dengan spasi untuk mengambil repoName, repoURL, dan chartName
	parts := strings.Split(line, " ")
	if len(parts) < 6 {
		return fmt.Errorf("Invalid repo line: %s. Format should be: helm repo add <repo_name> <repo_url> <chart_name>", line)
	}

	repoName := parts[3]
	repoURL := parts[4]
	chartName := parts[5]

	// Menambahkan repository ke Helm
	fmt.Printf("Adding repository: %s\n", repoName)
	cmd := exec.Command("helm", "repo", "add", repoName, repoURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error adding repo %s: %v", repoName, err)
	}

	// Tarik chart dalam bentuk .tgz jika belum ada
	fmt.Printf("Pulling chart for repo: %s, chart: %s\n", repoName, chartName)
	cmd = exec.Command("helm", "pull", repoName+"/"+chartName, "--destination", h.BasePath+"/tools")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error pulling chart %s: %v", chartName, err)
	}

	// Temukan file .tgz berdasarkan prefix nama chart di dalam folder tools
	tgzFile, err := h.findTGZFile(h.BasePath+"/tools", chartName)
	if err != nil {
		return fmt.Errorf("Error finding .tgz file for chart %s: %v", chartName, err)
	}

	// Ekstrak file .tgz
	if err := h.extractTGZ(tgzFile); err != nil {
		return fmt.Errorf("Error extracting chart %s: %v", chartName, err)
	}

	// Hapus file .tgz setelah diekstrak
	if err := os.Remove(tgzFile); err != nil {
		return fmt.Errorf("Error removing .tgz file %s: %v", tgzFile, err)
	}

	return nil
}

// findTGZFile mencari file .tgz berdasarkan prefix nama chart di direktori tools
func (h *Helm) findTGZFile(directory, chartName string) (string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		// Cari file yang diawali dengan chartName dan diakhiri dengan .tgz
		if strings.HasPrefix(file.Name(), chartName+"-") && strings.HasSuffix(file.Name(), ".tgz") {
			return filepath.Join(directory, file.Name()), nil
		}
	}

	return "", fmt.Errorf("no .tgz file found for chart %s", chartName)
}

// extractTGZ extracts all files from a .tgz archive directly into the tools directory, preserving internal structure
func (h *Helm) extractTGZ(tgzFile string) error {
	baseDir := filepath.Dir(tgzFile) // Extracts to the tools directory

	file, err := os.Open(tgzFile)
	if err != nil {
		return fmt.Errorf("Error opening .tgz file %s: %v", tgzFile, err)
	}
	defer file.Close()

	// Create a gzip reader
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("Error creating gzip reader: %v", err)
	}
	defer gzr.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzr)

	// Iterate over each item in the archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of archive
			break
		}
		if err != nil {
			return fmt.Errorf("Error reading tar archive: %v", err)
		}

		// Define the full path for the destination file or directory
		destPath := filepath.Join(baseDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory if it doesn't exist, including the entire path
			if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
				return fmt.Errorf("Error creating directory %s: %v", destPath, err)
			}
		case tar.TypeReg:
			// Ensure the parent directory exists
			if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
				return fmt.Errorf("Error creating directory %s: %v", filepath.Dir(destPath), err)
			}
			// Create and write the file
			outFile, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf("Error creating file %s: %v", destPath, err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("Error writing to file %s: %v", destPath, err)
			}
			outFile.Close()
		}
	}
	return nil
}

// cleanupListToolsFolder removes the list-tools directory if it exists
func (h *Helm) cleanupListToolsFolder() error {
	listToolsPath := filepath.Join(h.BasePath, "list-tools")
	if _, err := os.Stat(listToolsPath); err == nil {
		// Path exists, so remove it
		if err := os.RemoveAll(listToolsPath); err != nil {
			return fmt.Errorf("Error removing list-tools directory: %v", err)
		}
	}
	return nil
}
