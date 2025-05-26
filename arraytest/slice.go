package arraytest

import (
	"container/list"
	"fmt"
)

type Test struct {
	array []int
}

// 直接这样append，不会改变长度
func TestChangeLenth(array []int) {
	array = append(array, 1)
	array = append(array, 2)
	array = append(array, 3)
}

// 传指针可以改变长度，但是切片应该传指针吗？
func ChangeLenth(array *[]int) {
	*array = append(*array, 7)
	*array = append(*array, 8)
	*array = append(*array, 9)
}

// 直接改变索引是可以改变对应位置的值
func ChangeValue(array []int) {
	//must modify with index, if use range ,will not modify
	for i := 0; i < len(array); i++ {
		array[i] = array[i] + 1
	}

	// 这种不会修改array的值
	for _, val := range array {
		val *= 2
	}

	// 这样range + 索引也可以修改 array的值
	for index := range array {
		// 使用索引来修改原始切片中的值
		array[index] *= 2
	}
}

// map是可以改变值的
func ChangeMap(mp map[int]string) {
	mp[11] = "v11"
}

func ArrayTest() {
	fmt.Println("\033[1;32;40m  \n start ArrayTest--------------------------------------- \033[0m")
	test := &Test{}
	if test.array != nil {
		fmt.Println("array not nil")
	} else {
		fmt.Println("array is nil") //会到这里，未初始化的切片的默认值为nil
	}
	var testArray []int
	TestChangeLenth(testArray)
	fmt.Printf("1 after test change array len is:%d, val:%+v\n", len(testArray), testArray)

	ChangeLenth(&testArray)
	fmt.Printf("2 after change array len is:%d, val:%+v\n", len(testArray), testArray)

	ChangeValue(testArray)
	fmt.Printf("\n 3 after change value array len is:%d, val:%+v\n", len(testArray), testArray)

	mp := make(map[int]string)
	mp[0] = "v0"
	ChangeMap(mp)
	fmt.Printf("\n 1mp len is:%d,val:%+v\n", len(mp), mp)
}

// /////////////////////////////////////////////////////////////////////
func TwoDimensionalByteTest() {
	fmt.Println("\033[1;32;40m  \n start TwoDimensionalByteTest--------------------------------------- \033[0m")
	line := "........"
	var data [][]byte
	//这里是个大坑，切片是个指针，如果在这里初始化，然后再下面append，那么每行数据对应的指针是相同的，改变一行，每一行都变了(在c++里面可不会范这种低级错误)
	bline := []byte(line)
	for i := 0; i < 8; i++ {
		//bline := []byte(line)
		data = append(data, bline)
	}
	data[0][0] = 'Q'
	fmt.Printf("this init is error,change one char will change all line,data:%s\n", data)

	var data1 [][]byte
	for i := 0; i < 8; i++ {
		bline1 := []byte(line)
		data1 = append(data1, bline1)
	}
	data1[0][0] = 'Q'
	fmt.Printf("this change is ok,data1:%s\n", data1)
}

// /////////////////////////////////////////////////////////////////////
func MakeTest() {
	fmt.Println("\033[1;32;40m  \n start MakeTest--------------------------------------- \033[0m")
	data := make([][]byte, 5)
	for i := 0; i < 5; i++ {
		data[i] = make([]byte, 5)
	}
	fmt.Printf("make row:%d, col:%d, some pos:%d,data is:%s\n", len(data), len(data[0]), data[3][3], data)

	//map的key不支持数组
	//mp := make(map[[]int][]byte, 3)
	mp := make(map[int][]byte, 3)
	fmt.Printf("make len:%d, mp is:%+v\n", len(mp), mp)
}

// ////////////////////////////////////////////////////////////////////////////////////////
func StringTest() {
	fmt.Println("\033[1;32;40m  \n start StringTest--------------------------------------- \033[0m")
	str := "abcdefg"
	for i := 0; i < len(str); i++ {
		ch := str[i]
		fmt.Printf("%c,", ch)
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////
func QueueStackTest() {
	fmt.Println("\033[1;32;40m  \n start QueueStackTest--------------------------------------- \033[0m")
	stack := list.New()
	stack.PushBack(1)
	stack.PushBack(2)
	stack.PushBack(3)
	stack.PushBack(4)
	stack.PushBack(5)
	//fmt.Println(stack.Back().Value)
	for stack.Len() > 0 {
		element := stack.Back()
		val := stack.Remove(element)
		fmt.Println(element.Value.(int), "|", val)
		//val := stack.Remove(stack.Back())
		//fmt.Println("remove val:", val)
	}

	queue := list.New()
	queue.PushBack(1)
	queue.PushBack(2)
	queue.PushBack(3)
	queue.PushBack(4)
	queue.PushBack(5)
	for queue.Len() > 0 {
		element := queue.Front()
		val := queue.Remove(element)
		fmt.Println(element, "|", val)
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////
func floodFill(image [][]int, sr int, sc int, color int) [][]int {
	dir := []int{-1, 0, 1, 0, 0, -1, 0, 1} //上下左右四个方向
	queue := list.New()
	oldColor := image[sr][sc]
	image[sr][sc] = color
	queue.PushBack([]int{sr, sc})
	for queue.Len() > 0 {
		element := queue.Front()
		pos := element.Value.([]int)
		queue.Remove(element)
		//image[pos[0]][pos[1]] = color
		for i := 0; i < 4; i++ {
			newPosx := dir[i*2] + pos[0]
			newPosy := dir[i*2+1] + pos[1]
			if newPosx >= len(image) || newPosx < 0 ||
				newPosy >= len(image[0]) || newPosy < 0 || image[newPosx][newPosy] != oldColor {
				continue
			}
			if image[newPosx][newPosy] != color {
				continue
			}
			image[newPosx][newPosy] = color
			fmt.Printf("add pos x:%d,y:%d,old val:%d\n", newPosx, newPosy, image[newPosx][newPosy])
			queue.PushBack([]int{newPosx, newPosy})
		}
	}

	return image
}

func QueueTest() {
	fmt.Println("\033[1;32;40m  \n start QueueTest--------------------------------------- \033[0m")
	//[[1,1,1],[1,1,0],[1,0,1]]
	image := [][]int{{1, 1, 1}, {1, 1, 0}, {1, 0, 1}}
	ans := floodFill(image, 1, 1, 2)
	fmt.Printf("get result:%v\n", ans)
}

func MapAsSetTest() {
	fmt.Println("\033[1;32;40m  \n start MapAsSetTest--------------------------------------- \033[0m")
	//[[1,1,1],[1,1,0],[1,0,1]]
	mtest := make(map[uint64]struct{})
	mtest[1] = struct{}{}
	fmt.Printf("MapAsSetTest map:%v\n", mtest)
}
