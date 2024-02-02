package postgresql

import "fmt"

func ErrCommit(err error) error {
	return fmt.Errorf("failed to commit Tx due error: %v", err)
}

func ErrRollback(err error) error {
	return fmt.Errorf("failed to rollback Tx due error: %v", err)
}

func ErrCreateTx(err error) error {
	return fmt.Errorf("failed to create Tx due error: %v", err)
}

func ErrCreateQuery(err error) error {
	return fmt.Errorf("failed to create SQL Query due error: %v", err)
}

func ErrScan(err error) error {
	return fmt.Errorf("failed to scan due error: %v", err)
}

func ErrDoQuery(err error) error {
	return fmt.Errorf("failed to query due error: %v", err)
}
