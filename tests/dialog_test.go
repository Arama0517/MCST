package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/AlecAivazis/survey/v2"
)

func TestDialoger(t *testing.T) {
	// 保存原始的标准输入
	originalStdin := os.Stdin
	// 创建一个新的 Reader，将其赋值给标准输入
	os.Stdin.Write([]byte("你的输入"))
	var name string
	prompt := &survey.Input{
		Message: "What is your name?",
	}
	survey.AskOne(prompt, &name)
	survey.
	os.Stdin = originalStdin

	fmt.Printf("Hello, %s!\n", name)
}
