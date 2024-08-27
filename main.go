package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func splitString(input string) []string {
	// разделение текста на массив строк
	var result []string
	regex := regexp.MustCompile(`\b\w*-\w*\w*\b|\b\w*'\w*\w*\b|\(\w+\)|\(\w+,\s*\w+\)|\b\w+\b|[^\w\s]|\n`)
	result = regex.FindAllString(input, -1)
	return result
}

func toHex(words []string, index, count int) []string {
	if count == 0 {
		return append(words[:index], words[index+1:]...)
	}
	// Проверяем, находится ли индекс в пределах массива
	if index == 0 && count > 1 {
		words = append(words[:index], words[index+1:]...)
		index = index - 1
	} else if index < 1 {
		return append(words[:index], words[index+1:]...)

	}
	if index < count {
		count = index
	}
	if strings.ContainsAny(words[index-count], "\n") {
		count++
	}
	// проверяем текущее слово, является ли знаком препинания
	if strings.ContainsAny(words[index-count], "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") {
		count++
	}
	// Применяем команду к словам перед текущим индексом
	for i := count; i >= count; count-- {
		if count == 0 {
			break
		}
		if count >= 0 && count < len(words) { // Проверяем, находится ли индекс в пределах массива
			current := words[index-count]
			decimal, err := strconv.ParseInt(current, 16, 64)
			if err != nil {
				return append(words[:index], words[index+1:]...)
			} else {
				str := strconv.FormatInt(decimal, 10)
				words[index-count] = str
				return append(words[:index], words[index+1:]...)
			}
		}
	}
	return append(words[:index], words[index+1:]...)
}

func toBin(words []string, index, count int) []string {
	if count == 0 {
		return append(words[:index], words[index+1:]...)
	}
	// Проверяем, находится ли индекс в пределах массива
	if index == 0 && count > 1 {
		words = append(words[:index], words[index+1:]...)
		index = index - 1
	} else if index < 1 {
		return append(words[:index], words[index+1:]...)

	}
	if index < count {
		count = index
	}
	// Проверяем, находится ли индекс в пределах массива
	if count <= 0 {
		return words
	}
	if strings.ContainsAny(words[index-count], "\n") {
		count++
	}
	// проверяем текущее слово, является ли знаком препинания
	if strings.ContainsAny(words[index-count], "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") {
		count++
	}
	// Применяем команду к словам перед текущим индексом
	for i := count; i >= count; count-- {
		if count == 0 {
			break
		}
		if count >= 0 && count < len(words) { // Проверяем, находится ли индекс в пределах массива
			current := words[index-count]
			decimal, err := strconv.ParseInt(current, 2, 64)
			if err != nil {
				return append(words[:index], words[index+1:]...)
			} else {
				str := strconv.FormatInt(decimal, 10)
				words[index-count] = str
				return append(words[:index], words[index+1:]...)
			}
		}
	}
	return append(words[:index], words[index+1:]...)
}

func IsVowel(s string) bool {
	//гласные буквы как словарь
	vowels := "aeiouAEIOU"
	return strings.ContainsAny(s, vowels)
}

func removeSpacesInsideQuotes(input string) string {
	lines := strings.Split(input, "\n")
	var processedLines []string

	// Обработка каждой строки
	for _, line := range lines {
		// Пропускаем строки, содержащие двойные скобки
		if strings.Contains(line, "(") && strings.Contains(line, ")") {
			processedLines = append(processedLines, line)
			continue
		}

		// Убираем пробелы в начале и в конце строки
		line = strings.TrimSpace(line)

		// Убираем пробелы внутри кавычек на текущей строке
		re := regexp.MustCompile(`(\s*)'([^']*?)'(\s*)`)
		line = re.ReplaceAllString(line, "$1'$2'$3")

		// Добавляем обработанную строку в список
		processedLines = append(processedLines, line)
	}

	// Объединяем строки с добавлением переходов на новую строку
	result := strings.Join(processedLines, "\n")

	// Удаляем пробелы перед знаками препинания
	punctuation := regexp.MustCompile(`\s*([.,!?:;]+)[^\S\n]*(\s*)`)
	result = punctuation.ReplaceAllString(result, "$1 $2")

	// Удаляем пробелы между апострофами
	result = regexp.MustCompile(`(\w)'(\w)`).ReplaceAllString(result, "$1\x00$2")

	// Удаляем пробелы внутри кавычек
	result = regexp.MustCompile(`\s*'([^']*?)'\s*`).ReplaceAllString(result, "'$1'")

	// Заменяем одинарные кавычки на " ' " чтобы отделить их от слов
	result = regexp.MustCompile(`'`).ReplaceAllString(result, " ' ")

	// Удаляем пробелы перед и после кавычек
	result = regexp.MustCompile(`(\ +)?(')(\ +)?(.+?)(\ +)?(')(\ +)?`).ReplaceAllString(result, ` $2$4$6 `)

	// Восстанавливаем апострофы
	result = regexp.MustCompile(`\x00`).ReplaceAllString(result, "'")

	// Удаляем пробелы в начале и в конце строки
	result = strings.TrimSpace(result)

	return result
}

