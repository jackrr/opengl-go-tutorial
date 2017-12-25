package shader

import (
	"fmt"
  "strings"
  "io/ioutil"

  "github.com/go-gl/gl/v3.3-core/gl"
)

const NULL_TERM = "\x00"

type Shader struct {
  ID uint32
}

func compileShader(source string, shaderType uint32) (uint32, error) {
  shader := gl.CreateShader(shaderType)
  csources, free := gl.Strs(source)

  gl.ShaderSource(shader, 1, csources, nil)
  free()
  gl.CompileShader(shader)

  var status int32
  gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
  if status == gl.FALSE {
    var logLength int32
    gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

    log := strings.Repeat(NULL_TERM, int(logLength+1))
    gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

    return 0, fmt.Errorf("failed to compile %v: %v", source, log)
  }

  return shader, nil
}

func configureProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
  vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
  if err != nil {
    return 0, err
  }

  fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
  if err != nil {
    return 0, err
  }

  program := gl.CreateProgram()

  gl.AttachShader(program, vertexShader)
  gl.AttachShader(program, fragmentShader)
  gl.LinkProgram(program)

  var status int32
  gl.GetProgramiv(program, gl.LINK_STATUS, &status)
  if status == gl.FALSE {
    var logLength int32
    gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

    log := strings.Repeat(NULL_TERM, int(logLength+1))
    gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

    return 0, fmt.Errorf("failed to link program: %v", log)
  }

  gl.DeleteShader(vertexShader)
  gl.DeleteShader(fragmentShader)

  return program, nil
}

func readShaderFile(filepath string) (string, error) {
  data, err := ioutil.ReadFile(filepath)
  if err != nil {
    return "", err
  }
  return string(data) + NULL_TERM, nil
}

func NewShader(vertexPath string, fragmentPath string) (Shader, error) {
  vertexShaderSource, err := readShaderFile(vertexPath)
  if err != nil {
    return Shader{0}, err
  }
  fragmentShaderSource, err := readShaderFile(fragmentPath)
  if err != nil {
    return Shader{0}, err
  }

  program, err := configureProgram(vertexShaderSource, fragmentShaderSource)
  if err != nil {
    return Shader{0}, err
  }

  return Shader{program}, nil
}

func (s Shader) Use() {
  gl.UseProgram(s.ID)
}
