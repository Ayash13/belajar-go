package methods

import "fmt"

type Counter struct {
	Value int
}

// value receiver - works on a copy
func (c Counter) Print() {
	fmt.Println("Value is:", c.Value)
}

// pointer receiver - modifies the actual struct
func (c *Counter) Increment() {
	c.Value++
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

func Run() {
	fmt.Println("Counter Example:")
	c := Counter{Value: 0}
	c.Increment()
	c.Increment()
	c.Print()

	fmt.Println("\nRectangle Example:")
	rect := Rectangle{Width: 10, Height: 5}
	fmt.Println("Area:", rect.Area())
	
	rect.Scale(2)
	fmt.Println("Area after scaling by 2:", rect.Area())
}
