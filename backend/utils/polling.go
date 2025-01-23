/*
 * (c) 2025 Thales copyrights
 * This file is distributed under Apache-2.0 license.
 */

package utils

import (
	"context"
	"errors"
	"time"
)

// Poll polls a function periodically until it meets a success condition or times out.
func Poll(ctx context.Context, checkFunc func() (interface{}, error), successCondition func(result interface{}) bool, timeout time.Duration, frequency time.Duration) error {
	// Create a context with a timeout for the entire polling operation
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Loop until the context is done or the success condition is met
	for {
		select {
		case <-ctx.Done():
			return errors.New("timeout reached while polling")
		default:
			// Call the provided function to check the current status
			result, err := checkFunc()
			if err != nil {
				return err
			}

			// If the success condition is met, stop polling
			if successCondition(result) {
				return nil
			}

			// Wait for the specified frequency before the next poll
			select {
			case <-time.After(frequency):
				// Continue polling
			case <-ctx.Done():
				return errors.New("timeout reached while polling")
			}
		}
	}
}
