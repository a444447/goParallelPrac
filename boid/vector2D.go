package main

import "math"

type Vector2D struct {
	x, y float64
}

func (v1 Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{v1.x + v2.x, v1.y + v2.y}
}

func (v1 Vector2D) Subtract(v2 Vector2D) Vector2D {
	return Vector2D{v1.x - v2.x, v1.y - v2.y}
}

func (v1 Vector2D) Multiply(v2 Vector2D) Vector2D {
	return Vector2D{v1.x * v2.x, v1.y * v2.y}
}

func (v1 Vector2D) Division(v2 Vector2D) Vector2D {
	return Vector2D{v1.x / v2.x, v1.y / v2.y}
}

func (v1 Vector2D) AddScalar(s float64) Vector2D {
	return Vector2D{v1.x + s, v1.y + s}
}

func (v1 Vector2D) SubtractScalar(s float64) Vector2D {
	return Vector2D{v1.x - s, v1.y - s}
}

func (v1 Vector2D) MultiplyScalar(s float64) Vector2D {
	return Vector2D{v1.x * s, v1.y * s}
}

func (v1 Vector2D) DivisionScalar(s float64) Vector2D {
	return Vector2D{v1.x / s, v1.y / s}
}

func (v1 Vector2D) limit(lower, upper float64) Vector2D {
	return Vector2D{math.Min(math.Max(v1.x, lower), upper), math.Min(math.Max(v1.y, lower), upper)}
}

func (v1 Vector2D) Distance(v2 Vector2D) float64 {
	return math.Sqrt(math.Pow(v1.x-v2.x, 2) + math.Pow(v1.y-v2.y, 2))
}
