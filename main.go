package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard" // 用于处理键盘输入
)

// 游戏难度级别
const (
	DifficultyEasy   = 500 // 简单模式 - 500ms
	DifficultyNormal = 300 // 普通模式 - 300ms
	DifficultyHard   = 100 // 困难模式 - 100ms
)

// 默认游戏配置
const (
	defaultWidth  = 20
	defaultHeight = 10
	minWidth      = 10
	maxWidth      = 40
	minHeight     = 8
	maxHeight     = 20
	emptyCell     = "  " // 改为两个空格
	snakeBody     = "██" // 使用两个实心方块
	food          = "★ " // 食物后面加空格对齐
)

// GameConfig 游戏配置结构体
type GameConfig struct {
	width      int
	height     int
	difficulty time.Duration
}

// Point 结构体表示坐标点
type Point struct {
	x, y int // x和y坐标
}

// Game 结构体包含游戏的所有状态
type Game struct {
	snake     []Point    // 蛇身体，切片中的第一个元素是蛇头
	food      Point      // 食物的位置
	direction Point      // 蛇移动的方向
	score     int        // 当前得分
	gameOver  bool       // 游戏是否结束
	board     [][]string // 游戏板（二维数组）
	config    GameConfig // 游戏配置
}

// 获取用户输入的数字
func getNumericInput(prompt string, min, max int) int {
	var input string
	for {
		fmt.Print(prompt)
		fmt.Scanln(&input)

		// 检查是否要退出
		input = strings.TrimSpace(input)
		if input == "q" || input == "Q" {
			fmt.Println("\n已退出游戏!")
			os.Exit(0)
		}

		// 尝试转换为数字
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("输入无效，请输入一个介于 %d 和 %d 之间的数字，或输入 'q' 退出\n", min, max)
			continue
		}
		if num >= min && num <= max {
			return num
		}
		fmt.Printf("请输入一个介于 %d 和 %d 之间的数字，或输入 'q' 退出\n", min, max)
	}
}

// 获取游戏配置
func getGameConfig() GameConfig {
	var config GameConfig

	// 清屏
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println("=== 欢迎来到贪吃蛇游戏! ===")
	fmt.Println("\n【操作说明】")
	fmt.Println("- 在选择界面: 输入 'q' 退出游戏")
	fmt.Println("- 在游戏中: 按 'q' 或 'ESC' 退出游戏")
	fmt.Println("- WASD: 控制蛇的移动方向")

	fmt.Println("\n【难度选择】")
	fmt.Println("1 = 简单 (蛇移动较慢)")
	fmt.Println("2 = 普通 (蛇移动适中)")
	fmt.Println("3 = 困难 (蛇移动较快)")

	difficulty := getNumericInput("\n请选择难度 (1-3，或 'q' 退出): ", 1, 3)
	switch difficulty {
	case 1:
		config.difficulty = time.Duration(DifficultyEasy) * time.Millisecond
		fmt.Println("已选择: 简单模式")
	case 2:
		config.difficulty = time.Duration(DifficultyNormal) * time.Millisecond
		fmt.Println("已选择: 普通模式")
	case 3:
		config.difficulty = time.Duration(DifficultyHard) * time.Millisecond
		fmt.Println("已选择: 困难模式")
	}

	fmt.Println("\n【画面设置】")
	config.width = getNumericInput(fmt.Sprintf("请设置宽度 (%d-%d，或 'q' 退出): ", minWidth, maxWidth), minWidth, maxWidth)
	fmt.Printf("宽度已设置为: %d\n", config.width)

	config.height = getNumericInput(fmt.Sprintf("请设置高度 (%d-%d，或 'q' 退出): ", minHeight, maxHeight), minHeight, maxHeight)
	fmt.Printf("高度已设置为: %d\n", config.height)

	fmt.Println("\n游戏即将开始...")
	time.Sleep(2 * time.Second)

	// 清屏
	cmd = exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	return config
}

// NewGame 创建并初始化新游戏
func NewGame(config GameConfig) *Game {
	g := &Game{
		snake: []Point{
			{x: config.width / 2, y: config.height / 2}, // 蛇初始位置在屏幕中央
		},
		direction: Point{x: 1, y: 0}, // 初始移动方向向右
		board:     make([][]string, config.height),
		config:    config,
	}

	// 初始化游戏板
	for i := range g.board {
		g.board[i] = make([]string, config.width)
	}

	g.generateFood() // 生成第一个食物
	return g
}

// generateFood 在随机位置生成食物
func (g *Game) generateFood() {
	rand.Seed(time.Now().UnixNano()) // 设置随机数种子
	for {
		// 随机生成食物坐标
		x := rand.Intn(g.config.width)
		y := rand.Intn(g.config.height)
		occupied := false

		// 检查生成的位置是否与蛇身重叠
		for _, p := range g.snake {
			if p.x == x && p.y == y {
				occupied = true
				break
			}
		}

		// 如果位置未被占用，放置食物
		if !occupied {
			g.food = Point{x: x, y: y}
			break
		}
	}
}

