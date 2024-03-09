package period

import "slices"

type Period[T any] struct {
	Start T
	End   T
}

type Api[T any] struct {
	cmp func(T, T) int
}

func NewApi[T any](cmp func(T, T) int) *Api[T] {
	return &Api[T]{
		cmp: cmp,
	}
}

func (p *Api[T]) IsValidPeriod(period Period[T]) bool {
	return p.cmp(period.Start, period.End) < 0
}

func (p *Api[T]) ComparePeriods(a Period[T], b Period[T]) int {
	if d := p.cmp(a.Start, b.Start); d != 0 {
		return d
	}
	return p.cmp(a.End, b.End)
}

func (p *Api[T]) UnitePeriods(a Period[T], b Period[T]) Period[T] {
	start := a.Start
	if p.cmp(a.Start, b.Start) > 0 {
		start = b.Start
	}
	end := a.End
	if p.cmp(a.End, b.End) < 0 {
		end = b.End
	}
	return Period[T]{
		Start: start,
		End:   end,
	}
}

func (p *Api[T]) IntersectPeriods(a Period[T], b Period[T]) Period[T] {
	start := a.Start
	if p.cmp(a.Start, b.Start) < 0 {
		start = b.Start
	}
	end := a.End
	if p.cmp(a.End, b.End) > 0 {
		end = b.End
	}
	return Period[T]{
		Start: start,
		End:   end,
	}
}

func (p *Api[T]) SortAndUnitePeriods(periods []Period[T]) []Period[T] {
	if len(periods) < 2 {
		return periods
	}
	cloned := slices.Clone(periods)
	slices.SortFunc(cloned, p.ComparePeriods)
	lastPeriodIndex := 0
	for i := 1; i < len(cloned); i++ {
		lastPeriod := cloned[lastPeriodIndex]
		currentPeriod := cloned[i]
		if p.IsValidPeriod(p.IntersectPeriods(lastPeriod, currentPeriod)) {
			cloned[lastPeriodIndex] = p.UnitePeriods(lastPeriod, currentPeriod)
		} else {
			lastPeriodIndex++
			cloned[lastPeriodIndex] = currentPeriod
		}
	}
	return cloned[:lastPeriodIndex+1]
}

func (p *Api[T]) MakePeriodContainsCheck(period Period[T]) func(T) bool {
	return func(t T) bool {
		return p.cmp(period.Start, t) < 1 && p.cmp(t, period.End) < 1
	}
}

func (p *Api[T]) SubtractPeriods(a Period[T], b Period[T]) []Period[T] {
	intersection := p.IntersectPeriods(a, b)
	if !p.IsValidPeriod(intersection) {
		return []Period[T]{a}
	}
	isStartBeforeIntersection := p.cmp(a.Start, intersection.Start) < 0
	isEndAfterIntersection := p.cmp(intersection.End, a.End) < 0
	if isStartBeforeIntersection {
		if isEndAfterIntersection {
			return []Period[T]{
				{
					Start: a.Start,
					End:   intersection.Start,
				},
				{
					Start: intersection.End,
					End:   a.End,
				},
			}
		}
		return []Period[T]{
			{
				Start: a.Start,
				End:   intersection.Start,
			},
		}
	}
	if isEndAfterIntersection {
		return []Period[T]{
			{
				Start: intersection.End,
				End:   a.End,
			},
		}
	}
	return nil
}

func (p *Api[T]) SubtractPeriodsFromPeriods(
	periods []Period[T],
	periodsToSubtract []Period[T],
) []Period[T] {
	allPeriods := periods
	tmp := make([]Period[T], 0, len(periods))
	for _, breakPeriod := range periodsToSubtract {
		for _, period := range allPeriods {
			tmp = append(tmp, p.SubtractPeriods(period, breakPeriod)...)
		}
		allPeriods = tmp
	}
	return allPeriods
}
