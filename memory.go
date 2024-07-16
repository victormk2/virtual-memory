package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Page struct {
	size      int
	allocated int
	fileName  string
}

type Memory struct {
	pages     []Page
	maxSize   int
	allocated int
}

func (m *Memory) AllocateFileFirstFit(fileName string) bool {
	file, err := os.Stat(fileName)
	if err != nil {
		fmt.Print("Could not allocate file\n")
		return false
	}
	fileSize := int(file.Size())

	if m.allocated+fileSize > m.maxSize {
		return false
	}

	for i := 0; i < len(m.pages); i++ {
		// Tentar alocar em uma única página livre
		if m.pages[i].size >= fileSize && m.pages[i].fileName == "" {
			if m.pages[i].size > fileSize {
				newPage := Page{
					size:      m.pages[i].size - fileSize,
					allocated: 0,
					fileName:  "",
				}
				m.pages[i].size = fileSize
				m.pages = append(m.pages[:i+1], append([]Page{newPage}, m.pages[i+1:]...)...)
			}
			m.pages[i].allocated = fileSize
			m.pages[i].fileName = fileName
			m.allocated += fileSize
			return true
		}

		// Tentar alocar juntando páginas adjacentes livres
		totalFreeSize := 0
		startIdx := -1
		for j := i; j < len(m.pages) && m.pages[j].fileName == ""; j++ {
			if startIdx == -1 {
				startIdx = j
			}
			totalFreeSize += m.pages[j].size
			if totalFreeSize >= fileSize {
				// Juntar as páginas livres e alocar o arquivo
				m.pages[startIdx].size = totalFreeSize
				m.pages[startIdx].allocated = fileSize
				m.pages[startIdx].fileName = fileName
				m.allocated += fileSize

				// Remover as páginas extras
				m.pages = append(m.pages[:startIdx+1], m.pages[startIdx+1+j-startIdx:]...)
				return true
			}
		}
	}

	pagesSize := 0
	for j := 0; j < len(m.pages); j++ {
		pagesSize += m.pages[j].size
	}

	if (m.maxSize - pagesSize) < fileSize {
		fmt.Println("There are no available pages to allocate the file.")
		return false
	}

	// Caso não tenha encontrado um espaço adequado, criar uma nova página.
	if m.maxSize-m.allocated >= fileSize {
		m.pages = append(m.pages, Page{size: fileSize, allocated: fileSize, fileName: fileName})
		m.allocated += fileSize
		return true
	}

	return false
}

func (m *Memory) AllocateFileWorstFit(fileName string) bool {
	file, err := os.Stat(fileName)
	if err != nil {
		fmt.Print("Could not allocate file\n")
		return false
	}
	fileSize := int(file.Size())

	if m.allocated+fileSize > m.maxSize {
		return false
	}

	worstFitIndex := -1
	worstFitSize := 0

	// Encontrar a maior página livre
	for i := range m.pages {
		if m.pages[i].size >= fileSize && m.pages[i].fileName == "" {
			if m.pages[i].size > worstFitSize {
				worstFitIndex = i
				worstFitSize = m.pages[i].size
			}
		}
	}

	// Alocar na maior página livre encontrada
	if worstFitIndex != -1 {
		if m.pages[worstFitIndex].size > fileSize {
			newPage := Page{
				size:      m.pages[worstFitIndex].size - fileSize,
				allocated: 0,
				fileName:  "",
			}
			m.pages[worstFitIndex].size = fileSize
			m.pages = append(m.pages[:worstFitIndex+1], append([]Page{newPage}, m.pages[worstFitIndex+1:]...)...)
		}
		m.pages[worstFitIndex].allocated = fileSize
		m.pages[worstFitIndex].fileName = fileName
		m.allocated += fileSize
		return true
	} else {
		// Tentar juntar páginas adjacentes para encontrar espaço suficiente
		totalFreeSize := 0
		startIdx := -1
		for i := 0; i < len(m.pages); i++ {
			if m.pages[i].fileName == "" {
				if startIdx == -1 {
					startIdx = i
				}
				totalFreeSize += m.pages[i].size
				if totalFreeSize >= fileSize {
					// Juntar as páginas livres e alocar o arquivo
					m.pages[startIdx].size = totalFreeSize
					m.pages[startIdx].allocated = fileSize
					m.pages[startIdx].fileName = fileName
					m.allocated += fileSize

					// Remover as páginas extras
					m.pages = append(m.pages[:startIdx+1], m.pages[startIdx+1+(i-startIdx):]...)
					return true
				}
			} else {
				totalFreeSize = 0
				startIdx = -1
			}
		}
	}

	pagesSize := 0
	for j := 0; j < len(m.pages); j++ {
		pagesSize += m.pages[j].size
	}

	if (m.maxSize - pagesSize) < fileSize {
		fmt.Println("There are no available pages to allocate the file.")
		return false
	}

	// Caso não tenha encontrado um espaço adequado, criar uma nova página.
	if m.maxSize-m.allocated >= fileSize {
		m.pages = append(m.pages, Page{size: fileSize, allocated: fileSize, fileName: fileName})
		m.allocated += fileSize
		return true
	}

	return false
}

