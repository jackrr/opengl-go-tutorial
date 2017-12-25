package main

import (
	"fmt"
	"runtime"
  "math"

  "projects/graphics-tutorial/shader"

  "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
  glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

  win, err := glfw.CreateWindow(800, 600, "Hello world", nil, nil)

  win.SetFramebufferSizeCallback(framebufferSizeCb)

  if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

  win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

  shader, err := shader.NewShader("shader/shader.vs", "shader/shader.fs")
	if err != nil {
		panic(err)
	}

  triangleVao := prepDrawTriangle()
  // rectVao := prepDrawRectangle()

	for !win.ShouldClose() {
    handleInput(win)
	  gl.ClearColor(0, 0.5, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    drawVertexArray(shader, triangleVao, 3)
    // drawElementBuffer(program, rectVao, 6)

		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func init() {
  runtime.LockOSThread()
  if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}
}

func framebufferSizeCb(window *glfw.Window, width int, height int) {
  gl.Viewport(0, 0, int32(width), int32(height))
}

func handleInput(window *glfw.Window) {
  if (window.GetKey(glfw.KeyEscape) == glfw.Press) {
    window.SetShouldClose(true)
  }
}

func prepDrawTriangle() uint32 {
  // vertices := []float32{
    // 0.5, -0.5, 0.0, // bottom right
    // -0.5, -0.5, 0.0, // bottom left
    // 0.0,  0.75, 0.0} // top
  vertices := []float32{
    0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // bottom right
    -0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // bottom left
    0.0,  0.75, 0.0, 0.0, 0.0, 1.0} // top

  var vao, vbo uint32

  gl.GenVertexArrays(1, &vao)
  gl.GenBuffers(1, &vbo)

  gl.BindVertexArray(vao)

  gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
  gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

  gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
  gl.EnableVertexAttribArray(0)

  gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
  gl.EnableVertexAttribArray(1)
  return vao
}

func prepDrawRectangle() uint32 {
  vertices := []float32{
    0.5, 0.5, 0.0, // top right
    0.5, -0.5, 0.0, // bottom right
    -0.5, -0.5, 0.0, // bottom left
    -0.5,  0.5, 0.0} // top left


  indices := []uint32{
    0, 1, 3, // first triangle
    1, 2, 3}  // second triangle

  var vao, vbo, ebo uint32

  gl.GenVertexArrays(1, &vao)
  gl.GenBuffers(1, &vbo)
  gl.GenBuffers(1, &ebo)

  gl.BindVertexArray(vao)
  gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
  gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

  gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
  gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

  gl.EnableVertexAttribArray(0)
  gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

  return vao
}

func setColor(shader shader.Shader) {
  time := glfw.GetTime()
  green := float32((math.Sin(time) / 2) + 0.5)
  red := float32((math.Cos(time) / 2) + 0.5)
  vertexColorLocation := gl.GetUniformLocation(shader.ID, gl.Str("ourColor" + "\x00"))
  gl.Uniform4f(vertexColorLocation, red, green, 0.0, 1.0)
}

func drawVertexArray(shader shader.Shader, vao uint32, vertexCount int32) {
  shader.Use()
  setColor(shader)
  gl.BindVertexArray(vao)
  gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)
}

func drawElementBuffer(program uint32, vao uint32, vertexCount int32) {
  gl.UseProgram(program)
  gl.BindVertexArray(vao)
  gl.DrawElements(gl.TRIANGLES, vertexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

