package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

type Response struct{
	Page int `json:"page"`
	Per_Page int `json:"per_page"`
	Total int `json:"total"`
	Total_Pages int `json:"total_pages"`
	Data []CityDetails `json:"data"`
}

type CityDetails struct{
	Name string `json:"name"`
	Weather string `json:"weather"`
	Status []string `json:"status"`
}

func ExtractNumber(input string) int {
 
    pattern := `\d+`

    re := regexp.MustCompile(pattern)

    match := re.FindString(input)

    number, err := strconv.Atoi(match)
    if err != nil {
       fmt.Println(err.Error())
    }

    return number
}

func getResponse(url string) (Response, error) {
	resp, err := http.Get(url)

	if err != nil {
		return Response{},err
	}

	var response Response
	body,err := io.ReadAll(resp.Body)

	if err != nil {
		return Response{},err
	}

	err = json.Unmarshal(body,&response)

	if err != nil {
		return Response{},err
	}

	return response,nil

}

func fetch(url string) (interface{}, error) {
	response, err := getResponse(url)

	if err != nil {
		return nil,err
	}

	searchQuery := "&page="
	url = url + searchQuery
	
	var out []interface{}
	for i:= 1; i <= response.Total_Pages; i++ {

		searchUrl := url + strconv.Itoa(i)

		res, err := getResponse(searchUrl)

		if err != nil {
			return nil,err
		}
	
		for j := 0 ; j < len(res.Data); j++ {
			var array []interface{}
			array = append(array, res.Data[j].Name)
			array = append(array, ExtractNumber(res.Data[j].Weather))
			array = append(array, ExtractNumber(res.Data[j].Status[0]))
			array = append(array, ExtractNumber(res.Data[j].Status[1]))	
			out = append(out, array)
		}

	}
	
	return out, nil
}


func main() {
	
	url := "https://jsonmock.hackerrank.com/api/weather/search?name="

	var name string
	fmt.Println("Enter the name to search ::")
	fmt.Scanln(&name)

	url = url + name
	fmt.Println(url)
	out, err := fetch(url)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("//////////////////////////////////////////////")
	fmt.Println(out)
}