func (m *Memory) AllocateFileBestFit(fileName string) bool {
	file, err := os.Stat(fileName)
	if err != nil {
		fmt.Print("Could not allocate file\n")
		return false
	}
	fileSize := int(file.Size())

	if m.allocated+fileSize > m.maxSize {
		return false
	}

	bestFitIndex := -1
	bestFitSize := m.maxSize + 1

	// Encontrar a menor página livre que seja suficiente para alocar o arquivo
	for i := range m.pages {
		if m.pages[i].size >= fileSize && m.pages[i].fileName == "" {
			if m.pages[i].size < bestFitSize {
				bestFitIndex = i
				bestFitSize = m.pages[i].size
			}
		}
	}

	// Alocar na melhor página encontrada
	if bestFitIndex != -1 {
		if m.pages[bestFitIndex].size > fileSize {
			newPage := Page{
				size:      m.pages[bestFitIndex].size - fileSize,
				allocated: 0,
				fileName:  "",
			}
			m.pages[bestFitIndex].size = fileSize
			m.pages = append(m.pages[:bestFitIndex+1], append([]Page{newPage}, m.pages[bestFitIndex+1:]...)...)
		}
		m.pages[bestFitIndex].allocated = fileSize
		m.pages[bestFitIndex].fileName = fileName
		m.allocated += fileSize
		return true
	} else {
		// Tentar juntar páginas adjacentes para encontrar espaço suficiente
		totalFreeSize := 0
		startIdx := -1
		for i := 0; i < len(m.pages); i++ {
			if m.pages[i].fileName == "" {
				if startIdx == -1 {
					startIdx = i
				}
				totalFreeSize += m.pages[i].size
				if totalFreeSize >= fileSize {
					// Juntar as páginas livres e alocar o arquivo
					m.pages[startIdx].size = totalFreeSize
					m.pages[startIdx].allocated = fileSize
					m.pages[startIdx].fileName = fileName
					m.allocated += fileSize

					// Remover as páginas extras
					m.pages = append(m.pages[:startIdx+1], m.pages[startIdx+1+(i-startIdx):]...)
					return true
				}
			} else {
				totalFreeSize = 0
				startIdx = -1
			}
		}
	}

	pagesSize := 0
	for j := 0; j < len(m.pages); j++ {
		pagesSize += m.pages[j].size
	}

	if (m.maxSize - pagesSize) < fileSize {
		fmt.Println("There are no available pages to allocate the file.")
		return false
	}

	// Caso não tenha encontrado um espaço adequado, criar uma nova página.
	if m.maxSize-m.allocated >= fileSize {
		m.pages = append(m.pages, Page{size: fileSize, allocated: fileSize, fileName: fileName})
		m.allocated += fileSize
		return true
	}

	return false
}

func (m *Memory) DeallocateFile(fileName string) bool {
	for i := range m.pages {
		if m.pages[i].fileName == fileName {
			m.allocated -= m.pages[i].allocated
			m.pages[i].fileName = ""
			m.pages[i].allocated = 0
			return true
		}
	}
	return false
}

func createFile(fileName string, fileSize int) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	bytesWritten, err := file.Write(make([]byte, fileSize))
	if err != nil {
		return err
	}

	fmt.Printf("File created successfully with %d bytes written.\n", bytesWritten)
	return nil
}

func main() {
	createFile("file2", 2)
	createFile("file4", 4)
	createFile("file8", 8)

	memory := Memory{
		maxSize: 32,
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("1. Allocate File (First Fit)")
		fmt.Println("2. Allocate File (Worst Fit)")
		fmt.Println("3. Allocate File (Best Fit)")
		fmt.Println("4. Deallocate File")
		fmt.Println("5. View Memory")
		fmt.Println("6. Exit")
		fmt.Print("Enter your choice: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Enter file name: ")
			fileName, _ := reader.ReadString('\n')
			fileName = strings.TrimSpace(fileName)
			if memory.AllocateFileFirstFit(fileName) {
				fmt.Println("File allocated successfully.")
			} else {
				fmt.Println("Failed to allocate file.")
			}
		case "2":
			fmt.Print("Enter file name: ")
			fileName, _ := reader.ReadString('\n')
			fileName = strings.TrimSpace(fileName)
			if memory.AllocateFileWorstFit(fileName) {
				fmt.Println("File allocated successfully.")
			} else {
				fmt.Println("Failed to allocate file.")
			}
		case "3":
			fmt.Print("Enter file name: ")
			fileName, _ := reader.ReadString('\n')
			fileName = strings.TrimSpace(fileName)
			if memory.AllocateFileBestFit(fileName) {
				fmt.Println("File allocated successfully.")
			} else {
				fmt.Println("Failed to allocate file.")
			}
		case "4":
			fmt.Print("Enter file name: ")
			fileName, _ := reader.ReadString('\n')
			fileName = strings.TrimSpace(fileName)
			if memory.DeallocateFile(fileName) {
				fmt.Println("File deallocated successfully.")
			} else {
				fmt.Println("Failed to deallocate file.")
			}
		case "5":
			fmt.Printf("------***------\n")
			fmt.Println("Memory Status")
			fmt.Printf("Total Memory: %d bytes\n", memory.maxSize)
			fmt.Printf("Total Allocated: %d bytes\n", memory.allocated)
			fmt.Printf("Memory left: %d bytes\n", memory.maxSize-memory.allocated)
			fmt.Printf("Pages:\n")
			for i, page := range memory.pages {
				fmt.Printf("Page %d: File Name: %s, Size: %d bytes, Allocated: %d bytes\n", i+1, page.fileName, page.size, page.allocated)
			}
			fmt.Printf("------***------\n")
		case "6":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

		fmt.Println("")
	}
}
