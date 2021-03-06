package tracker

import (
	"errors"
	"strconv"
	"time"
)

type Day struct {
	Start  time.Time `json:"start"`
	Tasks  []Task    `json:"tasks"`
	Pauses []Pause   `json:"pauses"`
}

type Days []Day

type Pause struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type Task struct {
	End         time.Time `json:"end"`
	Description string    `json:"description"`
}

func (d *Day) Task(desc string) error {
	if d.paused() {
		return errors.New("Day is currently paused")
	}
	t := Task{End: Now, Description: desc}
	d.Tasks = append(d.Tasks, t)
	return nil
}

func (d *Day) Pause() error {
	if d.paused() {
		return errors.New("Day is already paused")
	}
	p := Pause{Start: Now}
	d.Pauses = append(d.Pauses, p)
	return nil
}

func (d *Day) Resume() error {
	size := len(d.Pauses)
	if size == 0 {
		return errors.New("Day hasn't been paused yet")
	}
	last := d.Pauses[size-1]
	if !last.End.IsZero() {
		return errors.New("Day is not paused")
	}
	last.End = Now
	d.Pauses[size-1] = last
	return nil
}

func (d *Day) Duration() time.Duration {
	var dur time.Duration
	size := len(d.Tasks)
	if d.Start.IsZero() || size == 0 {
		return 0
	}
	task := d.Tasks[size-1]
	dur = task.End.Sub(d.Start)
	for _, p := range d.Pauses {
		if !p.End.IsZero() {
			dur -= p.End.Sub(p.Start)
		}
	}
	return dur
}

func (d *Day) Status() string {
	dur := d.Duration()
	h := strconv.Itoa(int(dur.Hours()))
	m := strconv.Itoa(int(dur.Minutes()) % 60)
	return h + "h" + m + "m"
}

func (d *Day) paused() bool {
	size := len(d.Pauses)
	if size == 0 {
		return false
	}
	last := d.Pauses[size-1]
	return last.End.IsZero()
}

func (d *Day) SameDay() bool {
	return false
}

func (slice Days) Len() int {
	return len(slice)
}

func (slice Days) Less(a, b int) bool {
	return slice[a].Start.Before(slice[b].Start)
}

func (slice Days) Swap(a, b int) {
	slice[a], slice[b] = slice[b], slice[a]
}
