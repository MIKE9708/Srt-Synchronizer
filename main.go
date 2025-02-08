package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func get_file(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func generate_offset(path string, offset float64) error {
	path_array := strings.Split(path, "/")
	offset_file := path_array[len(path_array)-1]
	fmt.Printf("Creating file %s ...\n", offset_file)

	first_pattern, err := regexp.Compile(`(\d+:\d+:\d+.\d+).*--`)
	second_pattern, err := regexp.Compile(`> (\d+:\d+:\d+.\d+)`)

	of, err := os.Create(fmt.Sprint("sync_subs/", offset_file))
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		line_dotted := strings.Replace(line, ",", ".", 2)
		if first_pattern.MatchString(line_dotted) {
			offset1, _ := time.Parse("15:04:05.999", first_pattern.FindStringSubmatch(line_dotted)[1])
			offset1 = offset1.Add(time.Duration(offset * float64(time.Second)))
			offset1_string := fmt.Sprintf("%02d:%02d:%02d,%d", offset1.Hour(), offset1.Minute(), offset1.Second(), offset1.Nanosecond()/1e6)

			offset2, _ := time.Parse("15:04:05", second_pattern.FindStringSubmatch(line_dotted)[1])
			offset2 = offset2.Add(time.Duration(offset * float64(time.Second)))
			offset2_string := fmt.Sprintf("%02d:%02d:%02d,%d", offset2.Hour(), offset2.Minute(), offset2.Second(), offset2.Nanosecond()/1e6)

			updated_line := fmt.Sprintf("%s --> %s\n", offset1_string, offset2_string)
			updated_line = strings.Replace(updated_line, ".", ",", 2)

			of.Write([]byte(updated_line))
		} else {
			of.Write([]byte(line + "\n"))
		}
	}
	return nil
}

func main() {
	pflag := flag.String("p", "", "Specify the path of the srt file")
	offsetflag := flag.Float64("o", 0, "Specify the offset")

	flag.Parse()

	if *pflag == "" {
		fmt.Println("---------HELPER----------\n \t -p Specify Specify the path of the srt file")
		return
	}

	if *offsetflag == 0 {
		fmt.Println("No offset specified")
		return
	}

	_, err := get_file(*pflag)
	if err != nil {
		log.Fatal("Error for the specified path: \n", err)
	}

	err = generate_offset(*pflag, *offsetflag)
	if err != nil {
		log.Fatal("Error creating the offset\n")
	}

}
