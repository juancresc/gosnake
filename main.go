package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type position struct {
	x int
	y int
}

const (
	height int = 24
	width  int = 35
)

const (
	up    int = 119 // w
	down  int = 115 // s
	left  int = 97  // a
	right int = 100 // d
)

type snake struct {
	position  position
	direction int
	length    int
	sspeed    int
	colita    map[position]int
}

type comida struct {
	position position
}

func (s *snake) changeDirection(direction int) {
	switch direction {
	case up:
		if s.direction != down {
			s.direction = up
		}
	case down:
		if s.direction != up {
			s.direction = down
		}
	case left:
		if s.direction != right {
			s.direction = left
		}
	case right:
		if s.direction != left {
			s.direction = right
		}
	}
}

func mod(a, b int) int {
	m := a % b
	if a < 0 && b < 0 {
		m -= b
	}
	if a < 0 && b > 0 {
		m += b
	}
	return m
}

func (s *snake) move() {
	switch s.direction {
	case up:
		s.position.y = mod(s.position.y-1, height)
	case down:
		s.position.y = mod(s.position.y+1, height)
	case left:
		s.position.x = mod(s.position.x-1, width)
	case right:
		s.position.x = mod(s.position.x+1, width)
	}
}

func ticklogic(snake *snake, comidita *comida) {
	// check if snake is dead
	if _, ok := snake.colita[snake.position]; ok {
		fmt.Println("Game Over")
		os.Exit(0)
	}
	// check if snake ate
	if snake.position == comidita.position {
		snake.length++
		comidita.position = position{rand.Intn(width), rand.Intn(height)}
		if snake.sspeed > 50 {
			snake.sspeed -= 7
		} else if snake.sspeed < 20 && snake.sspeed > 10 {
			snake.sspeed -= 2
		}
	}
	// create colita
	snake.colita[snake.position] = snake.length
	// make colita slightly dissapear
	for key, colitaItem := range snake.colita {
		snake.colita[key] = colitaItem - 1 // colita TTL
		if snake.colita[key] == 0 {
			delete(snake.colita, key)
		}
	}
	// move snake
	snake.move()
	// print
	print(*snake, *comidita)
}

func print(snake snake, comidita comida) {
	// clear screen
	uniqueChar := "x"
	CallClear()
	fmt.Println(strings.Repeat(uniqueChar, width+2))
	// print snake
	for i := 0; i < height; i++ {
		fmt.Print(uniqueChar)
		for j := 0; j < width; j++ {
			if _, ok := snake.colita[position{j, i}]; ok {
				// colita
				fmt.Print(uniqueChar)
			} else if j == comidita.position.x && i == comidita.position.y {
				// comida
				fmt.Print(uniqueChar)
			} else if j == snake.position.x && i == snake.position.y {
				// head
				fmt.Print(uniqueChar)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println(uniqueChar)
	}
	fmt.Println(strings.Repeat(uniqueChar, width+2))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// init comida
	comidita := &comida{position{rand.Intn(width), rand.Intn(height)}}
	// init snake
	directions := []int{119, 115, 97, 100}
	randomDirection := rand.Intn(len(directions))
	initialDirection := directions[randomDirection]
	// frame rate
	sspeed := 150
	current_speed := sspeed
	mysnake := &snake{position{rand.Intn(width), rand.Intn(height)}, initialDirection, 1, sspeed, make(map[position]int)}
	//key reader
	keypress := make(chan int)
	go keyPress(keypress)
	tick := time.Tick(time.Duration(sspeed) * time.Millisecond)
	for {
		select {
		case <-tick:
			select {
			case keyCode := <-keypress:
				if keyCode == 27 {
					fmt.Println("Exit")
					os.Exit(0)
				}
				mysnake.changeDirection(keyCode)
				ticklogic(mysnake, comidita)
			default:
				// adjust frame rate
				if mysnake.sspeed != current_speed {
					tick = time.Tick(time.Duration(mysnake.sspeed) * time.Millisecond)
					current_speed = mysnake.sspeed
				}
				ticklogic(mysnake, comidita)
			}
		}
	}
}
