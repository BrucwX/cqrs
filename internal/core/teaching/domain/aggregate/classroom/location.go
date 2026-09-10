package classroom

import (
	"errors"
	"fmt"
)

// Location 教室物理位置值对象（不可变）
type Location struct {
	building string
	floor    int
	room     string
}

func NewLocation(building string, floor int, room string) (Location, error) {
	if building == "" || room == "" {
		return Location{}, errors.New("building and room cannot be empty")
	}
	return Location{
		building: building,
		floor:    floor,
		room:     room,
	}, nil
}

func (l Location) Building() string { return l.building }
func (l Location) Floor() int      { return l.floor }
func (l Location) Room() string     { return l.room }
func (l Location) FullName() string {
	return fmt.Sprintf("%s-%dF-%s", l.building, l.floor, l.room)
}
