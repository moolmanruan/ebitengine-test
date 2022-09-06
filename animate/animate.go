package animate

import (
	"time"
)

//@keyframes
//cancel - channel
//duration - time.Duration
//delay - duration (pos or negative)
//iteration-count
// - number
//direction
// - normal - The animation is played as normal (forwards). This is default
// - reverse - The animation is played in reverse direction (backwards)
// - alternate - The animation is played forwards first, then backwards
// - alternate-reverse - The animation is played backwards first, then forwards
//timing-function
// - ease - Specifies an animation with a slow start, then fast, then end slowly (this is default)
// - linear - Specifies an animation with the same speed from start to end
// - ease-in - Specifies an animation with a slow start
// - ease-out - Specifies an animation with a slow end
// - ease-in-out - Specifies an animation with a slow start and end
// - cubic-bezier(n,n,n,n) - Lets you define your own values in a cubic-bezier function

//*value
//from
//to

// a cancel channel is returned - (cancel chan struct{})

var FPS60Interval = time.Second / 60
var FPS30Interval = time.Second / 30

type options struct {
	interval time.Duration
	delay    time.Duration
}

func defaultOptions() options {
	return options{
		interval: FPS60Interval,
	}
}

type Option func(opts *options)

func WithInterval(i time.Duration) Option {
	return func(opts *options) {
		opts.interval = i
	}
}
func WithDelay(d time.Duration) Option {
	return func(opts *options) {
		opts.delay = d
	}
}

// Int a value from one value to another.
// value is a pointer to the value that will be animated
// from and to are the start and endpoints of the animation
// duration is the duration of the animation
// interval is the time between updates
func Int(update func(int), from, to int, duration time.Duration, opts ...Option) (cancel func()) {
	oo := defaultOptions()
	for _, o := range opts {
		o(&oo)
	}
	if duration.Nanoseconds() <= 0 {
		update(to)
		return nil
	}
	var stop bool
	go animateInt(update, from, to, duration, &stop, oo)
	return func() { stop = true }
}

func animateInt(update func(int), from, to int, duration time.Duration, stop *bool, opts options) {
	startTime := time.Now().Add(opts.delay)
	value := from
	for {
		if *stop {
			return
		}

		runTime := time.Since(startTime)
		if runTime < 0 {
			time.Sleep(opts.interval)
			continue
		}

		if duration.Nanoseconds() > 0 {
			progress := float64(runTime.Nanoseconds()) / float64(duration.Nanoseconds())
			if progress >= 1 {
				value = to
				update(value)
				return
			}
		}

		value++
		if value > to {
			value = from
		}
		update(value)
		time.Sleep(opts.interval)
	}
}

// Value a value from one value to another.
// value is a pointer to the value that will be animated
// from and to are the start and endpoints of the animation
// duration is the duration of the animation
// interval is the time between updates
func Value(value *float64, to float64, duration time.Duration, opts ...Option) (cancel func()) {
	oo := defaultOptions()
	for _, o := range opts {
		o(&oo)
	}
	if duration.Nanoseconds() == 0 {
		*value = to
		return nil
	}
	var stop bool
	go animate(value, to, duration, &stop, oo)
	return func() { stop = true }
}

func animate(value *float64, to float64, duration time.Duration, stop *bool, opts options) {
	startTime := time.Now().Add(opts.delay)
	from := *value // copy start value
	valueRange := to - from
	for {
		if *stop {
			return
		}

		runTime := time.Since(startTime)
		if runTime < 0 {
			time.Sleep(opts.interval)
			continue
		}

		progress := float64(runTime.Nanoseconds()) / float64(duration.Nanoseconds())
		if progress >= 1 {
			*value = to
			return
		}
		*value = from + progress*valueRange
		time.Sleep(opts.interval)
	}
}
