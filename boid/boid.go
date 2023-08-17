package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) calcAcceleratation() Vector2D {

	upper, lower := b.position.AddScalar(viewRadius), b.position.AddScalar(-viewRadius)
	avgPosition, avgVelcoity := Vector2D{0, 0}, Vector2D{0, 0}
	speration := Vector2D{0, 0}
	count := 0.0
	rwlock.RLock()
	for i := math.Max(0, lower.x); i <= math.Min(screenWidth, upper.x); i++ {
		for j := math.Max(0, lower.y); j <= math.Min(screenHeight, upper.y); j++ {
			if otherBoids := boidMap[int(i)][int(j)]; otherBoids != -1 && otherBoids != b.id {
				if dist := boids[otherBoids].position.Distance(b.position); dist < viewRadius {
					avgVelcoity = avgVelcoity.Add(boids[otherBoids].velocity)
					avgPosition = avgPosition.Add(boids[otherBoids].position)
					speration = speration.Add(b.position.Subtract(boids[otherBoids].position).DivisionScalar(dist))
					count++
				}
			}
		}
	}
	rwlock.RUnlock()
	accel := Vector2D{b.boraderBounce(b.position.x, screenWidth), b.boraderBounce(b.position.y, screenHeight)}
	accelV := Vector2D{0, 0}
	accelP := Vector2D{0, 0}
	accelS := Vector2D{0, 0}
	if count > 0 {
		avgVelcoity = avgVelcoity.DivisionScalar(count)
		avgPosition = avgPosition.DivisionScalar(count)
		accelV = avgVelcoity.Subtract(b.velocity).MultiplyScalar(adjRate)
		accelP = avgPosition.Subtract(b.position).MultiplyScalar(adjRate)
		accelS = speration.MultiplyScalar(adjRate)
		accel = accel.Add(accelP).Add(accelS).Add(accelV)
	}
	return accel
}

func (b *Boid) boraderBounce(pos, maxBorderPos float64) float64 {
	//pos与maxBorderPos- viewRadius的距离分别对应最左边和最右边
	if pos < viewRadius {
		return 1 / pos
	} else if pos > maxBorderPos-viewRadius {
		return 1 / (pos - maxBorderPos)
	}
	return 0
}

func (b *Boid) moveOne() {
	acceleration := b.calcAcceleratation()
	rwlock.Lock()
	b.velocity = b.velocity.Add(acceleration).limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	rwlock.Unlock()

}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(bid int) {
	boid := Boid{
		position: Vector2D{x: rand.Float64() * screenWidth, y: rand.Float64() * screenHeight},
		velocity: Vector2D{x: (rand.Float64() * 2) - 1.0, y: (rand.Float64() * 2) - 1.0},
		id:       bid,
	}
	boidMap[int(boid.position.x)][int(boid.position.y)] = bid
	boids[bid] = &boid
	go boid.start()
}
