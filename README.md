# Snake Game in Go

一个使用 Go 语言实现的终端版贪吃蛇游戏。

## 前置要求

- Go 1.21 或更高版本
- 终端支持 ANSI 转义序列

## 安装

1. 安装 Go（如果还没有安装）:
   - 访问 https://golang.org/dl/
   - 下载并安装适合你系统的 Go 版本

2. 克隆项目：
   ```bash
   git clone https://github.com/yourusername/snake-demo.git
   cd snake-demo
   ```

3. 安装依赖：
   ```bash
   go mod tidy
   ```

## 运行游戏

```bash
go run main.go
```

## 游戏控制

- `W`: 向上移动
- `S`: 向下移动
- `A`: 向左移动
- `D`: 向右移动
- `ESC`: 退出游戏

## 游戏规则

1. 使用 WASD 键控制蛇的移动方向
2. 吃到食物（★）可以增加分数和蛇的长度
3. 撞到墙壁或自己的身体会导致游戏结束
4. 尽可能获得更高的分数！