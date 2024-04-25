package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	PingCommand string = "ping"
	EchoCommand string = "echo"
)

type Message struct {
	command   string
	arguments []interface{}
}

func getItemEndingIndex(item *[]byte) int {
	i := 1
	for (*item)[i] != '\r' {
		i++
	}
	return i
}

func deserializeString(stringBytes []byte) (string, int) {
	stringEnding := getItemEndingIndex(&stringBytes)
	return string(stringBytes[1:stringEnding]), stringEnding + 2
}

func deserializeError(errorBytes []byte) (error, int) {
	errorEnding := getItemEndingIndex(&errorBytes)
	return fmt.Errorf("%s", errorBytes[1:errorEnding]), errorEnding + 2
}

func deserializeInt(intBytes []byte) (integer int, intLength int) {
	intEnding := getItemEndingIndex(&intBytes)
	integer, err := strconv.Atoi(string(intBytes[1:intEnding]))
	if err != nil {
		fmt.Println("Error parsing message, cannot parse int")
		os.Exit(1)
	}
	return integer, intEnding + 2
}

func deserializeBulkString(stringBytes []byte) (string, int) {
	stringLengthEnding := getItemEndingIndex(&stringBytes)
	stringLength, err := strconv.Atoi(string(stringBytes[1:stringLengthEnding]))
	if err != nil {
		fmt.Println("Error parsing message, cannot parse bulk string length")
		os.Exit(1)
	}
	stringStart := stringLengthEnding + 2
	return string(stringBytes[stringStart : stringStart+stringLength]), stringStart + stringLength + 2
}

func deserializeArray(arrayBytes []byte) (array []interface{}, arrayLength int) {
	currentElementStartIndex := 2 // First char is *, second is first digit of items amount
	for arrayBytes[currentElementStartIndex] != '\r' {
		currentElementStartIndex++
	}
	arrayLength, err := strconv.Atoi(string(arrayBytes[:currentElementStartIndex]))
	if err != nil {
		fmt.Println("An error occurred parsing array, cannot parse array length")
		os.Exit(1)
	}

	array = make([]interface{}, arrayLength)
	currentElementStartIndex += 2 // Point at the start of the first array element
	for i := 0; i < arrayLength; i++ {
		item, itemLength := deserializeItem(arrayBytes[currentElementStartIndex:])
		array[i] = item
		currentElementStartIndex += itemLength
	}
	return array, currentElementStartIndex
}

func deserializeItem(item []byte) (interface{}, int) {
	switch item[0] {
	case '+':
		return deserializeString(item)
	case '-':
		return deserializeError(item)
	case ':':
		return deserializeInt(item)
	case '$':
		return deserializeBulkString(item)
	default:
		return nil, 0
	}
}

func ParseMessage(message []byte) Message {
	messageArray, _ := deserializeArray(message)
	return Message{
		command:   messageArray[0].(string),
		arguments: messageArray[1:],
	}
}
