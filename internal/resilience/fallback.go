package resilience

// Fallback executes the primary function and falls back to the fallback function if it fails
func Fallback(primary func() (interface{}, error), fallback func(error) (interface{}, error)) (interface{}, error) {
	result, err := primary()
	if err != nil {
		// Primary function failed, try the fallback
		return fallback(err)
	}
	return result, nil
}

// FallbackWithCondition executes the primary function and falls back to the fallback function
// if it fails and the condition function returns true
func FallbackWithCondition(
	primary func() (interface{}, error),
	fallback func(error) (interface{}, error),
	shouldFallback func(error) bool) (interface{}, error) {

	result, err := primary()
	if err != nil && shouldFallback(err) {
		// Primary function failed and condition is met, try the fallback
		return fallback(err)
	} else if err != nil {
		// Primary function failed but condition not met, return the error
		return nil, err
	}
	return result, nil
}
