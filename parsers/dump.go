package parsers

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

type TraceDump struct {
	Frames []*Frame `json:"frames"`
}

type Frame struct {
	ID    int     `json:"id"`
	Calls []*Call `json:"calls"`
}

type Call struct {
	ID           string   `json:"id"`
	FunctionName string   `json:"functionName"`
	ParamNames   []string `json:"paramNames"`
	ParamValues  []string `json:"paramValues"`
	ReturnValue  string   `json:"returnValue"`
}

type ImageSet struct {
	Stencil string   `json:"stencil"`
	Depth   string   `json:"depth"`
	MRT     []string `json:"mrt"`
	Type    string   `json:"type"`
	CallID  string   `json:"callID"`
}

func ParseDump(dumpContents string) *TraceDump {

	td := new(TraceDump)

	frameNumber := 0

	scanner := bufio.NewScanner(strings.NewReader(dumpContents))

	frame := new(Frame)
	frame.ID = frameNumber

	shaderSource := false

	shaderSourceIndex := -1

	var str strings.Builder

	var c *Call

	for scanner.Scan() {
		line := scanner.Text()

		if !shaderSource {
			// this is probably the beginning of a shader source declaration
			if strings.Contains(line, "#") && strings.Contains(line, "version") {

				c = new(Call)
				// example line for this scenario
				// 193 glShaderSource(shader = 1, count = 1, string = &"#version 330

				shaderSource = true

				// 193
				callID := strings.Fields(line)[0]

				lineComponents := strings.Split(line, "(")

				// 193 glShaderSource
				callIDAndFunctionName := lineComponents[0]

				// shader = 1, count = 1, string = &"#version 330
				someParametersAndReturn := lineComponents[1]

				someParameters := strings.Split(someParametersAndReturn, ")")

				returnValue := ""

				if len(someParameters) > 1 {
					returnValue = someParameters[1]
				}
				// glShaderSource
				functionName := strings.Split(callIDAndFunctionName, " ")[1]

				parameterPairs := strings.Split(someParameters[0], " = ")

				// "attribs", {kCGLPFAAccelerated, kCGLPFAClosestPolicy, kCGLPFAOpenGLProfile, kCGLOGLPVersion_3_2_Core, kCGLPFAColorSize, 24, kCGLPFAAlphaSize, 8, kCGLPFADepthSize, 24, kCGLPFAStencilSize, 8, kCGLPFADoubleBuffer, kCGLPFASampleBuffers, 0, 0}, pix

				// first param name
				c.ParamNames = append(c.ParamNames, strings.TrimSpace(parameterPairs[0]))

				var i int

				for i = 1; i < len(parameterPairs); i++ {

					p := parameterPairs[i]

					nv := strings.Split(p, ", ")

					//n := nv[0]
					//
					//v := nv[1]

					if len(nv) > 1 {
						previousValue := strings.Join(nv[0:len(nv)-1], " ")
						nextName := nv[len(nv)-1]

						c.ParamNames = append(c.ParamNames, strings.TrimSpace(nextName))
						c.ParamValues = append(c.ParamValues, previousValue)
					}

				}

				// last value
				c.ParamValues = append(c.ParamValues, parameterPairs[len(parameterPairs)-1])

				if strings.Contains(parameterPairs[len(parameterPairs)-1], "#") && strings.Contains(parameterPairs[len(parameterPairs)-1], "version") {
					shaderSourceIndex = i
				}

				c.ID = callID
				c.FunctionName = functionName
				c.ReturnValue = returnValue

			} else {
				// regular line, nothing fancy

				line = strings.TrimSpace(line)
				//fmt.Println(line)

				if strings.HasPrefix(line, "//") {
					log.Println("Comment, skipping")
					log.Println(line)
					continue
				}

				// if it's an empty line, new frame!
				if len(line) == 0 {
					// add the current frame to the trace dump array of frames
					td.Frames = append(td.Frames, frame)

					frame = new(Frame)
					frameNumber++

					frame.ID = frameNumber

				} else {
					c = new(Call)

					// 209
					callID := strings.Fields(line)[0]

					lineComponents := strings.Split(line, "(")

					//fmt.Println(lineComponents)

					if len(lineComponents) < 2 {
						log.Println("Skipping this line (, possibly malformed")
						log.Println(line)
						continue
					}

					// 209 CGLFlushDrawable
					callIDAndFunctionName := lineComponents[0]

					// ctx = 0x805a200) = kCGLNoError
					someParametersAndReturn := lineComponents[1]

					someParameters := strings.Split(someParametersAndReturn, ")")

					if len(someParameters) < 2 {
						log.Println("Skipping this line ), possibly malformed")
						log.Println(line)
						continue
					}

					//  = kCGLNoError
					returnValue := someParameters[1]

					// CGLFlushDrawable
					functionName := strings.Split(callIDAndFunctionName, " ")[1]

					// attribs = {kCGLPFAAccelerated, kCGLPFAClosestPolicy, kCGLPFAOpenGLProfile, kCGLOGLPVersion_3_2_Core, kCGLPFAColorSize, 24, kCGLPFAAlphaSize, 8, kCGLPFADepthSize, 24, kCGLPFAStencilSize, 8, kCGLPFADoubleBuffer, kCGLPFASampleBuffers, 0, 0}, pix = &0x4700020, npix = &2
					parameterPairs := strings.Split(someParameters[0], " = ")

					// "attribs", {kCGLPFAAccelerated, kCGLPFAClosestPolicy, kCGLPFAOpenGLProfile, kCGLOGLPVersion_3_2_Core, kCGLPFAColorSize, 24, kCGLPFAAlphaSize, 8, kCGLPFADepthSize, 24, kCGLPFAStencilSize, 8, kCGLPFADoubleBuffer, kCGLPFASampleBuffers, 0, 0}, pix

					// first param name
					c.ParamNames = append(c.ParamNames, strings.TrimSpace(parameterPairs[0]))

					var i int

					for i = 1; i < len(parameterPairs); i++ {

						p := parameterPairs[i]

						nv := strings.Split(p, ", ")

						//n := nv[0]
						//
						//v := nv[1]

						if len(nv) > 1 {
							previousValue := strings.Join(nv[0:len(nv)-1], " ")
							nextName := nv[len(nv)-1]

							c.ParamNames = append(c.ParamNames, strings.TrimSpace(nextName))
							c.ParamValues = append(c.ParamValues, previousValue)
						}

					}

					// last value
					c.ParamValues = append(c.ParamValues, parameterPairs[len(parameterPairs)-1])

					if strings.Contains(parameterPairs[len(parameterPairs)-1], "#") && strings.Contains(parameterPairs[len(parameterPairs)-1], "version") {
						shaderSourceIndex = i
					}

					c.ID = callID
					c.FunctionName = functionName
					c.ReturnValue = returnValue

					for i := 0; i < len(c.ParamNames); i++ {
						if len(c.ParamNames[i]) == 0 {
							c.ParamNames = append(c.ParamNames[:i], c.ParamNames[i+1:]...)
							c.ParamValues = append(c.ParamValues[:i], c.ParamValues[i+1:]...)
						}
					}

					frame.Calls = append(frame.Calls, c)
				}
			}
		} else {
			// inside a shader source

			//fmt.Println(line)

			if strings.Contains(line, "\", length =") && strings.Contains(line, ")") {

				// }", length = NULL, something = TEST)
				removedParens := strings.Split(line, ")")

				// }", length = NULL
				remainingParameters := strings.Split(removedParens[0], ",")
				// "}", length", "NULL"

				str.WriteString(remainingParameters[0])

				outstandingParameters := remainingParameters[1]

				//fmt.Println(len(c.ParamValues))
				//fmt.Println(c.ParamValues)
				//fmt.Println(shaderSourceIndex)
				c.ParamValues[shaderSourceIndex-2] = fmt.Sprintf("%s\\n%s", c.ParamValues[shaderSourceIndex-2], str.String())

				parameterPairs := strings.Split(outstandingParameters, " = ")

				// first param name
				c.ParamNames = append(c.ParamNames, strings.TrimSpace(parameterPairs[0]))

				var i int

				for i = 1; i < len(parameterPairs); i++ {

					p := parameterPairs[i]

					nv := strings.Split(p, ", ")

					if len(nv) > 1 {
						previousValue := strings.Join(nv[0:len(nv)-1], " ")
						nextName := nv[len(nv)-1]

						c.ParamNames = append(c.ParamNames, strings.TrimSpace(nextName))
						c.ParamValues = append(c.ParamValues, previousValue)
					}

				}

				// last value
				c.ParamValues = append(c.ParamValues, parameterPairs[len(parameterPairs)-1])

				shaderSource = false

				for i := 0; i < len(c.ParamNames); i++ {
					if len(c.ParamNames[i]) == 0 {
						c.ParamNames = append(c.ParamNames[:i], c.ParamNames[i+1:]...)
						c.ParamValues = append(c.ParamValues[:i], c.ParamValues[i+1:]...)
					}
				}

				frame.Calls = append(frame.Calls, c)
			} else {
				str.WriteString(fmt.Sprintf("%s\\n", line))
			}
		}

	}

	return td
}

