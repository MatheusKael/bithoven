package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func Load(path string) string {

	path, err := disc(path)

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	godotenv.Load(path)
	key := os.Getenv("GEMINI_KEY")
	url := os.Getenv("GEMINI_URL")

	if len(url) == 0 {
		fmt.Print("empty url env")
		os.Exit(1)
	}

	return fmt.Sprintf("%s%s", url, key)
}

// disc procura um arquivo .env subindo a árvore de diretórios
// a partir de dir. Retorna o diretório que contém o .env.
func disc(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		envPath := filepath.Join(abs, ".env")

		if info, err := os.Stat(envPath); err == nil && !info.IsDir() {
			return envPath, nil
		}

		parent := filepath.Dir(abs)
		if parent == abs { // chegamos na raiz
			return "", errors.New(".env não encontrado em nenhum diretório pai")
		}
		abs = parent
	}
}
