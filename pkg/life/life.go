package life

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type World struct {
	Height int // высота сетки
	Width  int // ширина сетки
	Cells  [][]bool
}

func (w *World) SaveState(filename string) error {
	// Открываем файл для записи
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Проходим по всем клеткам и записываем их состояние в файл
	for i, row := range w.Cells {
		for _, cell := range row {
			var value byte
			if cell {
				value = '1'
			} else {
				value = '0'
			}
			// Записываем байт в файл
			if _, err := file.Write([]byte{value}); err != nil {
				return err
			}
		}
		// Добавляем символ новой строки после каждой строки сетки, кроме последней
		if i < len(w.Cells)-1 {
			if _, err := file.Write([]byte{'\n'}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *World) LoadState(filename string) error {
	// Открываем файл для чтения
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Читаем файл построчно
	scanner := bufio.NewScanner(file)
	var cells [][]bool
	var width int
	// Считываем каждую строку из файла
	for scanner.Scan() {
		line := scanner.Text()
		var row []bool
		// Проверяем длину строки
		if width == 0 {
			width = len(line)
		} else if len(line) != width {
			return fmt.Errorf("некорректные размеры сетки в файле")
		}
		// Парсим каждый символ строки и добавляем в текущую строку
		for _, char := range line {
			if char == '1' {
				row = append(row, true)
			} else if char == '0' {
				row = append(row, false)
			} else {
				return fmt.Errorf("некорректные символы в файле")
			}
		}
		cells = append(cells, row)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Устанавливаем размерность сетки и клетки
	w.Height = len(cells)
	w.Width = width
	w.Cells = cells

	return nil
}

func (w *World) String() string {
	var buf bytes.Buffer
	// Проходим по всем клеткам и формируем строку для отображения
	for _, row := range w.Cells {
		for _, cell := range row {
			if cell {
				buf.WriteString("◼️") // Черный квадрат для живой клетки
			} else {
				buf.WriteString("◻️") // Белый квадрат для мертвой клетки
			}
		}
		buf.WriteString("\n") // Переход на новую строку после каждой строки сетки
	}
	return buf.String()
}

func (w *World) Neighbours(x, y int) int {
	// Массив смещений для обхода соседних клеток
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	count := 0
	// Перебираем соседние клетки
	for i := 0; i < 8; i++ {
		nx := x + dx[i]
		ny := y + dy[i]

		// Проверяем, что соседние координаты находятся в пределах сетки
		if nx >= 0 && nx < w.Width && ny >= 0 && ny < w.Height {
			if w.Cells[ny][nx] {
				count++
			}
		}
	}
	return count
}

func (w *World) Next(x, y int) bool {
	n := w.Neighbours(x, y)      // получим количество живых соседей
	alive := w.Cells[y][x]       // текущее состояние клетки
	if n < 4 && n > 1 && alive { // если соседей двое или трое, а клетка жива
		return true // то следующее состояние — жива
	}
	if n == 3 && !alive { // если клетка мертва, но у неё трое соседей
		return true // клетка оживает
	}

	return false // в любых других случаях — клетка мертва
}

func (w *World) Seed() {
	// снова переберём все клетки
	for _, row := range w.Cells {
		for i := range row {
			//rand.Intn(10) возвращает случайное число из диапазона	от 0 до 9
			if rand.Intn(10) == 1 {
				row[i] = true
			}
		}
	}
}

func NextState(oldWorld, newWorld *World) {
	// переберём все клетки, чтобы понять, в каком они состоянии
	for i := 0; i < oldWorld.Height; i++ {
		for j := 0; j < oldWorld.Width; j++ {
			// для каждой клетки получим новое состояние
			newWorld.Cells[i][j] = oldWorld.Next(j, i)
		}
	}
}

func NewWorld(height, width int) *World {
	// создаём тип World с количеством слайсов hight (количество строк)
	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width) // создаём новый слайс в каждой строке
	}
	return &World{
		Height: height,
		Width:  width,
		Cells:  cells,
	}
}

// Start
func GoLife(height int, width int) {
	// объект для хранения текущего состояния сетки
	currentWorld := NewWorld(height, width)
	// объект для хранения следующего состояния сетки
	nextWorld := NewWorld(height, width)
	// установим начальное состояние
	currentWorld.Seed()
	for { // цикл для вывода каждого состояния
		// выведем текущее состояние на экран
		fmt.Println(currentWorld)
		// рассчитываем следующее состояние
		NextState(currentWorld, nextWorld)
		// изменяем текущее состояние
		currentWorld = nextWorld
		// делаем паузу
		time.Sleep(100 * time.Millisecond)
		// специальная последовательность для очистки экрана после каждого шага
		fmt.Print("\033[H\033[2J")
	}
}
