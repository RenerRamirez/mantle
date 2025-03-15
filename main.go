package main

import (
  "fmt"
  "log"
  "os"
  // "io"
  "bufio"
  "os/exec"
  "strings"
  "encoding/json"
)

type GlossObj struct {
  Pos string `json:"pos"`
  Gloss string `json:"gloss"`
}

type WordJson struct {
  Reading string `json:"reading"`
  Text string `json:"text"`
  Kana string `json:"kana"`
  Score int `json:"score"`
  Seq int `json:"seq"`
  Gloss []GlossObj `json:"gloss"`
  Conj []interface{} `json:"conj"`
}

func main() {
  file, err := os.Open("../mantle/sentences.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    ex_sentence := scanner.Text()

    cmd := exec.Command("ichiran-cli", "-f", ex_sentence)
    var out strings.Builder
    cmd.Stdout = &out

    if err := cmd.Run(); err != nil {
      log.Fatal(err)
    }

    parsed_json := out.String()

    if isJson := json.Valid([]byte(parsed_json)); !isJson {
      fmt.Println("Invalid JSON!");
    }

    /*
    type Words struct {
      word string
      wordJson WordJson
      emptyArr []any
    }

    type InnerInnerArray struct {
      words []Words
    }

    type InnerArray struct {
      innerinnerArray []InnerInnerArray
      score int
      conj []string
    }

    type Array struct {
      innerArr []InnerArray
      period string
    }
    */

    const jsonStream = `
      [
        "hawai",
        {
          "reading": "ハワイ",
          "text": "ハワイ",
          "kana": "ハワイ",
          "score": 384,
          "seq": 1096400,
          "gloss": [
            {
              "pos": "[n]",
              "gloss": "Hawaii; Hawai'i"
            }
          ],
          "conj": []
        },
        []
      ]
`
    parseFinalArray(jsonStream)
    break
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }
}

func parseFinalArray(jsonStream string) {
  dec := json.NewDecoder(strings.NewReader(jsonStream))

  /* ["word",{},[]], ["word",{},[]], ... */
  var data []interface{}
  err := dec.Decode(&data)
  if err != nil {
    log.Fatalf("Failed to decode JSON: %v", err)
  }

  term, ok := data[0].(string) 
  if !ok {
    log.Fatalf("Expected string at data[0]: %v", err)
  }
  _ = term

  objMap, ok := data[1].(map[string]interface{})
  if !ok {
    log.Fatalf("Expected obj at data[1]: %v", err)
  }

  var myword WordJson
  wordbytes, err := json.Marshal(objMap)
  if err != nil {
    log.Fatalf("Failed to marshal %v", err)
  }

  err = json.Unmarshal(wordbytes, &myword)
  if err != nil {
    log.Fatalf("Failed to unmarshal %v", err)
  }

  emptyArr, ok := data[2].([]interface{})
  if !ok {
    log.Fatalf("Expected empty[] at data[3]: %v", err)
  }
  _ = emptyArr

  fmt.Println(myword.Seq)
  /* ["word",{},[]], ["word",{},[]], ... */
}

