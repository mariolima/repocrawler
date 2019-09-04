package utils


import(
	log "github.com/sirupsen/logrus"
	// "io/ioutil"
	"os"
)

func SaveLineToFile(line, file_name string) error {
	f, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		f, err = os.Create(file_name)
		if err!= nil {
			log.Fatal("Error Creating file: ",file_name)
		}
	}
	defer f.Close()
	_, err = f.WriteString(line)
	// err := ioutil.WriteFile(file_name, []byte(line), 0644)
	if err != nil {
		log.Fatal("Error writing to file: ",file_name," ", err)
	}
	return err
}