func ParseImageDumpFile(imageDump string) *ImageSet {

	scanner := bufio.NewScanner(strings.NewReader(imageDump))

	imageSet := new(ImageSet)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Wrote") {
			filenameComponents := strings.Split(line, fmt.Sprintf("Wrote "))
			fileName := filenameComponents[1]

			// main.0000003635-mrt0.png
			fileIDAndType := strings.Split(fileName, "-")

			// 0000003635
			fileID := fileIDAndType[0]

			// mrt0.png
			imageTypeAndExtension := fileIDAndType[1]

			imageTypeAndExtensionComponents := strings.Split(imageTypeAndExtension, ".")

			// mrt0
			imageType := imageTypeAndExtensionComponents[0]

			// png
			imageExtension := imageTypeAndExtensionComponents[1]

			imageSet.CallID = fileID
			addImageToImageSet(imageSet, imageExtension, imageType)
			// add this line to the image set

		}

	}

	return imageSet
}

func addImageToImageSet(is *ImageSet, imageExtension, imageType string) {
	is.Type = imageExtension

	if strings.Contains(imageType, "mrt") {
		is.MRT = append(is.MRT, imageType)
	} else if strings.Contains(imageType, "z") {
		is.Depth = imageType
	} else if strings.Contains(imageType, "s") {
		is.Stencil = imageType
	} else {
		log.Printf("Unknown image type %s", imageType)
	}

}
