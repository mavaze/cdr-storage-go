package storage

import (
	"time"
)

// RotationPolicy interface to define rotation strategy
type RotationPolicy interface {
	// order() int
	rotate(temp *TempFile) bool
}

// FileRotator is a decorator consolidating multiple rotation stategies
type FileRotator struct {
	policies []RotationPolicy
}

func (f *FileRotator) addPolicy(policy RotationPolicy) {
	f.policies = append(f.policies, policy)
	// sort.Slice(f.policies[:], func(i, j int) bool {
	// 	return f.policies[i].order() < f.policies[j].order()
	// })
}

func (f *FileRotator) rotate(temp *TempFile) bool {
	var toBeRotated bool
	for _, policy := range f.policies {
		if toBeRotated = toBeRotated || policy.rotate(temp); toBeRotated {
			return true
		}
	}
	return false
}

// CountBasedRotator strategy to rotate file on reaching the number of cdrs written to file
type CountBasedRotator struct {
	maxCount uint64
}

func (c *CountBasedRotator) rotate(temp *TempFile) bool {
	//fmt.Println("Evaluating count based rotation policy")
	return temp.count >= c.maxCount
}

// VolumeBasedRotator strategy to rotate file on reaching size limit
type VolumeBasedRotator struct {
	maxSize int64
}

func (v *VolumeBasedRotator) rotate(temp *TempFile) bool {
	//fmt.Println("Evaluating volume based rotation policy")
	if info, err := temp.File.Stat(); err == nil {
		return info.Size() >= v.maxSize
	}
	return true
}

// TimeBasedRotator strategy to rotate file after an interval since first byte was written
type TimeBasedRotator struct {
	intervalInMillis int64
}

func (t *TimeBasedRotator) rotate(temp *TempFile) bool {
	//fmt.Println("Evaluating time based rotation policy")
	interval := time.Duration(t.intervalInMillis)
	return temp.createdAt.Add(interval * time.Millisecond).Before(time.Now())
}
