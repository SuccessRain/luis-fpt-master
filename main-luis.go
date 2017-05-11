package main

import (
	"fmt"
	luis "github.com/SuccessRain/luis-fpt-sdk"
	"os"
	"bufio"
	"strings"
	"flag"
)

/*
	Created by CongVV
*/

func testLuis(e *luis.Luis, text string, intent string) bool {
	res, err := e.Predict(text)
	if err != nil {
		fmt.Println("Error happen on :", err.Err)
	}
	result := luis.GetBestScoreIntent(luis.NewPredictResponse(res))
	if strings.Compare(strings.TrimSpace(result.Name), strings.TrimSpace(intent)) == 0 {
		return true
	}
	return false
}

func TrainPost(e *luis.Luis){
	res, err := e.Train()
	if err != nil {
		fmt.Println("Error happen on :", err.Err)
	}
	fmt.Println("Got response:", string(res))
}

func TrainGet(e *luis.Luis){
	res, err := e.TrainStatus()
	if err != nil {
		fmt.Println("Error happen on :", err.Err)
	}
	fmt.Println("Got response:", string(res))
}

func main() {
	/*
	e := luis.NewLuis("2205967740474306901736ed5754e1ff", "38737645-cd67-4dd8-b4ba-ac2f3484eac3")

	TrainPost(e)
	TrainGet(e)

	res3, err := e.Predict("text")
	if err != nil {
		fmt.Println("Error happen on :", err.Err)
	}
	result := luis.GetBestScoreIntent(luis.NewPredictResponse(res3))
	//testLuis(e,"text","intent")
	fmt.Print("RESULTTTTTTT:\t"); fmt.Println(result)
	if strings.Compare(strings.TrimSpace(result.Name), strings.TrimSpace("intent")) == 0 {
		fmt.Print("OKKKK:\t"); fmt.Println(result.Name)
	}else{
		fmt.Print("NOT OKKKK:\t"); fmt.Println(result.Name)
	}
	*/

	trainCmd := flag.NewFlagSet("train", flag.ExitOnError)
	AddSharedFlags(trainCmd)

	testCmd := flag.NewFlagSet("test", flag.ExitOnError)
	AddSharedFlags(testCmd)

	if len(os.Args) < 2 {
		fmt.Println("Error: Input is not enough")
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	command := os.Args[1]
	if command == "train"{
		trainCmd.Parse(os.Args[2:])
	}else if command == "test" {
		testCmd.Parse(os.Args[2:])
	}else if command == "help" {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	if target != "intent" && target != "entity" {
		fmt.Println("Error: You must choose intent or entity")
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	if inputFP == "" {
		fmt.Println("Error: Input file is required but empty")
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	if key == "" {
		fmt.Println("Error: Key is required")
		fmt.Println(helpMessage)
		os.Exit(1)
	}
	if appid == "" {
		fmt.Println("Error: Appid is required")
		fmt.Println(helpMessage)
		os.Exit(1)
	}
	client := luis.NewClient(key)
	obj, err2 := ReadIntentsFromFile(inputFP)
	if err2 != nil{
		fmt.Println(err2)
	}

	if trainCmd.Parsed() {
		multi := recycleInentName(obj)
		var intents []luis.Intent
		var utterances []luis.Utterance
		for _, a := range multi {
			intents = append(intents, luis.Intent{a.intentName})
			for _,b := range a.values{
				if len(b) > 500{
					b = string(b[0:500])
				}
				utterances = append(utterances, luis.Utterance{b, a.intentName, []string{}})
			}
		}
		client.CreateApp(appid, intents, utterances)
	}

	if testCmd.Parsed() {
		var count int = 0
		var i int = 0
		e := luis.NewLuis(key, appid)
		for _,a := range obj{
			text := a.Name
			if len(text) > 500{
				text = string(text[0:500])
			}
			if testLuis(e, text, a.Intent){
				fmt.Print(i+1); fmt.Println(".\tCorrect")
				count ++
			}else{
				fmt.Print(i+1); fmt.Println(".\tIncorrect")
			}
			i++
		}
		fmt.Printf("Success: %f \n", float64(count) / float64(len(obj)) * 100)
	}
}

var target string
var inputFP string
var key string
var appid string

func AddSharedFlags(fs *flag.FlagSet) {
	fs.StringVar(&target, "t", "", "required, intent or entity")
	fs.StringVar(&inputFP, "i", "", "required, path to the input file")
	fs.StringVar(&key, "key", "", "required, LUIS.AI key")
	fs.StringVar(&appid, "appid", "", "required, LUIS.AI appid")
}

type ObjectFile struct{
	Intent string
	Name string
	Response string
}

func ReadIntentsFromFile(inputFP string) ([]ObjectFile, error) {
	input, err := os.Open(inputFP)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	var context []ObjectFile
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		var con ObjectFile
		text := strings.Replace(scanner.Text(), `"`, ``, -1)
		tokens := strings.SplitN(text, ",", 2)
		con.Intent, con.Name = strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])
		context = append(context, con)
	}

	return context, nil
}

type ObjectMutilValue struct{
	intentName string
	values []string
}

func recycleInentName(objs []ObjectFile)  []ObjectMutilValue {

	var omvs []ObjectMutilValue

	data := make(map[string][]string)
	for _, a := range objs {
		samples := data[a.Intent]
		samples = append(samples, a.Name)
		data[a.Intent] = samples
	}

	for key, value := range data {

		var omv ObjectMutilValue
		omv.intentName = key
		omv.values = value
		omvs = append(omvs, omv)
	}

	return omvs
}

const helpMessage string = `No help ahihihihihi`