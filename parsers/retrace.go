package parsers

import (
	"encoding/json"
	"fmt"
	"log"
)

type RetraceData struct {
	Parameters                map[string]interface{}  `json:"parameters"`
	Shaders                   map[string]interface{}  `json:"shaders"`
	Uniforms                  map[string]interface{}  `json:"uniforms"`
	Buffers                   map[string]interface{}  `json:"buffers"`
	ShaderStorageBufferBlocks map[string]interface{}  `json:"shaderstoragebufferblocks"`
	Textures                  map[string]interface{}  `json:"textures"`
	FrameBufer                map[string]*FrameBuffer `json:"framebuffer"`
}

type FrameBuffer struct {
	Class  string `json:"__class__"`
	Width  int    `json:"__width__"`
	Height int    `json:"__height__"`
	Depth  int    `json:"__depth__"`
	Format string `json:"__format__"`
	Data   string `json:"__data__"`
}

func NewRetraceData(retraceOutput string) RetraceData {

	var rd RetraceData

	if err := json.Unmarshal([]byte(retraceOutput), &rd); err != nil {
		log.Println("CANT DECODE")
		log.Println(err.Error())
	}

	fmt.Println("PRINTING BUFFERS")
	fmt.Printf("%+v\n", rd.Buffers)
	fmt.Println("PRINTING UNIFORMS")
	fmt.Printf("%+v\n", rd.Uniforms)

	return rd
}
