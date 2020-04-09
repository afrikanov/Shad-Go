// +build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	for {
	loopGet:
		oldValue, err := c.Get(&kvapi.GetRequest{Key: key})
		var newValue string
		var apiError *kvapi.APIError
		var authError *kvapi.AuthError
	loopKeyNotFound:
		if errors.As(err, &authError) {
			return &kvapi.APIError{
				Method: "get",
				Err:    authError,
			}
		}
		if errors.Is(err, kvapi.ErrKeyNotFound) ||
			(errors.Is(err, errors.Unwrap(kvapi.ErrKeyNotFound)) && errors.Unwrap(kvapi.ErrKeyNotFound) != nil) {
			newValue, err = updateFn(nil)
		} else if errors.As(err, &apiError) {
			goto loopGet
		} else {
			newValue, err = updateFn(&oldValue.Value)
		}
		if err != nil {
			return err
		}
		if errors.As(err, &apiError) {
			goto loopGet
		}
		newVersion := uuid.Must(uuid.NewV4())
	loopSet:
		if oldValue != nil {
			_, err = c.Set(&kvapi.SetRequest{
				Key:        key,
				Value:      newValue,
				OldVersion: oldValue.Version,
				NewVersion: newVersion,
			})
		} else {
			_, err = c.Set(&kvapi.SetRequest{
				Key:        key,
				Value:      newValue,
				OldVersion: uuid.UUID{},
				NewVersion: newVersion,
			})
		}
		if errors.As(err, &authError) {
			return &kvapi.APIError{
				Method: "set",
				Err:    authError,
			}
		}
		var conflictError *kvapi.ConflictError
		if errors.As(err, &conflictError) {
			if newVersion == conflictError.ExpectedVersion {
				return nil
			}
			oldValue.Version = conflictError.ExpectedVersion
			goto loopGet
		}
		if errors.Is(err, kvapi.ErrKeyNotFound) ||
			(errors.Is(err, errors.Unwrap(kvapi.ErrKeyNotFound)) && errors.Unwrap(kvapi.ErrKeyNotFound) != nil) {
			if oldValue == nil {
				oldValue = &kvapi.GetResponse{
					Value:   "",
					Version: uuid.UUID{},
				}
			}
			oldValue.Value = newValue
			oldValue.Version = uuid.UUID{}
			goto loopKeyNotFound
		}
		if errors.As(err, &apiError) {
			goto loopSet
		}
		return nil
	}
}
