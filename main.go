package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func main() {
	// Define as flags de linha de comando
	linkFlag := flag.String("u", "", "URL para baixar")
	listFlag := flag.String("l", "", "Arquivo contendo uma lista de URLs")

	flag.Parse()

	// Verifica se a flag -u ou -l foi especificada
	if *linkFlag == "" && *listFlag == "" {
		fmt.Println("Uso: download_and_extract -u URL | -l arquivo.txt")
		os.Exit(1)
	}

	// Se a flag -u foi especificada, baixe a página da URL
	if *linkFlag != "" {
		downloadAndExtract(*linkFlag)
	}

	// Se a flag -l foi especificada, leia a lista de URLs do arquivo e baixe as páginas
	if *listFlag != "" {
		urls, err := readURLsFromFile(*listFlag)
		if err != nil {
			fmt.Println("Erro ao ler o arquivo de URLs:", err)
			os.Exit(1)
		}

		for _, url := range urls {
			downloadAndExtract(url)
		}
	}
}

func downloadAndExtract(url string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// Aplicar expressões regulares para encontrar links de mídias sociais
	twitterRegex := regexp.MustCompile(`https?://twitter\.com/\w+`)
	facebookRegex := regexp.MustCompile(`https?://(www\.)?facebook\.com/\w+`)

	twitterLinks := twitterRegex.FindAllString(string(body), -1)
	facebookLinks := facebookRegex.FindAllString(string(body), -1)

	if len(twitterLinks) > 0 {
		fmt.Printf("Links do Twitter encontrados em %s:\n", url)
		for _, link := range twitterLinks {
			color.Set(color.FgGreen)
			fmt.Println(link)
			color.Unset()
		}
	}

	if len(facebookLinks) > 0 {
		fmt.Printf("Links do Facebook encontrados em %s:\n", url)
		for _, link := range facebookLinks {
			color.Set(color.FgGreen)
			fmt.Println(link)
			color.Unset()
		}
	}
}

func readURLsFromFile(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			urls = append(urls, line)
		}
	}

	return urls, nil
}