func correct(words []string) string {
	var result strings.Builder

	for i, word := range words {
		// проверяем текущее слово, является ли знаком препинания
		isPunctuation := unicode.IsPunct([]rune(word)[len(word)-1])

		if i > 0 && !isPunctuation {
			result.WriteString(" ")
		}
		// исправляем артикли перед словами, начинающимися с гласной или 'h'
		if strings.ToLower(word) == "a" && i < len(words)-1 && IsVowel(words[i+1]) && len(words[i+1]) > 1 {
			if word == "A" {
				result.WriteString("An")
			} else {
				result.WriteString("an")
			}
		} else {
			// добавляем текущее слово
			result.WriteString(word)
		}
	}
	return result.String()
}

func capitalize(words []string, index, count int) []string {
	if count == 0 {
		return append(words[:index], words[index+1:]...)
	}
	// Проверяем, находится ли индекс в пределах массива
	if index == 0 && count > 1 {
		words = append(words[:index], words[index+1:]...)
		index = index - 1
	} else if index < 1 {
		return append(words[:index], words[index+1:]...)

	}
	if index < count {
		count = index
	}
	// Проверяем, находится ли индекс в пределах массива
	if count <= 0 {
		return words
	}
	if strings.ContainsAny(words[index-count], "\n") {
		count++
	}
	// проверяем текущее слово, является ли знаком препинания
	if strings.ContainsAny(words[index-count], "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") {
		count++
	}
	// Применяем команду к словам перед текущим индексом
	for i := count; i >= count; count-- {
		if count == 0 {
			break
		}
		if count >= 0 && count < len(words) { // Проверяем, находится ли индекс в пределах массива
			s := words[index-count]
			capitalized := strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
			words[index-count] = capitalized
		}
	}
	return append(words[:index], words[index+1:]...)
}

func ToUpper(words []string, index, count int) []string {
	// Проверяем, находится ли индекс в пределах массива
	if index == 0 && count > 1 {
		words = append(words[:index], words[index+1:]...)
		index = index - 1
	} else if index < 1 {
		return append(words[:index], words[index+1:]...)

	}
	if index < count {
		count = index
	}
	// Проверяем, находится ли индекс в пределах массива
	if count <= 0 {
		return words
	}
	if strings.ContainsAny(words[index-count], "\n") {
		count++
	}
	// проверяем текущее слово, является ли знаком препинания
	if strings.ContainsAny(words[index-count], "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") {
		count++
	}
	// Применяем команду к словам перед текущим индексом
	for i := count; i >= count; count-- {
		if count == 0 {
			break
		}
		if count >= 0 && count < len(words) { // Проверяем, находится ли индекс в пределах массива
			s := words[index-count]
			capitalized := strings.ToUpper(string(s[0:]))
			words[index-count] = capitalized
		}
	}
	return append(words[:index], words[index+1:]...)
}

func ToLower(words []string, index, count int) []string {
	// Проверяем, находится ли индекс в пределах массива
	if index == 0 && count > 1 {
		words = append(words[:index], words[index+1:]...)
		index = index - 1
	} else if index < 1 {
		return append(words[:index], words[index+1:]...)

	}
	if index < count {
		count = index
	}
	// Проверяем, находится ли индекс в пределах массива
	if count <= 0 {
		return words
	}
	if strings.ContainsAny(words[index-count], "\n") {
		count++
	}
	// проверяем текущее слово, является ли знаком препинания
	if strings.ContainsAny(words[index-count], "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") {
		count++
	}
	// Применяем команду к словам перед текущим индексом
	for i := count; i >= count; count-- {
		if count == 0 {
			break
		}
		if count >= 0 && count < len(words) { // Проверяем, находится ли индекс в пределах массива
			s := words[index-count]
			capitalized := strings.ToLower(s[0:])
			words[index-count] = capitalized
		}
	}
	return append(words[:index], words[index+1:]...)
}

func applyCommands(words []string) []string {
	var count int
	var command string

	for index := 0; index < len(words); index++ {
		for index := 0; index < len(words); index++ {
			// Проверяем, есть ли команды в текущем слове
			if strings.HasPrefix(words[index], "(") && strings.HasSuffix(words[index], ")") {
				parts := strings.Split(words[index][1:len(words[index])-1], ",")
				command = strings.TrimSpace(parts[0])
				count = 1
				if len(parts) == 2 {
					count, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
				}
				switch command {
				case "hex":
					words = toHex(words, index, count)
				case "bin":
					words = toBin(words, index, count)
				case "up":
					words = ToUpper(words, index, count)
				case "low":
					words = ToLower(words, index, count)
				case "cap":
					words = capitalize(words, index, count)
				}
			}
		}
	}
	return words
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <source_file> <destination_file>")
		os.Exit(1)
	}

	srcFile := os.Args[1]  // Исходный файл
	destFile := os.Args[2] // Целевой файл

	file, err := os.ReadFile(srcFile)
	if err != nil {
		fmt.Println("Ошибка!\n", err)
		return
	}

	// Прочтенный текст в одну строку
	text := string(file)
	words := splitString(text)
	preresult := applyCommands(words)
	cor := correct(preresult)
	res := removeSpacesInsideQuotes(cor)

	// Запись обработанного текста в файл
	err = os.WriteFile(destFile, []byte(res), 0644)
	if err != nil {
		fmt.Println("Ошибка при записи файла:", err)
		return
	}

	fmt.Println("Обработка завершена. Результат записан в файл", destFile)
}