// update 更新游戏状态
func (g *Game) update() {
	if g.gameOver {
		return
	}

	// 计算蛇头的新位置
	newHead := Point{
		x: g.snake[0].x + g.direction.x,
		y: g.snake[0].y + g.direction.y,
	}

	// 检查是否撞墙
	if newHead.x < 0 || newHead.x >= g.config.width || newHead.y < 0 || newHead.y >= g.config.height {
		g.gameOver = true
		return
	}

	// 检查是否撞到自己
	for _, p := range g.snake {
		if p.x == newHead.x && p.y == newHead.y {
			g.gameOver = true
			return
		}
	}

	// 移动蛇（在头部添加新位置）
	g.snake = append([]Point{newHead}, g.snake...)

	// 检查是否吃到食物
	if newHead.x == g.food.x && newHead.y == g.food.y {
		g.score++        // 增加分数
		g.generateFood() // 生成新的食物
	} else {
		// 如果没有吃到食物，删除尾部（保持长度不变）
		g.snake = g.snake[:len(g.snake)-1]
	}
}

// restart 重新开始游戏
func (g *Game) restart() {
	g.snake = []Point{
		{x: g.config.width / 2, y: g.config.height / 2}, // 蛇初始位置在屏幕中央
	}
	g.direction = Point{x: 1, y: 0} // 初始移动方向向右
	g.score = 0                     // 重置分数
	g.gameOver = false              // 重置游戏状态
	g.generateFood()                // 生成新的食物
}

// draw 绘制游戏界面
func (g *Game) draw() {
	// 清空游戏板
	for i := range g.board {
		for j := range g.board[i] {
			g.board[i][j] = emptyCell
		}
	}

	// 绘制蛇
	for _, p := range g.snake {
		if p.y >= 0 && p.y < g.config.height && p.x >= 0 && p.x < g.config.width {
			g.board[p.y][p.x] = snakeBody
		}
	}

	// 绘制食物
	g.board[g.food.y][g.food.x] = food

	// 清屏
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// 打印游戏信息
	fmt.Printf("分数: %d    难度: %s    大小: %dx%d\n",
		g.score,
		getDifficultyName(g.config.difficulty),
		g.config.width,
		g.config.height,
	)

	// 打印上边界
	fmt.Print("┌")
	for i := 0; i < g.config.width*2; i++ {
		fmt.Print("─")
	}
	fmt.Println("┐")

	// 打印游戏主体区域
	for _, row := range g.board {
		fmt.Print("│") // 左边界
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println("│") // 右边界
	}

	// 打印下边界
	fmt.Print("└")
	for i := 0; i < g.config.width*2; i++ {
		fmt.Print("─")
	}
	fmt.Println("┘")

	// 如果游戏结束，显示最终得分和重新开始选项
	if g.gameOver {
		fmt.Println("游戏结束! 最终得分:", g.score)
		fmt.Println("\n按 'r' 重新开始游戏")
		fmt.Println("按 'q' 或 'ESC' 退出游戏")
	}
}

// getDifficultyName 获取难度名称
func getDifficultyName(d time.Duration) string {
	switch d {
	case time.Duration(DifficultyEasy) * time.Millisecond:
		return "简单"
	case time.Duration(DifficultyNormal) * time.Millisecond:
		return "普通"
	case time.Duration(DifficultyHard) * time.Millisecond:
		return "困难"
	default:
		return "未知"
	}
}

func main() {
	// 获取游戏配置
	config := getGameConfig()

	fmt.Println("游戏控制:")
	fmt.Println("W: 向上移动")
	fmt.Println("S: 向下移动")
	fmt.Println("A: 向左移动")
	fmt.Println("D: 向右移动")
	fmt.Println("R: 重新开始游戏")
	fmt.Println("Q/ESC: 退出游戏")
	fmt.Println("\n准备开始游戏...")
	time.Sleep(2 * time.Second) // 给玩家时间阅读说明

	// 初始化键盘
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close() // 确保程序结束时关闭键盘

	game := NewGame(config)
	// 创建定时器，控制游戏速度
	ticker := time.NewTicker(config.difficulty)
	defer ticker.Stop()

	// 启动协程处理键盘输入
	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			// ESC键或Q键退出游戏
			if key == keyboard.KeyEsc || char == 'q' {
				fmt.Println("\n游戏已退出!")
				os.Exit(0)
			}

			// 处理方向键输入
			// 确保蛇不能直接向反方向移动
			switch char {
			case 'w': // 向上移动
				if game.direction.y != 1 {
					game.direction = Point{x: 0, y: -1}
				}
			case 's': // 向下移动
				if game.direction.y != -1 {
					game.direction = Point{x: 0, y: 1}
				}
			case 'a': // 向左移动
				if game.direction.x != 1 {
					game.direction = Point{x: -1, y: 0}
				}
			case 'd': // 向右移动
				if game.direction.x != -1 {
					game.direction = Point{x: 1, y: 0}
				}
			case 'r': // 重新开始游戏
				if game.gameOver {
					game.restart()
				}
			}
		}
	}()

	// 主游戏循环
	for range ticker.C {
		game.update() // 更新游戏状态
		game.draw()   // 绘制游戏界面
	}
}
