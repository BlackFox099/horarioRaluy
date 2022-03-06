package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"math/rand"
	"time"
)

type Time struct {
	Day string `json:"day"`
	Hour string `json:"hour"`
}

type Materia struct {
	Name string `json:"name"`
	Time Time	`json:"time"`
}

type Materias struct {
	Materias []Materia `json:"materias"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func deleteFiles() {
	for range time.Tick(time.Second * 60){
		files , _ := ioutil.ReadDir("./tempFiles")
		if len(files) >= 10 {

			os.RemoveAll("./tempFiles")
			os.Mkdir("./tempFiles", os.ModePerm)

		}
	}
}

func main () {

	os.Mkdir("./tempFiles", os.ModePerm)
	
	go deleteFiles()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, DELETE",
	}))

	app.Get("dowloadJSON/:Id", func(c *fiber.Ctx) error {
		Id := c.Params("Id")

		File := "./tempFiles/calendario_" + Id + ".json"

		if _, err := os.Stat(File); os.IsNotExist(err) {
			return c.Status(500).JSON("Error: File not found")
		}

		return c.Download(File , "calendario.json")

	})

	app.Post("/createJSON", func(c *fiber.Ctx) error {

		materias := new(Materias)

		if err := c.BodyParser(materias); err != nil {
			log.Fatal(err)
		}

		file, _ := json.MarshalIndent(materias, "", " ")

		Id := RandStringBytes(16)

		fileName := "./tempFiles/calendario_"+ Id + ".json"

		err := os.WriteFile(fileName, file, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(Id)
	})

	app.Post("/uploadJSON" , func(c *fiber.Ctx) error {

		file, err := c.FormFile("document")

		if err != nil {
			log.Fatal(err)
		}

		fileName := "./tempFiles/" + file.Filename

		c.SaveFile(file, fmt.Sprintf("./%s", fileName))

		jsonFile, err := os.Open(fileName)

		if err != nil {
			log.Fatal(err)
		}

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var materias Materias

		json.Unmarshal(byteValue, &materias)

		jsonFile.Close()

		err = os.Remove(fileName)

		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(Materias{Materias: materias.Materias})

	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	err := app.Listen(":" + port)

	if err != nil {
		log.Fatal(err)
	}

}