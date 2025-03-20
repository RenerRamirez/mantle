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
  Compound []string `json:"compound,omitempty"`
  Components []interface{} `json:"components,omitempty"`
  Conj []interface{} `json:"conj"`
}

func main() {
  file, err := os.Open("sentences.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)

  cnt := 0
  _ = cnt
  for scanner.Scan() {
    ex_sentence := scanner.Text()

    cmd := exec.Command("ichiran-cli", "-f", ex_sentence)
    fmt.Println("【",ex_sentence,"】")
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

    parseFinalArray(parsed_json)
    // if cnt == 5 {
    //   break
    // }
    // cnt++
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }
}

func parseFinalArray(jsonStream string) {
  dec := json.NewDecoder(strings.NewReader(jsonStream))

  /* ["word",{},[]], ["word",{},[]], ... */
  var sentencedot []interface{}
  err := dec.Decode(&sentencedot)
  if err != nil {
    log.Fatalf("Failed to decode JSON: %v", err)
  }
  
  template  := sentencedot[0] .([]interface{}) 
  sentences := template[0]    .([]interface{})

  /* ["word",{},[]], ["word",{},[]], ... */
  data := sentences[0].([]interface{})

  for k, _ := range data {
    wordinfo, ok := data[k].([]interface{})
    if !ok {
       log.Fatalf("Expected wordinfo at data[0]: %v", err)
    }

    term, ok := wordinfo[0].(string) 
    if !ok {
      log.Fatalf("Expected string at data[0]: %v", err)
    }

    objMap, ok := wordinfo[1].(map[string]interface{})
    if !ok {
      log.Fatalf("Expected obj at data[1]: %v", err)
    }

    wordObj := objMap
    alt, ok := objMap["alternative"].([]interface{})
    if ok {
      // only check the first entry in "alternative"
      a, ok := alt[0].(map[string]interface{})
      if ok {
        wordObj = a
      } // !ok should be unreachable
    }

    var word WordJson
    wordBytes, err := json.Marshal(wordObj) // objMap is the same as 'a' in else
    if err != nil {
      log.Fatalf("Failed to marshal %v", err)
    }

    err = json.Unmarshal(wordBytes, &word)
    if err != nil {
      log.Fatalf("Failed to unmarshal %v", err)
    }

    emptyArr, ok := wordinfo[2].([]interface{})
    if !ok {
      log.Fatalf("Expected empty[] at wordinfo[2]: %v", err)
    }
    _ = emptyArr
    _ = term

    fmt.Println(word.Text)
  }
  /* ["word",{},[]], ["word",{},[]], ... */
}

